package mq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"ddd/adaptor"
	"base/log"
)

/*
DOC:https://github.com/apache/rocketmq/tree/master/docs/cn
 */

type RocketMQConsumer struct {
	rocketmq.PushConsumer
	Topic string
}

func (C *RocketMQConsumer) Watch(newProcessor func([]byte) adaptor.MessageProcessor) error {
	//consumer.MessageSelector{Type: consumer.TAG, Expression:"Role_S1"}
	err := C.Subscribe(C.Topic, consumer.MessageSelector{}, func(ctx context.Context,
		msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//orderlyCtx, _ := primitive.GetOrderlyCtx(ctx)
		//fmt.Printf("orderly context: %v\n", orderlyCtx)
		//fmt.Printf("subscribe orderly callback: %v \n", msgs)
		for _, msg := range msgList{
			//if msg.GetTags() != "Role_S1"{
			//	continue
			//}

			newProcessor(msg.Body).Command()
			//msg.GetKeys()
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil{
		return err
	}
	return C.Start()
}

func (C *RocketMQConsumer) Stop() error {
	return C.Shutdown()
}



func NewRocketMQProducer()(saga *rmqProducer, err error){
	saga = &rmqProducer{}
	resolver := primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})
	saga.p, err = rocketmq.NewProducer(
		producer.WithNsResovler(resolver),
		producer.WithRetry(2),
	)
	err = saga.p.Start()

	return
}

type rmqProducer struct {
	ReplyTopic string
	p rocketmq.Producer

}

func (saga *rmqProducer)Stop()error  {
	return saga.p.Shutdown()
}


func (saga *rmqProducer)Produce(data []byte) error {
	msg := primitive.NewMessage(saga.ReplyTopic, data)
	res, err := saga.p.SendSync(context.Background(), msg)
	if err != nil {
		log.Errorf("send message error: %s\n", err)
	} else {
		log.Debugf("send message success: result=%s\n", res.String())
	}
	return err
}