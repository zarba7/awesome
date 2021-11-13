package application

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"ddd"
	"ddd/dto"
	"base/log"
	"src/services/application/session"
)



type apiGateWay struct {
	newProcessor func()*session.Context
	rpc          *grpc.Server
}
func (the *apiGateWay) stop() {	the.rpc.GracefulStop()}
func (the *apiGateWay) start(s *Service, opt core.Option)  error  {
	tcpAddr, err := net.ResolveTCPAddr("tcp", opt.Application.RPCAddr)
	if err != nil {
		log.Error(err)
		return err
	}
	L, err := net.Listen("tcp", fmt.Sprint(":",tcpAddr.Port))
	if err != nil {
		log.Error(err)
		return err
	}

	the.rpc = grpc.NewServer()
	dto.RegisterApplicationServer(the.rpc, the)
	the.newProcessor = func() *session.Context {return &session.Context{s: s}}
	go the.rpc.Serve(L)
	return nil
}


func (the *apiGateWay) Dispatch(_ context.Context, pack *dto.AppPacket) (rsp *dto.AppResponse, err error) {
	proc :=  the.newProcessor()
	proc.From = 1
	proc.Request = pack
	rsp = proc.Query(pack)
	if rsp != nil {
		return
	}
	err = proc.Command()
	rsp = &proc.Result
	return
}

func (the *apiGateWay) Query(_ context.Context, pack *dto.AppPacket) (rsp *dto.AppResponse, err error) {
	rsp = the.newProcessor().Query(pack)
	if rsp != nil {
		return
	}
	err = fmt.Errorf("protocol no impl @%v", *pack)
	return
}
