package inventory

import "go-lite/schema"

type Service struct {
	schema.UnimplementedInventoryServiceServer
}

func NewService() *Service {
	return &Service{}
}
