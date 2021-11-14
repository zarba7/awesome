package modules

import "google.golang.org/protobuf/proto"

type IFace interface {
	Uid()uint64
	PublishDomainEvent(eTid uint32, eMsg proto.Message)
}
