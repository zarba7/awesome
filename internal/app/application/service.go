package application

import "C"
import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"ddd"
	"ddd/adaptor"
	"ddd/application"
	"ddd/dto"
	"base/json"
	"base/log"
	"base/mq"
	"src/services/application/session"

	_ "dubbo.apache.org/dubbo-go/v3/registry/etcdv3"
)

func Main(){


}


type Service struct {
	protocol application.Dispatcher
	consumer adaptor.Consumer
	api      apiGateWay
}

func (Srv *Service) OnInit(opt core.Option) error {
	opt.Application.Name = "Application"
	opt.RocketMQ.NameServer = []string{"127.0.0.1:9876"}
	Srv.protocol = protocolTable()
	var err error
	if err = Srv.newConsumer(opt); err != nil {
		log.Error(err)
		return err
	}
	if err = Srv.api.start(Srv, opt); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (Srv *Service) OnQuit() {
	Srv.api.stop()
	var err error
	err = Srv.consumer.Stop()
	log.Error(err)
}

func (Srv *Service)newProcessor(data []byte) adaptor.MessageProcessor {
	proc :=  &session.Context{}
	proc.Request = &dto.AppPacket{}
	json.UnmarshalPanic(data, proc.Request)
	return proc
}

func (Srv *Service)newContext(data []byte) *session.Context {
	proc :=  &session.Context{}
	var pack =
	C := the.s.dispatcher.FindCommandHandler(pack.Tid, the.from)
	if C != nil{
		json.UnmarshalPanic(pack.Content, C)
		C.Command(&the.Result)
		the.Result.Packet = nil
	}
	return proc
}

func (Srv *Service)newConsumer(opt core.Option) error  {
	var	C = &mq.RocketMQConsumer{}
	var err error
	C.Topic = opt.Application.Name
	resolver := primitive.NewPassthroughResolver(opt.RocketMQ.NameServer)
	C.PushConsumer, err = rocketmq.NewPushConsumer(
		consumer.WithNamespace(opt.NameSpace),
		consumer.WithGroupName(C.Topic),
		consumer.WithNsResovler(resolver),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		//consumer.WithConsumerOrder(true),
	)
	if err = C.Watch(Srv.newProcessor); err != nil {
		log.Error(err)
		return err
	}
	Srv.consumer = C
	return err
}