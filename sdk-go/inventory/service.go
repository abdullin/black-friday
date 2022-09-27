package inventory

import "sdk-go/protos"

type Service struct {
	protos.UnimplementedInventoryServiceServer
}

func NewService() *Service {
	return &Service{}
}
