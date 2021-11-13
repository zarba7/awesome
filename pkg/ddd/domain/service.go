package domain

import (
	"ddd/adaptor"
	dto2 "ddd/dto"
	"ddd/po"
	"ddd/util"
	"fmt"
	"log"
	"reflect"
	"time"
)

func RunService(opts Options) error {
	var srv = service{
		repo : opts.Repo,
		namespace : opts.Namespace,
		ftySession : opts.SessionFactory,
		ftyEvent : opts.EventFactory,
		ftyCommand : opts.CommandFactory,
		outBox : opts.OutBox,
		codec : opts.Codec,
		input :opts.Input,
	}
	return srv.run()
}

type service struct {
	repo       adaptor.DomainRepo
	namespace  string
	cache      cache
	codec      Codec
	ftySession func(adaptor.DomainAggregate) Session
	ftyCommand util.SameNameStructFactory
	ftyEvent   util.SameNameStructFactory
	outBox     adaptor.DomainOutBox
	input  Input
}

func (srv *service) Quit() error {
	//srv.saga.Quit()
	return srv.input.Stop()
}

func (srv *service) run() error {
	go srv.cache.run(srv.repo)
	return srv.input.Start(srv.command)
}

func (srv *service) command(msg *dto2.Message) (err error) {
	agg, ok := srv.cache.Find(msg.AggID)
	if !ok {
		agg.srv = srv
		agg.po.ID = msg.AggID
		agg.po.Tag = srv.namespace
		agg.ss = srv.ftySession(agg)
	}
	agg.clear()
	if err = srv.load(agg); err != nil {
		return
	}
	var handler interface{}
	if handler, err = srv.ftyCommand.FindStructByUnmarshal(msg.Tid, msg.UnmarshalContent); err != nil {
		return
	}
	if msg.Result, err = agg.ss.HandleCommand(handler); err != nil {
		return
	}

	var evts []*dto2.Event

	for _, evt := range agg.events {
		applier := srv.ftyEvent.FindStructByFieldVal(evt.tid, evt.args)
		agg.ss.ApplyEvent(applier)
		agg.po.CurrVersion += 1

		poEvt := &po.DomainEvent{
			Uid:       agg.po.CurrVersion,
			TrxID:     msg.TrxID,
			CreatedAt: time.Now().UnixMilli(),
			Tid:       evt.tid,
			Tag:       fmt.Sprintf("%s:%T", agg.po.Tag, evt.args),
		}
		poEvt.Content, _ = srv.codec.Marshal(evt.args)

		agg.po.Events = append(agg.po.Events, poEvt)

		evts = append(evts, &dto2.Event{
			Uid:       poEvt.Uid,
			TrxID:     poEvt.TrxID,
			CreatedAt: poEvt.CreatedAt,
			Tid:       poEvt.Tid,
			Tag:       reflect.TypeOf(evt.args).String(),
			Content:   poEvt.Content,
			AggID:     agg.po.ID,
			AggTag:    agg.po.Tag,
		})
	}

	if err = srv.repo.Save(&agg.po); err != nil {
		srv.cache.Remove(agg)
		return
	}
	if err2 := srv.outBox.Publish(agg.po.Tag, evts); err2 != nil {
		log.Printf("Error %e", err2)
	}

	{

		//若保存版本不对 补发
		//todo outbox 补发
		//go func() {
		///
		//}
	}
	return
}

func (srv *service) load(agg *aggregate) (err error) {
	poAgg := &agg.po
	if poAgg.CurrVersion == 0 {
		if err = srv.repo.LoadAll(poAgg); err != nil {
			return
		}
		if poAgg.Snapshot != nil {
			if err = agg.ss.Load(poAgg.Snapshot); err != nil {
				return
			}
			poAgg.Snapshot = poAgg.Snapshot[:]
		}
	} else {
		if err = srv.repo.Load(poAgg); err != nil {
			return
		}
	}
	for _, evt := range poAgg.Events {
		applier, _ := srv.ftyEvent.FindStructByUnmarshal(evt.Tid, evt.UnmarshalContent)
		agg.ss.ApplyEvent(applier)
	}
	poAgg.Events = poAgg.Events[:]
	return
}
