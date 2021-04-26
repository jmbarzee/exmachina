package pbconvert

import (
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
)

// NewPBNode Builds a Node for grpc requests
func NewPBNode(n node.Node) *pb.Node {
	id := n.GetID()
	pbNode := &pb.Node{
		ID:       id[:],
		Type:     n.GetType(),
		Children: NewPBNodeList(n.GetChildren()),
	}
	return pbNode
}

// NewPBNodeList Builds a NodeList for grpc requests
func NewPBNodeList(nodes []node.Node) []*pb.Node {
	pbNodes := make([]*pb.Node, len(nodes))
	for i, node := range nodes {
		pbNodes[i] = NewPBNode(node)
	}
	return pbNodes
}
