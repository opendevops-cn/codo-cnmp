package service

import (
	"context"

	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
)

// GreeterService is a greeter service.
type GreeterService struct {
	pb.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{Message: "Hello " + g.Hello}, nil
}
