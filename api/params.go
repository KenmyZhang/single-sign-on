package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	PAGE_DEFAULT     = 0
	PER_PAGE_DEFAULT = 10
	PER_PAGE_MAXIMUM = 200
)

type ApiParams struct {
	UserId               string
	Email                string
	Username             string
	Region               string
	Term                 string
	Service              string
	Mobile               string
	Page                 int
	PerPage              int
}

func ApiParamsFromRequest(r *http.Request) *ApiParams {
	params := &ApiParams{}

	props := mux.Vars(r)

	if val, ok := props["user_id"]; ok {
		params.UserId = val
	}

	if val, ok := props["email"]; ok {
		params.Email = val
	}

	if val, ok := props["username"]; ok {
		params.Username = val
	}

	if val, ok := props["service"]; ok {
		params.Service = val
	}

	params.Region = r.URL.Query().Get("region")
	params.Term = r.URL.Query().Get("term")
	params.Mobile = r.URL.Query().Get("mobile")

	if val, err := strconv.Atoi(r.URL.Query().Get("page")); err != nil || val < 0 {
		params.Page = PAGE_DEFAULT
	} else {
		params.Page = val
	}

	if val, err := strconv.Atoi(r.URL.Query().Get("per_page")); err != nil || val < 0 {
		params.PerPage = PER_PAGE_DEFAULT
	} else if val > PER_PAGE_MAXIMUM {
		params.PerPage = PER_PAGE_MAXIMUM
	} else {
		params.PerPage = val
	}

	return params
}
