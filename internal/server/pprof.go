package server

import (
	"context"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync/atomic"

	"codo-cnmp/internal/conf"
)

type PprofServer struct {
	conf *conf.PprofConfig

	started  uint32
	listener net.Listener
}

func NewPprofServer(bc *conf.Bootstrap) (*PprofServer, error) {
	c := bc.PPROF
	svr := &PprofServer{
		conf: c,
	}
	if svr.conf.GetENABLE() {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
		addr := svr.conf.GetADDR()
		listener, err := net.Listen(svr.conf.GetNETWORK(), addr)
		if err != nil {
			return nil, err
		}
		svr.listener = listener
	}
	return svr, nil
}

func (x *PprofServer) Start(ctx context.Context) error {
	if x.listener == nil {
		return nil
	}
	if !atomic.CompareAndSwapUint32(&x.started, 0, 1) {
		return nil
	}
	return http.Serve(x.listener, nil)
}

func (x *PprofServer) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&x.started, 1, 0) {
		return x.listener.Close()
	}
	return nil
}
