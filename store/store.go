package store

import (
	"time"

	l4g "github.com/alecthomas/log4go"

	"github.com/KenmyZhang/single-sign-on/model"
)

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

func Must(sc StoreChannel) interface{} {
	r := <-sc
	if r.Err != nil {
		l4g.Close()
		time.Sleep(time.Second)
		panic(r.Err)
	}

	return r.Data
}

type Store interface {
	Audit() AuditStore
	System() SystemStore
	Doctor() DoctorStore
	UnactivatedAccount() UnactivatedAccountStore
	MarkSystemRanUnitTests()
	Close()
}

type DoctorStore interface {
	Save(doctor *model.Doctor) StoreChannel
	Get(id string) StoreChannel
	SearchDoctors(term string, offset, limit int) StoreChannel
}

type AuditStore interface {
	Save(audit *model.Audit) StoreChannel
	Get(user_id string, offset int, limit int) StoreChannel
	PermanentDeleteByUser(userId string) StoreChannel
}

type SystemStore interface {
	Save(system *model.System) StoreChannel
	SaveOrUpdate(system *model.System) StoreChannel
	Update(system *model.System) StoreChannel
	Get() StoreChannel
	GetByName(name string) StoreChannel
}

type UnactivatedAccountStore interface {
	Save(unactivatedAccount *model.UnactivatedAccount) StoreChannel
	Get(doctorId string) StoreChannel
}
