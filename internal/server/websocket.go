package server

import (
	"codo-cnmp/internal/conf"
	"codo-cnmp/internal/service"
	_ "github.com/go-kratos/kratos/v2/log"

	"github.com/opendevops-cn/codo-golang-sdk/adapter/kratos/transport/websocket"
)

func NewWebSocketServer(bc *conf.Bootstrap, podLogSvc *service.PodLogWebsocketService,
	podCommandSvc *service.PodCommandWebsocketService) (*websocket.Server, error) {
	wsConf := bc.WS
	addr := wsConf.GetADDR()
	return websocket.NewServer(addr, websocket.WithHandlerBuilders(podLogSvc), websocket.WithHandlerBuilders(podCommandSvc))
}
