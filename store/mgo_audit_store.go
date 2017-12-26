package store

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/KenmyZhang/single-sign-on/model"
)

type MgoAuditStore struct {
	*MgoStore
}

func NewMgoAuditStore(mgoStore *MgoStore) AuditStore {
	return &MgoAuditStore{mgoStore}
}

func AuditC(s *MgoSession) *mgo.Collection {
	return s.DB().C("audits")
}

func (s MgoAuditStore) CreateIndexesIfNotExists() {
	AuditC(s.masterSession).EnsureIndexKey("userId")
}

func (s MgoAuditStore) Save(audit *model.Audit) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}

		audit.Id = model.NewId()
		audit.CreateAt = model.GetMillis()

		if err := AuditC(ms).Insert(audit); err != nil {
			result.Err = model.NewLocAppError(
				"MgoAuditStore.Save", "store.mgo_audit.save.saving.app_error", nil,
				"user_id="+audit.UserId+" action="+audit.Action,
			)
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoAuditStore) Get(userId string, offset int, limit int) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}

		if limit > 1000 {
			limit = 1000
			result.Err = model.NewLocAppError(
				"MgoAuditStore.Get", "store.mgo_audit.get.limit.app_error", nil, "user_id="+userId,
			)
			storeChannel <- result
			close(storeChannel)
			ms.Close()
			return
		}

		q := bson.M{}
		if len(userId) > 0 {
			q["userId"] = userId
		}

		var audits model.Audits
		if err := AuditC(ms).Find(q).Sort("-createAt").Skip(offset).Limit(limit).All(&audits); err != nil {
			result.Err = model.NewLocAppError(
				"MgoAuditStore.Get", "store.mgo_audit.get.finding.app_error", nil, "user_id="+userId,
			)
		} else {
			result.Data = audits
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoAuditStore) PermanentDeleteByUser(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}

		if err := AuditC(ms).Remove(bson.M{"userId": userId}); err != nil {
			result.Err = model.NewLocAppError("MgoAuditStore.Delete",
				"store.mgo_audit.permanent_delete_by_user.app_error", nil, "user_id="+userId,
			)
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}
