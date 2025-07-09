package dep

//import (
//	"codo-cnmp/internal/conf"
//	"codo-cnmp/pb"
//	"github.com/go-kratos/kratos/v2/transport/grpc"
//)
//
//func NewGreeterClient(ctx context.Context, logger log.Logger, bc *conf.Bootstrap) (pb.GreeterClient, func(), error) {
//	c := bc.GreeterRpcConf
//	connGRPC, err := grpc.DialInsecure(
//		ctx,
//		grpc.WithEndpoint(c.Addr),
//		grpc.WithTimeout(c.Timeout.AsDuration()),
//	)
//
//	if err != nil {
//		return nil, nil, err
//	}
//
//	return pb.NewGreeterClient(connGRPC), func() {
//		connGRPC.Close()
//	}, nil
//}
