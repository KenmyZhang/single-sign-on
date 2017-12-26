package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mssola/user_agent"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

func AuthenticateUserForLogin(id, loginId, password string) (*model.User, *model.AppError) {
	if len(password) == 0 {
		err := model.NewLocAppError("AuthenticateUserForLogin",
			"api.user.login.blank_pwd.app_error", nil, "",
		)
		err.StatusCode = http.StatusBadRequest
		return nil, err
	}

	var user *model.User
	var err *model.AppError

	if len(id) != 0 {
		if user, err = GetUser(id); err != nil {
			err.StatusCode = http.StatusBadRequest
			return nil, err
		}
	} else {
		if user, err = GetUserForLogin(loginId); err != nil {
			return nil, err
		}
	}

	if user, err = authenticateUser(user, password); err != nil {
		return nil, err
	}

	return user, nil
}

func DoLogin(w http.ResponseWriter, r *http.Request, user *model.User, deviceId string) *model.AppError {
	ua := user_agent.New(r.UserAgent())

	plat := ua.Platform()
	if plat == "" {
		plat = "unknown"
	}

	os := ua.OS()
	if os == "" {
		os = "unknown"
	}

	bname, bversion := ua.Browser()
	if bname == "" {
		bname = "unknown"
	}

	if bversion == "" {
		bversion = "0.0"
	}

	JwtConfig := utils.Cfg.GetJWTService()
	maxAge := *JwtConfig.ExpireTimeLengthInDays * 60 * 60 * 24
	customClaims := model.CustomClaims{
		UserId: user.Id,
		Roles:  user.Roles,
		StandardClaims: jwt.StandardClaims{
			Issuer:    JwtConfig.Issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(maxAge) * time.Second).Unix(),
		},
	}
	customClaims.AddProp(model.SESSION_PROP_PLATFORM, plat)
	customClaims.AddProp(model.SESSION_PROP_OS, os)
	customClaims.AddProp(model.SESSION_PROP_BROWSER, fmt.Sprintf("%v/%v", bname, bversion))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	tokenString, err := token.SignedString([]byte(JwtConfig.Secret))
	if err != nil {
		return model.NewAppError("DoLogin", "app.do_login.signed_string.app_error", nil, "err="+err.Error(), model.StatusBadRequest)
	}

	w.Header().Set(model.HEADER_TOKEN, tokenString)

	secure := false
	if GetProtocol(r) == "https" {
		secure = true
	}

	expiresAt := time.Unix(model.GetMillis()/1000+int64(maxAge), 0)
	sessionCookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    tokenString,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
	}

	userCookie := &http.Cookie{
		Name:    model.SESSION_COOKIE_USER,
		Value:   user.Id,
		Path:    "/",
		MaxAge:  maxAge,
		Expires: expiresAt,
		Secure:  secure,
	}

	http.SetCookie(w, sessionCookie)
	http.SetCookie(w, userCookie)

	return nil
}

func GetProtocol(r *http.Request) string {
	if r.Header.Get(model.HEADER_FORWARDED_PROTO) == "https" || r.TLS != nil {
		return "https"
	} else {
		return "http"
	}
}
