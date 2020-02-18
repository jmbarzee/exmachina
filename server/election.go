package server

import (
	"context"
	"time"
)

type Election struct {
	// Votest
	Votes      int
	Ballots    []Ballot
	Start      time.Time
	SelfRun    bool
	ClosePolls context.CancelFunc
}

type Ballot struct {
	Accept      bool
	Proficiency int32
	UUID        string
}

func (e Election) Succeeded(totalVoters int64, requiredPercentage float64) bool {
	numberAccepted := len(e.Ballots)
	percentAccepted := float64(numberAccepted) / float64(totalVoters)
	if percentAccepted > requiredPercentage {
		// Election can select a winner
		return true
	}
	return false
}

func (e Election) Failed(totalVoters int64, requiredPercentage float64) bool {
	numberRejected := e.Votes - len(e.Ballots)
	percentRejected := float64(numberRejected) / float64(totalVoters)
	if percentRejected > 1.0-requiredPercentage {
		// Election can NOT select a winner
		return true
	}
	return false
}

func (d *Domain) beginElection(ctx context.Context, serviceName string, election *Election) {
	d.elections[serviceName] = election
	go d.endElectionAfter(ctx, serviceName)
}

func (d *Domain) endElectionAfter(ctx context.Context, serviceName string) {
	ctx, cancel := context.WithTimeout(ctx, d.config.ServiceHierarchyConfig.ElectionTimeout)
	defer cancel()

	<-ctx.Done()

	d.debugf(debugLocksElections, "endElectionAfter() pre-lock()\n")
	d.electionsLock.Lock()
	{
		d.debugf(debugLocksElections, "endElectionAfter() in-lock()\n")
		if election, ok := d.elections[serviceName]; ok {
			if election.SelfRun {
				d.Logf("Owned election ended: %s\n", serviceName)
				election.ClosePolls()
				delete(d.elections, serviceName)
			} else {
				d.Logf("Unowned election ended: %s\n", serviceName)
				delete(d.elections, serviceName)
			}
		}
	}
	d.electionsLock.Unlock()
	d.debugf(debugLocksElections, "endElectionAfter() post-lock()\n")
}

func (d *Domain) hostElection(ctx context.Context, serviceName string) {
	holdElection := false
	var electionCtx context.Context
	var closePolls context.CancelFunc

	d.debugf(debugLocksElections, "hostElection() pre-lock-0()\n")
	d.electionsLock.Lock()
	{
		d.debugf(debugLocksElections, "hostElection() in-lock-0()\n")
		if _, ok := d.elections[serviceName]; ok {
			// pending election found
			d.Logf("Found a pending Election for: %s", serviceName)
			holdElection = false

		} else {
			// no known elections
			serviceConfig, err := d.serviceConfigFromName(serviceName)
			if err != nil {
				d.Logf("Could not find config for  \"%s\": %v", serviceName, err)
			} else {
				holdElection = true
				electionCtx, closePolls = context.WithCancel(ctx)
				election := Election{
					SelfRun:    true,
					Votes:      1,
					Start:      time.Now(),
					ClosePolls: closePolls,
					Ballots: []Ballot{
						Ballot{
							Accept:      true,
							Proficiency: d.getProficiencyForService(serviceConfig),
							UUID:        d.config.UUID,
						},
					},
				}
				d.beginElection(electionCtx, serviceName, &election)
				d.Logf("Holding Election for %s", serviceName)
			}
		}
	}
	d.electionsLock.Unlock()
	d.debugf(debugLocksElections, "hostElection() post-lock-0()\n")

	if !holdElection {
		return
	}

	ballots := make(chan Ballot)

	d.peerMap.Range(func(uuid string, peer *Peer) bool {
		d.debugf(debugLocks, "hostElection() pre-lock-0(%v)\n", uuid)
		peer.RLock()
		{
			go d.rpcOpenPosition(electionCtx, peer, serviceName, ballots)

		}
		peer.RUnlock()
		d.debugf(debugLocks, "hostElection() post-lock-0(%v)\n", uuid)
		return true
	})

	totalVoters := int64(d.peerMap.SizeEstimate() + 1)
	requiredPercentage := d.config.ServiceHierarchyConfig.RequiredVotePercentage
Loop:
	for {
		select {
		case ballot, ok := <-ballots:
			if !ok {
				// channel closed
				break Loop
			}
			electionHasConclusion := false
			d.debugf(debugLocksElections, "hostElection() pre-lock-1()\n")
			d.electionsLock.Lock()
			{
				d.debugf(debugLocksElections, "hostElection() in-lock-1()\n")
				if election, ok := d.elections[serviceName]; ok {
					election.Votes++

					// election found
					if ballot.Accept {
						election.Ballots = append(election.Ballots, ballot)
					}
					if election.Succeeded(totalVoters, requiredPercentage) {
						d.Logf("Election for \"%s\" Succeeded!\n", serviceName)
						electionHasConclusion = true
					}
					if election.Failed(totalVoters, requiredPercentage) {
						d.Logf("Election for \"%s\" Failed!\n", serviceName)
						electionHasConclusion = true
					}

				} else {
					// election ended somehow?
					d.Logf("Election ended unexpectedly, ma: %s", serviceName)
				}
			}
			d.electionsLock.Unlock()
			d.debugf(debugLocksElections, "hostElection() post-lock-1()\n")

			if electionHasConclusion {
				break Loop
			}

		case <-electionCtx.Done():
			// election ended
			break Loop
		}
	}

	// conclued election, either we have enough votes or it timedout
	winnerUUID := ""

	d.debugf(debugLocksElections, "hostElection() pre-lock-2()\n")
	d.electionsLock.Lock()
	{
		d.debugf(debugLocksElections, "hostElection() in-lock-2()\n")

		if election, ok := d.elections[serviceName]; ok {
			if election.Succeeded(totalVoters, requiredPercentage) {
				bestBallot := election.Ballots[0]
				for _, ballot := range election.Ballots {
					if ballot.Proficiency > bestBallot.Proficiency {
						bestBallot = ballot
					}
				}
				winnerUUID = bestBallot.UUID
			}
			delete(d.elections, serviceName)
		}
	}
	d.electionsLock.Unlock()
	d.debugf(debugLocksElections, "hostElection() post-lock-2()\n")

	if winnerUUID == "" {
		d.Logf("no winner found!")
		return
	}

	d.peerMap.Range(func(uuid string, peer *Peer) bool {
		d.debugf(debugLocks, "hostElection() pre-lock-1(%v)\n", uuid)
		peer.RLock()
		{
			go d.rpcClosePosition(electionCtx, peer, serviceName, winnerUUID == peer.UUID)
		}
		peer.RUnlock()
		d.debugf(debugLocks, "hostElection() post-lock-1(%v)\n", uuid)
		return true
	})

	if winnerUUID == d.config.UUID {
		serviceConfig, err := d.serviceConfigFromName(serviceName)
		if err == nil {
			err = d.startService(serviceConfig)
			if err != nil {
				d.Logf("Failed to start service after acceptance: %v", err)
			}
		}
	}

}
