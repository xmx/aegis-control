package repository

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type All interface {
	DB() *mongo.Database
	Client() *mongo.Client

	Agent() Agent
	AgentRelease() AgentRelease
	Broker() Broker
	BrokerRelease() BrokerRelease
	Certificate() Certificate
	FS() FS
	Setting() Setting

	CreateIndex(ctx context.Context) error
}

func NewAll(db *mongo.Database) All {
	return &allRepo{
		db:            db,
		agent:         NewAgent(db),
		agentRelease:  NewAgentRelease(db),
		broker:        NewBroker(db),
		brokerRelease: NewBrokerRelease(db),
		certificate:   NewCertificate(db),
		fs:            NewFS(db),
		setting:       NewSetting(db),
	}
}

type allRepo struct {
	db            *mongo.Database
	agent         Agent
	agentRelease  AgentRelease
	broker        Broker
	brokerRelease BrokerRelease
	certificate   Certificate
	fs            FS
	setting       Setting
}

func (ar *allRepo) DB() *mongo.Database   { return ar.db }
func (ar *allRepo) Client() *mongo.Client { return ar.db.Client() }

func (ar *allRepo) Agent() Agent                 { return ar.agent }
func (ar *allRepo) AgentRelease() AgentRelease   { return ar.agentRelease }
func (ar *allRepo) Broker() Broker               { return ar.broker }
func (ar *allRepo) BrokerRelease() BrokerRelease { return ar.brokerRelease }
func (ar *allRepo) Certificate() Certificate     { return ar.certificate }
func (ar *allRepo) FS() FS                       { return ar.fs }
func (ar *allRepo) Setting() Setting             { return ar.setting }

func (ar *allRepo) CreateIndex(ctx context.Context) error {
	rv := reflect.ValueOf(ar)
	rt := rv.Type()
	for i := range rv.NumMethod() {
		mv := rv.Method(i)
		mt := rt.Method(i)
		ci := ar.reflectCall(mt, mv)
		if ci == nil {
			continue
		}

		if err := ci.CreateIndex(ctx); err != nil {
			return err
		}
	}

	return nil
}

type IndexCreator interface {
	CreateIndex(context.Context) error
}

func (ar *allRepo) reflectCall(method reflect.Method, mv reflect.Value) IndexCreator {
	mt := method.Type
	if mt.NumIn() != 1 || mt.NumOut() != 1 {
		return nil
	}

	rets := mv.Call([]reflect.Value{})
	if len(rets) != 1 {
		return nil
	}
	ret := rets[0]
	val := ret.Interface()
	if ic, ok := val.(IndexCreator); ok {
		return ic
	}

	return nil
}
