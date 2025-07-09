package server

//import (
//	"codo-cnmp/internal/conf"
//	"codo-cnmp/internal/service"
//	v1 "codo-cnmp/pb"
//	"github.com/go-kratos/kratos/v2/log"
//	"github.com/go-kratos/kratos/v2/middleware/logging"
//	"github.com/go-kratos/kratos/v2/middleware/metrics"
//	"github.com/go-kratos/kratos/v2/middleware/recovery"
//	"github.com/go-kratos/kratos/v2/middleware/tracing"
//	"github.com/go-kratos/kratos/v2/transport/grpc"
//	"go.opentelemetry.io/otel/metric"
//	"go.opentelemetry.io/otel/trace"
//)
//
//// NewGRPCServer new a gRPC server.
//func NewGRPCServer(bc *conf.Bootstrap, greeter *service.GreeterService, logger log.Logger,
//	mp metric.MeterProvider, tp trace.TracerProvider,
//) (*grpc.Server, error) {
//	c := bc.Server
//	meter := mp.Meter("server.grpc")
//
//	counter, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
//	if err != nil {
//		return nil, err
//	}
//	seconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
//	if err != nil {
//		return nil, err
//	}
//
//	opts := []grpc.ServerOption{
//		grpc.Middleware(
//			recovery.Recovery(),
//			tracing.Server(
//				tracing.WithTracerProvider(tp),
//			),
//			logging.Server(logger),
//			metrics.Server(
//				metrics.WithRequests(counter),
//				metrics.WithSeconds(seconds),
//			),
//		),
//	}
//	if c.Grpc.Network != "" {
//		opts = append(opts, grpc.Network(c.Grpc.Network))
//	}
//	if c.Grpc.Addr != "" {
//		opts = append(opts, grpc.Address(c.Grpc.Addr))
//	}
//	if c.Grpc.Timeout != nil {
//		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
//	}
//	srv := grpc.NewServer(opts...)
//	v1.RegisterGreeterServer(srv, greeter)
//	return srv, nil
//}
