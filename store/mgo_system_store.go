package store

import (
	mgo "gopkg.in/mgo.v2"

	"github.com/KenmyZhang/single-sign-on/model"
)

type MgoSystemStore struct {
	*MgoStore
}

func NewMgoSystemStore(mgoStore *MgoStore) SystemStore {
	return &MgoSystemStore{mgoStore}
}

func SystemC(s *MgoSession) *mgo.Collection {
	return s.DB().C("systems")
}

func (s MgoSystemStore) CreateIndexesIfNotExists() {
}

func (s MgoSystemStore) Save(system *model.System) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}

		if err := SystemC(ms).Insert(system); err != nil {
			result.Err = model.NewLocAppError(
				"MgoSystemStore.Save", "store.mgo_system.save.app_error", nil, "",
			)
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoSystemStore) SaveOrUpdate(system *model.System) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		systemC := SystemC(ms)
		result := StoreResult{}

		if n, _ := systemC.FindId(system.Name).Count(); n == 1 {
			if err := systemC.UpdateId(system.Name, system); err != nil {
				result.Err = model.NewLocAppError(
					"MgoSystemStore.SaveOrUpdate", "store.mgo_system.update.app_error", nil, "",
				)
			}
		} else {
			if err := systemC.Insert(system); err != nil {
				result.Err = model.NewLocAppError(
					"MgoSystemStore.SaveOrUpdate", "store.mgo_system.save.app_error", nil, "",
				)
			}
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoSystemStore) Update(system *model.System) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}

		if err := SystemC(ms).UpdateId(system.Name, system); err != nil {
			result.Err = model.NewLocAppError(
				"MgoSystemStore.Update", "store.mgo_system.update.app_error", nil, err.Error(),
			)
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoSystemStore) Get() StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}

		systems := []model.System{}
		props := make(model.StringMap)

		if err := SystemC(ms).Find(nil).All(&systems); err != nil {
			result.Err = model.NewLocAppError(
				"MgoSystemStore.Get", "store.mgo_system.get.app_error", nil, "",
			)
		} else {
			for _, prop := range systems {
				props[prop.Name] = prop.Value
			}
			result.Data = props
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoSystemStore) GetByName(name string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}

		var system model.System

		if err := SystemC(ms).FindId(name).One(&system); err != nil {
			result.Err = model.NewLocAppError(
				"MgoSystemStore.GetByName", "store.mgo_system.get_by_name.app_error", nil, "",
			)
		} else {
			result.Data = &system
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}
