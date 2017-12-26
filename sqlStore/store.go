package sqlStore

import (
	"time"

	"github.com/KenmyZhang/single-sign-on/model"
	l4g "github.com/alecthomas/log4go"
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
	User() UserStore
	System() SystemStore
	Token() TokenStore
	Close()
}

type UserStore interface {
	Save(user *model.User) StoreChannel
	Update(user *model.User, allowRoleUpdate bool) StoreChannel
	Get(id string) StoreChannel
	GetByEmail(email string) StoreChannel
	GetByUsername(username string) StoreChannel
	GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail, allowSignInWithMobile, weixinEnabled bool) StoreChannel
	GetTotalUsersCount() StoreChannel
	UpdateFailedPasswordAttempts(userId string, attempts int) StoreChannel
	GetByAuth(authData *string, authService string) StoreChannel
	UpdateLastPictureUpdate(userId string) StoreChannel
	GetProfileByMobile(mobile string) StoreChannel
	UpdatePassword(userId, hashedPassword string) StoreChannel
	SearchUsers(term string, options map[string]bool) StoreChannel
}

type TokenStore interface {
	Save(recovery *model.Token) StoreChannel
	Delete(token string) StoreChannel
	GetByToken(token string) StoreChannel
	GetByExtra(extra string) StoreChannel
	GetTokenCountByExtra(extra string) StoreChannel	
	Cleanup()
}

type SystemStore interface {
	Save(system *model.System) StoreChannel
	SaveOrUpdate(system *model.System) StoreChannel
	Update(system *model.System) StoreChannel
	Get() StoreChannel
	GetByName(name string) StoreChannel
}
