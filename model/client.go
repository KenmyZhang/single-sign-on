package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	HEADER_REQUEST_ID         = "X-Request-ID"
	HEADER_VERSION_ID         = "X-Version-ID"
	HEADER_ETAG_SERVER        = "ETag"
	HEADER_ETAG_CLIENT        = "If-None-Match"
	HEADER_FORWARDED          = "X-Forwarded-For"
	HEADER_REAL_IP            = "X-Real-IP"
	HEADER_FORWARDED_PROTO    = "X-Forwarded-Proto"
	HEADER_TOKEN              = "Token"
	HEADER_BEARER             = "BEARER"
	HEADER_AUTH               = "Authorization"
	HEADER_REQUESTED_WITH     = "X-Requested-With"
	HEADER_REQUESTED_WITH_XML = "XMLHttpRequest"
	STATUS                    = "status"
	STATUS_OK                 = "OK"
	LAST_PICTURE_UPDATE		  = "Last_picture_update"

	CLIENT_DIR = "webapp/dist"

	API_URL_SUFFIX = "/sso"
)

type Response struct {
	StatusCode    int
	Error         *AppError
	RequestId     string
	Etag          string
	ServerVersion string
}

type Client struct {
	Url        string       // The location of the server, for example  "http://localhost:8065"
	ApiUrl     string       // The api location of the server, for example "http://localhost:8065/api/v4"
	HttpClient *http.Client // The http client
	AuthToken  string
	AuthType   string
}

func NewAPIClient(url string) *Client {
	return &Client{url, url + API_URL_SUFFIX, &http.Client{}, "", ""}
}

func BuildErrorResponse(r *http.Response, err *AppError) *Response {
	if r == nil {
		return &Response{StatusCode: 0, Error: err}
	} else {
		return &Response{StatusCode: r.StatusCode, Error: err}
	}
}

func BuildResponse(r *http.Response) *Response {
	return &Response{
		StatusCode:    r.StatusCode,
		RequestId:     r.Header.Get(HEADER_REQUEST_ID),
		Etag:          r.Header.Get(HEADER_ETAG_SERVER),
		ServerVersion: r.Header.Get(HEADER_VERSION_ID),
	}
}

func closeBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}

func (c *Client) SetOAuthToken(token string) {
	c.AuthToken = token
	c.AuthType = HEADER_TOKEN
}

func (c *Client) ClearOAuthToken() {
	c.AuthToken = ""
	c.AuthType = HEADER_BEARER
}

func (c *Client) GetUsersRoute() string {
	return fmt.Sprintf("/users")
}

func (c *Client) GetUserRoute(userId string) string {
	return fmt.Sprintf(c.GetUsersRoute()+"/%v", userId)
}

func (c *Client) DoApiGet(url string, etag string) (*http.Response, *AppError) {
	return c.DoApiRequest(http.MethodGet, c.ApiUrl+url, "", etag)
}

func (c *Client) DoApiPost(url string, data string) (*http.Response, *AppError) {
	return c.DoApiRequest(http.MethodPost, c.ApiUrl+url, data, "")
}

func (c *Client) DoApiPut(url string, data string) (*http.Response, *AppError) {
	return c.DoApiRequest(http.MethodPut, c.ApiUrl+url, data, "")
}

func (c *Client) DoApiDelete(url string) (*http.Response, *AppError) {
	return c.DoApiRequest(http.MethodDelete, c.ApiUrl+url, "", "")
}

func (c *Client) DoApiRequest(method, url, data, etag string) (*http.Response, *AppError) {
	rq, _ := http.NewRequest(method, url, strings.NewReader(data))
	rq.Close = true

	if len(etag) > 0 {
		rq.Header.Set(HEADER_ETAG_CLIENT, etag)
	}

	if len(c.AuthToken) > 0 {
		rq.Header.Set(HEADER_AUTH, c.AuthType+" "+c.AuthToken)
	}

	if rp, err := c.HttpClient.Do(rq); err != nil || rp == nil {
		return nil, NewLocAppError(url, "model.client.connecting.app_error", nil, err.Error())
	} else if rp.StatusCode == 304 {
		return rp, nil
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return rp, AppErrorFromJson(rp.Body)
	} else {
		return rp, nil
	}
}

func (c *Client) Login(loginId string, password string) (*User, *Response) {
	m := make(map[string]string)
	m["login_id"] = loginId
	m["password"] = password
	return c.login(m)
}

func (c *Client) login(m map[string]string) (*User, *Response) {
	if r, err := c.DoApiPost("/users/login", MapToJson(m)); err != nil {
		return nil, BuildErrorResponse(r, err)
	} else {
		c.AuthToken = r.Header.Get(HEADER_TOKEN)
		c.AuthType = HEADER_BEARER
		defer closeBody(r)
		return DecodeUserFromJson(r.Body), BuildResponse(r)
	}
}

// CreateUser creates a user in the system based on the provided user struct.
func (c *Client) CreateUser(user *User) (*User, *Response) {
	if r, err := c.DoApiPost(c.GetUsersRoute(), user.ToJson()); err != nil {
		return nil, BuildErrorResponse(r, err)
	} else {
		defer closeBody(r)
		return DecodeUserFromJson(r.Body), BuildResponse(r)
	}
}