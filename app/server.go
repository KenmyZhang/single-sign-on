package app

import (
	"net"
	"net/http"
	"strings"
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/sqlStore"
	"github.com/KenmyZhang/single-sign-on/utils"
)

type Server struct {
	SqlStore       sqlStore.Store
	Router         *mux.Router
	GracefulServer *graceful.Server
}

var allowedMethods []string = []string{
	"POST",
	"GET",
	"OPTIONS",
	"PUT",
	"PATCH",
	"DELETE",
}

type RecoveryLogger struct {
}

func (rl *RecoveryLogger) Println(i ...interface{}) {
	l4g.Error("Please check the std error output for the stack trace")
	l4g.Error(i)
}

type CorsWrapper struct {
	router *mux.Router
}

func (cw *CorsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(*utils.Cfg.ServiceSettings.AllowCorsFrom) > 0 {
		origin := r.Header.Get("Origin")
		if *utils.Cfg.ServiceSettings.AllowCorsFrom == "*" ||
			strings.Contains(*utils.Cfg.ServiceSettings.AllowCorsFrom, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			}
		}
	}

	if r.Method == "OPTIONS" {
		return
	}

	cw.router.ServeHTTP(w, r)
}

const TIME_TO_WAIT_FOR_CONNECTIONS_TO_CLOSE_ON_SERVER_SHUTDOWN = time.Second

var Srv *Server

func NewServer() {
	l4g.Info(utils.T("api.server.new_server.init.info"))

	Srv = &Server{}
}

func InitStores() {
	Srv.SqlStore = sqlStore.NewLayeredStore()
}

type VaryBy struct{}

func (m *VaryBy) Key(r *http.Request) string {
	return utils.GetIpAddress(r)
}

func redirectHTTPToHTTPS(w http.ResponseWriter, r *http.Request) {
	if r.Host == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	url := r.URL
	url.Host = r.Host
	url.Scheme = "https"
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func StartServer() {
	l4g.Info(utils.T("api.server.start_server.starting.info"))

	var handler http.Handler = &CorsWrapper{Srv.Router}

	Srv.GracefulServer = &graceful.Server{
		Timeout: TIME_TO_WAIT_FOR_CONNECTIONS_TO_CLOSE_ON_SERVER_SHUTDOWN,
		Server: &http.Server{
			Addr:         utils.Cfg.ServiceSettings.ListenAddress,
			Handler:      handlers.RecoveryHandler(handlers.RecoveryLogger(&RecoveryLogger{}), handlers.PrintRecoveryStack(true))(handler),
			ReadTimeout:  time.Duration(*utils.Cfg.ServiceSettings.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(*utils.Cfg.ServiceSettings.WriteTimeout) * time.Second,
		},
	}
	l4g.Info(utils.T("api.server.start_server.listening.info"), utils.Cfg.ServiceSettings.ListenAddress)

	if *utils.Cfg.ServiceSettings.Forward80To443 {
		go func() {
			listener, err := net.Listen("tcp", ":80")
			if err != nil {
				l4g.Error("Unable to setup forwarding")
				return
			}
			defer listener.Close()

			http.Serve(listener, http.HandlerFunc(redirectHTTPToHTTPS))
		}()
	}

	go func() {
		var err error
		if *utils.Cfg.ServiceSettings.ConnectionSecurity == model.CONN_SECURITY_TLS {
			err = Srv.GracefulServer.ListenAndServeTLS(*utils.Cfg.ServiceSettings.TLSCertFile, *utils.Cfg.ServiceSettings.TLSKeyFile)
		} else {
			err = Srv.GracefulServer.ListenAndServe()
		}
		if err != nil {
			l4g.Critical(utils.T("api.server.start_server.starting.critical"), err)
			time.Sleep(time.Second)
		}
	}()
}

func StopServer() {
	l4g.Info(utils.T("api.server.stop_server.stopping.info"))
	Srv.GracefulServer.Stop(TIME_TO_WAIT_FOR_CONNECTIONS_TO_CLOSE_ON_SERVER_SHUTDOWN)
	Srv.SqlStore.Close()
	l4g.Info(utils.T("api.server.stop_server.stopped.info"))
}
