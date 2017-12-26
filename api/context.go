package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	l4g "github.com/alecthomas/log4go"
	"github.com/dgrijalva/jwt-go"
	goi18n "github.com/nicksnyder/go-i18n/i18n"

	"github.com/KenmyZhang/single-sign-on/app"
	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

type Context struct {
	CustomClaims  model.CustomClaims
	TokenString   string
	Params        *ApiParams
	Err           *model.AppError
	T             goi18n.TranslateFunc
	RequestId     string
	IpAddress     string
	Path          string
	siteURLHeader string
}

func ApiHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{
		handleFunc:          h,
		requireCustomClaims: false,
		trustRequester:      false,
		isApi:               true,
	}
}

func ApiCustomClaimsRequired(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{
		handleFunc:          h,
		requireCustomClaims: true,
		trustRequester:      false,
		isApi:               true,
	}
}

func AppHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{
		handleFunc:          h,
		requireCustomClaims: false,
		trustRequester:      false,
		isApi:               false,
	}
}

type handler struct {
	handleFunc          func(*Context, http.ResponseWriter, *http.Request)
	requireCustomClaims bool
	trustRequester      bool
	isApi               bool
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l4g.Debug("%v - %v", r.Method, r.URL.Path)

	c := &Context{}
	c.T, _ = utils.GetTranslationsAndLocale(w, r)
	c.RequestId = model.NewId()
	c.IpAddress = utils.GetIpAddress(r)
	c.Params = ApiParamsFromRequest(r)

	tokenString := ""

	authHeader := r.Header.Get(model.HEADER_AUTH)
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == model.HEADER_BEARER {
		tokenString = authHeader[7:]
	} else if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == model.HEADER_TOKEN {
		tokenString = authHeader[6:]
	}

	if len(tokenString) == 0 {
		if cookie, err := r.Cookie(model.SESSION_COOKIE_TOKEN); err == nil {
			tokenString = cookie.Value

			if h.requireCustomClaims && !h.trustRequester {
				if r.Header.Get(model.HEADER_REQUESTED_WITH) != model.HEADER_REQUESTED_WITH_XML {
					c.Err = model.NewLocAppError("ServeHTTP",
						"api.context.session_expired.app_error", nil,
						"tokenString="+tokenString+" Appears to be a CSRF attempt",
					)
					tokenString = ""
				}
			}
		}
	}

	if len(tokenString) == 0 {
		tokenString = r.URL.Query().Get("access_token")
	}

	c.SetSiteURLHeader(app.GetProtocol(r) + "://" + r.Host)

	w.Header().Set(model.HEADER_REQUEST_ID, c.RequestId)
	w.Header().Set(model.HEADER_VERSION_ID, fmt.Sprintf("%v.%v.%v", model.CurrentVersion, model.BuildNumber, utils.ClientCfgHash))

	if !h.isApi {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'self'")
	} else {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			w.Header().Set("Expires", "0")
		}
	}

	if len(tokenString) != 0 && h.requireCustomClaims {
		token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(utils.Cfg.JWTSettings.Secret), nil
		})

		if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
			c.CustomClaims = *claims
			c.TokenString = tokenString
		} else {
			c.Err = model.NewLocAppError("ServeHTTP",
				"api.context.jwt_parse.app_error", nil, "err="+err.Error()+",tokenString="+tokenString,
			)
			c.Err.StatusCode = http.StatusUnauthorized
		}
	}

	c.Path = r.URL.Path

	if c.Err == nil && h.requireCustomClaims {
		c.CustomClaimsRequired()
	}

	if c.Err == nil {
		h.handleFunc(c, w, r)
	}

	if c.Err != nil {
		c.Err.Translate(c.T)
		c.Err.RequestId = c.RequestId
		c.LogError(c.Err)
		c.Err.Where = r.URL.Path

		if !*utils.Cfg.ServiceSettings.EnableDeveloper {
			c.Err.DetailedError = ""
		}

		if h.isApi {
			w.WriteHeader(c.Err.StatusCode)
			RenderJson(w, c.Err)
		} else {
			if c.Err.StatusCode == http.StatusUnauthorized {
				http.Redirect(w, r,
					c.GetSiteURLHeader()+"/?redirect="+url.QueryEscape(r.URL.Path), http.StatusTemporaryRedirect,
				)
			} else {
				utils.RenderWebError(c.Err, w, r)
			}
		}
	}
}

