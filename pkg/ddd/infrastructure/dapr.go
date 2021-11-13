package infrastructure

import (
	"bytes"
	"context"
	"ddd/internal/dto"
	"ddd/internal/util/mutex"
	"errors"
	nr "github.com/dapr/components-contrib/nameresolution"
	"github.com/dapr/go-sdk/service/common"
	sdk "github.com/dapr/go-sdk/service/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"sync"
)

type DaprProxy struct {
	mtx mutex.Factory
	msgPool sync.Pool
	commander func(msg *dto.Message) error
	srv common.Service
}

func (the *DaprProxy) command(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	if in == nil {
		err = errors.New("nil invocation parameter")
		return
	}
	msg := the.msgPool.Get().(*dto.Message)

	if err = msg.RequestUnmarshal(in.Data); err != nil {
		return
	}
	if err = the.doCommand(ctx, msg); err != nil {
		return
	}
	buf := bytes.NewBuffer(in.Data)
	if err = msg.ResponseMarshal(buf); err != nil {
		return
	}
	out = &common.Content{
		Data:        buf.Bytes(),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	msg.Reset()
	the.msgPool.Put(msg)
	return
}

func (the *DaprProxy) doCommand(ctx context.Context, msg *dto.Message) error {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		keys := md.Get("consistent-hash")
		if len(keys) > 0 {
			locker := the.mtx.LocalLocker(keys[0])
			locker.Lock()
			defer locker.Unlock()
		}
	}
	return the.commander(msg)
}

func (the *DaprProxy) Start(commander func(msg *dto.Message) error) (err error) {
	appID := os.Getenv(nr.AppID)
	port := os.Getenv(nr.AppPort)
	log.Printf("app port from env [%s:%s, %s:%s]", nr.AppID, appID, nr.AppPort, port)
	// create a Dapr service server
	the.srv, err = sdk.NewService("127.0.0.1:" + port)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
	// add a service to service invocation handler
	if commander != nil {
		the.commander = commander
		the.msgPool.New = func() interface{} { return &dto.Message{} }
		the.mtx = mutex.NewFactory()
		if err = the.srv.AddServiceInvocationHandler("command", the.command); err != nil {
			log.Fatalf("error adding invocation handler: %v", err)
		}
	}
	// start the server
	return the.srv.Start()
}

func (the *DaprProxy) Stop() (err error) {
	return the.srv.Stop()
}