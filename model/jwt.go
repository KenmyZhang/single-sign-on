package model

import (
	"github.com/dgrijalva/jwt-go"
	"strings"
)

const (
	SESSION_COOKIE_TOKEN  = "AUTHTOKEN"
	SESSION_COOKIE_USER   = "USERID"
	SESSION_CACHE_SIZE    = 35000
	SESSION_PROP_PLATFORM = "platform"
	SESSION_PROP_OS       = "os"
	SESSION_PROP_BROWSER  = "browser"
)

type CustomClaims struct {
	UserId string            `json:"user_id"`
	Roles  string            `json:"roles"`
	Props  map[string]string `json:"props"`
	jwt.StandardClaims
}

func (me *CustomClaims) IsExpired() bool {
	if me.ExpiresAt <= 0 {
		return false
	}

	if GetMillis() > me.ExpiresAt {
		return true
	}

	return false
}

func (me *CustomClaims) SetExpireInDays(days int) {
	me.ExpiresAt = GetMillis() + (1000 * 60 * 60 * 24 * int64(days))
}

func (me *CustomClaims) AddProp(key string, value string) {
	if me.Props == nil {
		me.Props = make(map[string]string)
	}

	me.Props[key] = value
}

func (me *CustomClaims) GetUserRoles() []string {
	return strings.Fields(me.Roles)
}

//func (me *CustomClaims) RevokeCustomClaims() (string, error) {
// jwt.TimeFunc = func() time.Time {
// return time.Unix(0, 0)
// }
// token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// return j.SigningKey, nil
// })
// if err != nil {
// return "", err
// }
// if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
// jwt.TimeFunc = time.Now
// claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
// return j.CreateToken(*claims)
// }
// return "", TokenInvalid
//}
