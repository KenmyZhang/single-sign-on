package store

import (
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	l4g "github.com/alecthomas/log4go"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

const (
	EXIT_DB_OPEN = 101
)

func ConvertIDToBsonID(id string) (bson.ObjectId, bool) {
	if bson.IsObjectIdHex(id) {
		return bson.ObjectIdHex(id), true
	}
	return bson.ObjectId(""), false
}

type MgoSession struct {
	session *mgo.Session
	dbname  string
}

func (ms *MgoSession) Copy() *MgoSession {
	return &MgoSession{
		session: ms.session.Copy(),
		dbname:  ms.dbname,
	}
}

func (ms *MgoSession) DB() *mgo.Database {
	return ms.session.DB(ms.dbname)
}

func (ms *MgoSession) Close() {
	ms.session.Close()
}

type MgoStore struct {
	masterSession      *MgoSession
	audit              AuditStore
	system             SystemStore
	doctor             DoctorStore
	unactivatedAccount UnactivatedAccountStore
	SchemaVersion      string
}

func NewMgoStore() Store {
	session, err := mgo.Dial(utils.Cfg.MgoSettings.DialURL)

	if err != nil {
		os.Exit(EXIT_DB_OPEN)
	}

	mgoStore := &MgoStore{
		masterSession: &MgoSession{
			session: session,
			dbname:  utils.Cfg.MgoSettings.DatabaseName,
		},
	}

	mgoStore.SchemaVersion = mgoStore.GetCurrentSchemaVersion()

	mgoStore.doctor = NewMgoDoctorStore(mgoStore)
	mgoStore.audit = NewMgoAuditStore(mgoStore)
	mgoStore.system = NewMgoSystemStore(mgoStore)
	mgoStore.unactivatedAccount = NewMgoUnactivatedAccountStore(mgoStore)

	mgoStore.doctor.(*MgoDoctorStore).CreateIndexesIfNotExists()
	mgoStore.audit.(*MgoAuditStore).CreateIndexesIfNotExists()
	mgoStore.system.(*MgoSystemStore).CreateIndexesIfNotExists()
	mgoStore.unactivatedAccount.(*MgoUnactivatedAccountStore).CreateIndexesIfNotExists()
	return mgoStore
}

func (s *MgoStore) GetWorkerSession() *MgoSession {
	return s.masterSession.Copy()
}

func (s *MgoStore) GetCurrentSchemaVersion() string {
	var system model.System
	s.masterSession.DB().C("systems").Find(bson.M{"name": "Version"}).One(&system)
	return system.Value
}

func (s *MgoStore) MarkSystemRanUnitTests() {
	if result := <-s.System().Get(); result.Err == nil {
		props := result.Data.(model.StringMap)
		unitTests := props[model.SYSTEM_RAN_UNIT_TESTS]
		if len(unitTests) == 0 {
			systemTests := &model.System{Name: model.SYSTEM_RAN_UNIT_TESTS, Value: "1"}
			<-s.System().Save(systemTests)
		}
	}
}

func (s *MgoStore) Close() {
	l4g.Info(utils.T("store.mgo.closing.info"))
	s.masterSession.Close()
}

func (s *MgoStore) Doctor() DoctorStore {
	return s.doctor
}

func (s *MgoStore) Audit() AuditStore {
	return s.audit
}

func (s *MgoStore) System() SystemStore {
	return s.system
}

func (s *MgoStore) UnactivatedAccount() UnactivatedAccountStore {
	return s.unactivatedAccount
}
