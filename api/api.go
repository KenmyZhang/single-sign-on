package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"

	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/mux"

	"github.com/KenmyZhang/single-sign-on/app"
	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

type Routes struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'sso'

	Users *mux.Router // 'sso/users'
	User  *mux.Router // 'sso/users/{user_id:[A-Za-z0-9]+}'
}

var BaseRoutes *Routes

func InitRouter() {
	app.Srv.Router = mux.NewRouter()
	app.Srv.Router.NotFoundHandler = http.HandlerFunc(Handle404)
}

func InitApi() {
	BaseRoutes = &Routes{}
	BaseRoutes.Root = app.Srv.Router
	BaseRoutes.ApiRoot = app.Srv.Router.PathPrefix(model.API_URL_SUFFIX).Subrouter()

	BaseRoutes.Users = BaseRoutes.ApiRoot.PathPrefix("/users").Subrouter()
	BaseRoutes.User = BaseRoutes.ApiRoot.PathPrefix("/users/{user_id:[A-Za-z0-9]+}").Subrouter()

	InitUser()
	InitSystem()
	InitOauth()

	utils.InitHTML()

	app.Srv.Router.Handle("/sso/{anything:.*}", http.HandlerFunc(Handle404))
}

func HandleEtag(etag string, routeName string, w http.ResponseWriter, r *http.Request) bool {
	if et := r.Header.Get(model.HEADER_ETAG_CLIENT); len(etag) > 0 {
		if et == etag {
			w.Header().Set(model.HEADER_ETAG_SERVER, etag)
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}

	return false
}

func Handle404(w http.ResponseWriter, r *http.Request) {
	err := model.NewLocAppError("Handle404", "api.context.404.app_error", nil, "")
	err.Translate(utils.T)
	err.StatusCode = http.StatusNotFound

	l4g.Debug("%v: code=404 ip=%v", r.URL.Path, utils.GetIpAddress(r))

	w.WriteHeader(err.StatusCode)
	err.DetailedError = "There doesn't appear to be an api call for the url='" + r.URL.Path + "'."
	RenderJson(w, err)
}

func ReturnStatusOK(w http.ResponseWriter) {
	m := make(map[string]string)
	m[model.STATUS] = model.STATUS_OK
	RenderJson(w, m)
}

func RenderJson(w http.ResponseWriter, o interface{}) {
	if b, err := json.Marshal(o); err != nil {
		w.Write([]byte(""))
	} else {
		w.Write(b)
	}
}

func BindJson(data io.Reader, dest interface{}) error {
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return errors.New("BindJSON not a pointer")
	}

	decoder := json.NewDecoder(data)

	if err := decoder.Decode(dest); err != nil {
		l4g.Debug(err)
		return err
	}

	return nil
}
