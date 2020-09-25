package dominion

import (
	"time"

	"github.com/jmbarzee/dominion/identity"
)

type Dominion struct {
	//DominionIdentity holds the identifying information of the service
	identity.DominionIdentity

	// LastContact is the last time a service replied to a rpc
	LastContact time.Time
}
