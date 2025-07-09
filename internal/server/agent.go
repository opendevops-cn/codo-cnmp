package server

//import (
//	"codo-agent/embed"
//	"context"
//	"fmt"
//	"codo-cnmp/internal/conf"
//)

//type AgentServer struct {
//	agent *embed.CodoAgentEmbed
//}
//
//func NewAgentServer(ctx context.Context, bc *conf.Bootstrap) (*AgentServer, error) {
//	cfg := bc.AGENT_SERVER
//	if cfg == nil {
//		print(">>> agent server config is nil \n")
//		return nil, nil
//	}
//	enabled := cfg.GetENABLED()
//	if !enabled {
//		print(">>> agent server config is disabled \n")
//		return nil, nil
//	}
//	serverAddr := cfg.GetSERVER_ADDR()
//	if serverAddr == "" {
//		return nil, fmt.Errorf("agent server server address is empty")
//	}
//	meshAddr := cfg.GetMESH_ADDR()
//	if meshAddr == "" {
//		return nil, fmt.Errorf("agent server mesh address is empty")
//	}
//	nodeType := cfg.GetNODE_TYPE()
//	if nodeType == "" {
//		nodeType = "normal"
//	}
//
//	logLevel := bc.OTEL.LOG.GetLEVEL()
//	if logLevel == "" {
//		logLevel = "INFO"
//	}
//	cf := embed.Config{
//		ServerAddress: serverAddr,
//		NodeType:      nodeType,
//		MeshAddr:      meshAddr,
//		LogLevel:      logLevel,
//	}
//
//	agent, err := embed.NewCodoAgentEmbed(ctx, cf)
//	if err != nil {
//		return nil, err
//	}
//	return &AgentServer{agent: agent}, nil
//}
//
//func (x *AgentServer) Start(ctx context.Context) error {
//	if x != nil && x.agent != nil {
//		return x.agent.Start(ctx)
//	}
//	return nil
//}
//
//func (x *AgentServer) Stop(ctx context.Context) error {
//	if x != nil && x.agent != nil {
//		return x.agent.Stop(ctx)
//	}
//	return nil
//}
