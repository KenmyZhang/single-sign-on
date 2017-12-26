package store

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/KenmyZhang/single-sign-on/model"
	"net/http"
)

type MgoUnactivatedAccountStore struct {
	*MgoStore
}

func NewMgoUnactivatedAccountStore(mgoStore *MgoStore) UnactivatedAccountStore {
	return &MgoUnactivatedAccountStore{mgoStore}
}

func UnactivatedAccountC(s *MgoSession) *mgo.Collection {
	return s.DB().C("unactivated_accounts")
}

func (s MgoUnactivatedAccountStore) CreateIndexesIfNotExists() {
	unactivatedAccountC := UnactivatedAccountC(s.masterSession)
	unactivatedAccountC.EnsureIndexKey("username")
	unactivatedAccountC.EnsureIndexKey("doctorId")
	unactivatedAccountC.EnsureIndex(mgo.Index{
		Key:    []string{"doctorId"},
		Unique: true,
	})
}

func (s MgoUnactivatedAccountStore) Save(unactivatedAccount *model.UnactivatedAccount) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}
		ms := s.GetWorkerSession()

		unactivatedAccount.PreSave()

		if err := UnactivatedAccountC(ms).Insert(unactivatedAccount); err != nil {
			result.Err = model.NewLocAppError("MgoUnactivatedAccountStore.Save",
				"store.mgo_unactivated_account_store.save.app_error", nil, err.Error(),
			)
		} else {
			result.Data = unactivatedAccount
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoUnactivatedAccountStore) Get(doctorId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		var unactivatedAccount model.UnactivatedAccount
		result := StoreResult{}

		if err := UnactivatedAccountC(ms).Find(bson.M{"doctorId": doctorId}).One(&unactivatedAccount); err != nil {
			result.Err = model.NewAppError("MgoUnactivatedAccountStore.Get",
				"store.mgo_unactivated_account_store.get.app_error", nil, "doctorId=" + doctorId, http.StatusNotFound,
			)
		} else {
			result.Data = &unactivatedAccount
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

