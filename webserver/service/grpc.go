package service

import (
	"context"

	"github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/ident"
	"github.com/jmbarzee/dominion/service/dominion"
	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/dominion/system/connect"
)

func (s WebServer) rpcGetDomains(ctx context.Context) ([]ident.DomainRecord, error) {
	rpcName := "GetDomains"
	domainRecords := []ident.DomainRecord{}

	err := s.Dominion.LatchWrite(func(dominion *dominion.Dominion) error {
		err := connect.CheckConnection(ctx, dominion)
		if err != nil {
			return err
		}

		serviceRequest := &grpc.Empty{}

		system.LogRPCf(rpcName, "Sending request")
		dominionClient := grpc.NewDominionClient(dominion.Conn)
		reply, err := dominionClient.GetDomains(ctx, serviceRequest)
		if err != nil {
			return err
		}
		system.LogRPCf(rpcName, "Received reply")

		domainRecords, err = ident.NewDomainRecordList(reply.GetDomainRecords())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return domainRecords, nil
}
