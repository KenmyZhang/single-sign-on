package store

import (
	"net/http"

	"github.com/KenmyZhang/single-sign-on/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MISSING_DOCTOR_ERROR = "store.mgo_doctor.missing_doctor.const"
)

type MgoDoctorStore struct {
	*MgoStore
}

func NewMgoDoctorStore(mgoStore *MgoStore) DoctorStore {
	return &MgoDoctorStore{mgoStore}
}

func DoctorC(s *MgoSession) *mgo.Collection {
	return s.DB().C("doctors")
}

func (s MgoDoctorStore) CreateIndexesIfNotExists() {
	doctorC := DoctorC(s.masterSession)
	doctorC.EnsureIndexKey("name")
}

func (s MgoDoctorStore) Save(doctor *model.Doctor) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}
		ms := s.GetWorkerSession()
		doctorC := DoctorC(ms)

		doctor.PreSave()

		if err := doctorC.Insert(doctor); err != nil {
			result.Err = model.NewLocAppError(
				"MgoDoctorStore.Save", "store.mgo_doctor_store.save.app_error", nil,
				"name="+doctor.Name+", "+err.Error(),
			)
		} else {
			result.Data = doctor
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoDoctorStore) Get(id string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		var doctor model.Doctor
		result := StoreResult{}

		if err := DoctorC(ms).FindId(id).One(&doctor); err != nil {
			if err == mgo.ErrNotFound {
				result.Err = model.NewAppError("MgoDoctorStore.Get",
					MISSING_DOCTOR_ERROR, nil, "doctor_id="+id, http.StatusNotFound,
				)
			} else {
				result.Err = model.NewLocAppError("MgoDoctorStore.Get",
					"store.mgo_doctor.get.app_error", nil, "doctor_id="+id+", "+err.Error(),
				)
			}
		} else {
			result.Data = &doctor
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}

func (s MgoDoctorStore) SearchDoctors(term string, offset, limit int) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		ms := s.GetWorkerSession()
		result := StoreResult{}
		doctors := &[]model.Doctor{}
		query := bson.M{"name": bson.M{
				"$regex":   "^" + term,
				"$options": "i",
				},
			}

		if err := DoctorC(ms).Find(query).Sort("createAt").Skip(offset).Limit(limit).All(doctors); err != nil {
			result.Err = model.NewLocAppError("MgoDoctorStore.SearchDoctors",
				"store.mgo_doctor_store.search_doctors.app_error", nil, "term="+term+", "+err.Error(),
			)
		} else {
			result.Data = doctors
		}

		storeChannel <- result
		close(storeChannel)
		ms.Close()
	}()

	return storeChannel
}
