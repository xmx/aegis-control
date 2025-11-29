package repository

import (
	"context"
	"log/slog"
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
	Firewall() Firewall
	FS() FS
	Setting() Setting
	VictoriaMetrics() VictoriaMetrics

	CreateIndex(ctx context.Context) error
}

func NewAll(db *mongo.Database, log *slog.Logger) All {
	return &allRepo{
		db:              db,
		log:             log,
		agent:           NewAgent(db),
		agentRelease:    NewAgentRelease(db),
		broker:          NewBroker(db),
		brokerRelease:   NewBrokerRelease(db),
		certificate:     NewCertificate(db),
		firewall:        NewFirewall(db),
		fs:              NewFS(db),
		setting:         NewSetting(db),
		victoriaMetrics: NewVictoriaMetrics(db),
	}
}

type allRepo struct {
	db  *mongo.Database
	log *slog.Logger

	agent           Agent
	agentRelease    AgentRelease
	broker          Broker
	brokerRelease   BrokerRelease
	certificate     Certificate
	firewall        Firewall
	fs              FS
	setting         Setting
	victoriaMetrics VictoriaMetrics
}

func (ar *allRepo) DB() *mongo.Database   { return ar.db }
func (ar *allRepo) Client() *mongo.Client { return ar.db.Client() }

func (ar *allRepo) Agent() Agent                     { return ar.agent }
func (ar *allRepo) AgentRelease() AgentRelease       { return ar.agentRelease }
func (ar *allRepo) Broker() Broker                   { return ar.broker }
func (ar *allRepo) BrokerRelease() BrokerRelease     { return ar.brokerRelease }
func (ar *allRepo) Certificate() Certificate         { return ar.certificate }
func (ar *allRepo) Firewall() Firewall               { return ar.firewall }
func (ar *allRepo) FS() FS                           { return ar.fs }
func (ar *allRepo) Setting() Setting                 { return ar.setting }
func (ar *allRepo) VictoriaMetrics() VictoriaMetrics { return ar.victoriaMetrics }

func (ar *allRepo) CreateIndex(ctx context.Context) error {
	rv := reflect.ValueOf(ar)
	for i := range rv.NumMethod() {
		mv := rv.Method(i)
		coll, idx := ar.reflectCall(mv, mv.Type())
		if idx == nil {
			continue
		}
		if err := idx.CreateIndex(ctx); err != nil {
			ar.log.Error("索引创建错误", "name", coll, "error", err)
			return err
		}

		ar.log.Info("索引创建完毕", "name", coll)
	}

	return nil
}

func (ar *allRepo) reflectCall(mv reflect.Value, mt reflect.Type) (string, indexCreator) {
	if mt.NumIn() != 0 || mt.NumOut() != 1 {
		return "", nil
	}

	rets := mv.Call(nil)
	if len(rets) != 1 {
		return "", nil
	}
	ret := rets[0]
	if ret.IsNil() || !ret.IsValid() {
		return "", nil
	}

	val := ret.Interface()
	ic, ok := val.(indexCreator)
	if !ok {
		return "", nil
	}
	coll := ret.Type().Name()
	if ni, yes := val.(interface{ Name() string }); yes {
		coll = ni.Name()
	}

	return coll, ic
}

type indexCreator interface {
	CreateIndex(context.Context) error
}