func (c *Context) SetSiteURLHeader(url string) {
	c.siteURLHeader = strings.TrimRight(url, "/")
}

func (c *Context) GetSiteURLHeader() string {
	return c.siteURLHeader
}

func (c *Context) RemoveCustomClaimsCookie(w http.ResponseWriter, r *http.Request) {
	sessionCookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	userCookie := &http.Cookie{
		Name:   model.SESSION_COOKIE_USER,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, sessionCookie)
	http.SetCookie(w, userCookie)
}

func (c *Context) SetInvalidParam(parameter string) {
	c.Err = NewInvalidParamError(parameter)
}

func (c *Context) SetInvalidUrlParam(parameter string) {
	c.Err = NewInvalidUrlParamError(parameter)
}

func NewInvalidParamError(parameter string) *model.AppError {
	err := model.NewLocAppError("Context",
		"api.context.invalid_body_param.app_error", map[string]interface{}{"Name": parameter}, parameter,
	)
	err.StatusCode = http.StatusBadRequest
	return err
}

func NewInvalidUrlParamError(parameter string) *model.AppError {
	err := model.NewLocAppError("Context",
		"api.context.invalid_url_param.app_error", map[string]interface{}{"Name": parameter}, "",
	)
	err.StatusCode = http.StatusBadRequest
	return err
}

func (c *Context) IsSystemAdmin() bool {
	return app.CustomClaimsHasPermissionTo(c.CustomClaims, model.PERMISSION_MANAGE_SYSTEM)
}

func (c *Context) CustomClaimsRequired() {
	if len(c.CustomClaims.UserId) == 0 {
		c.Err = model.NewAppError("",
			"api.context.session_expired.app_error", nil, "UserRequired", http.StatusUnauthorized,
		)
		return
	}
}

func (c *Context) RequireUserId() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.UserId == model.ME {
		c.Params.UserId = c.CustomClaims.UserId
	}

	if len(c.Params.UserId) < 24 {
		c.SetInvalidUrlParam("user_id")
	}

	return c
}

func (c *Context) RequireService() *Context {
	if c.Err != nil {
		return c
	}

	if len(c.Params.Service) == 0 {
		c.SetInvalidUrlParam("service")
	}

	return c
}

func (c *Context) LogError(err *model.AppError) {
	if err.Id == "web.check_browser_compatibility.app_error" {
		c.LogDebug(err)
	} else {
		l4g.Error(utils.TDefault("api.context.log.error"), c.Path, err.Where, err.StatusCode,
			c.RequestId, c.CustomClaims.UserId, c.IpAddress, err.SystemMessage(utils.TDefault), err.DetailedError,
		)
	}
}

func (c *Context) LogDebug(err *model.AppError) {
	l4g.Debug(utils.TDefault("api.context.log.error"), c.Path, err.Where, err.StatusCode,
		c.RequestId, c.CustomClaims.UserId, c.IpAddress, err.SystemMessage(utils.TDefault), err.DetailedError,
	)
}

func IsApiCall(r *http.Request) bool {
	return strings.Index(r.URL.Path, "/api/") == 0
}

func (c *Context) SetPermissionError(permission *model.Permission) {
	c.Err = model.NewLocAppError("Permissions", "api.context.permissions.app_error", nil, "userId="+c.CustomClaims.UserId+", "+"permission="+permission.Id)
	c.Err.StatusCode = http.StatusForbidden
}
