package api

import (
	"github.com/KenmyZhang/single-sign-on/app"
	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
	"net/http"
	"strings"
)

func InitOauth() {
	BaseRoutes.ApiRoot.Handle("/oauth/{service:[A-Za-z0-9]+}/login", ApiHandler(loginWithOauth)).Methods("GET")
	BaseRoutes.ApiRoot.Handle("/oauth/{service:[A-Za-z0-9]+}/signup", ApiHandler(signupWithOauth)).Methods("POST")
	BaseRoutes.ApiRoot.Handle("/signup/{service:[A-Za-z0-9]+}/complete", ApiHandler(completeOAuth)).Methods("GET")
	BaseRoutes.ApiRoot.Handle("/oauth/{service:[A-Za-z0-9]+}/mobile", ApiHandler(completeOAuthMobile)).Methods("POST")
}

func loginWithOauth(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireService()
	if c.Err != nil {
		return
	}

	if authUrl, err := app.GetOAuthLoginEndpoint(w, r, c.Params.Service, model.OAUTH_ACTION_LOGIN); err != nil {
		c.Err = err
		return
	} else {
		http.Redirect(w, r, authUrl, http.StatusFound)
	}
}

func signupWithOauth(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireService()
	if c.Err != nil {
		return
	}

	state := r.URL.Query().Get("state")

	var props map[string]string
	BindJson(r.Body, &props)

	loginId := props["username"]
	password := props["password"]

	if user, err := app.SignupWithOauth(w, r, c.Params.Service, state, loginId, password); err != nil {
		c.Err = err
		return
	} else {
		err := app.DoLogin(w, r, user, "")
		if err != nil {
			err.Translate(c.T)
			c.Err = err
			return
		}

		user.Sanitize(map[string]bool{})
		RenderJson(w, user)
	}
}

func completeOAuth(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireService()
	if c.Err != nil {
		return
	}

	service := c.Params.Service

	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		http.Redirect(w, r, c.GetSiteURLHeader()+"/error?type=oauth_missing_code&service="+strings.Title(service), http.StatusTemporaryRedirect)
		return
	}

	state := r.URL.Query().Get("state")

	uri := c.GetSiteURLHeader() + "/sso/signup/" + service + "/complete"

	redirectUrlInState, body, err := app.AuthorizeOAuthUser(w, r, service, code, state, uri)

	if err != nil {
		http.Redirect(w, r, c.GetSiteURLHeader()+"/error?message="+err.Message, http.StatusTemporaryRedirect)
		return
	}

	exist, user, err := app.CompleteOAuth(service, body)
	if err != nil {
		http.Redirect(w, r, c.GetSiteURLHeader()+"/error?message="+err.Message, http.StatusTemporaryRedirect)
		return
	}

	var redirectUrl string
	if exist {
		err := app.DoLogin(w, r, user, "")
		if err != nil {
			err.Translate(c.T)
			c.Err = err
			return
		}
		redirectUrl = utils.GetSiteURL() + "/sso"
	} else {
		if state, err := app.SetCookieAndToken(w, r, user); err != nil {
			c.Err = err
			return
		} else {
			redirectUrl = utils.GetSiteURL() + "/sso/wechat?state=" + state
		}
	}
	if redirectUrlInState != "" {
		//reserved
	}
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func completeOAuthMobile(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)
	if len(props["access_token"]) > 120 {
		c.SetInvalidParam("access_token")
		return		
	}

	if len(props["open_id"]) > 120 {
		c.SetInvalidParam("open_id")
		return		
	}		

	accessToken := props["access_token"]
	openId := props["open_id"]

	sso := utils.Cfg.GetSSOService(model.SERVICE_WEIXIN)
	if sso == nil {
		err := model.NewLocAppError("completeOAuthMobile", "api.user.authorize_oauth_user.unsupported.app_error", nil, "service=" + model.SERVICE_WEIXIN)
		c.Err = err
		return
	}

    userApiEndpoint := sso.UserApiEndpoint + "?access_token=" + accessToken + "&openid=" + openId
    req, _ := http.NewRequest("GET", userApiEndpoint, strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var resp *http.Response
	var userApiErr error
	if resp, userApiErr = utils.HttpClient(true).Do(req); userApiErr != nil {
		err := model.NewLocAppError("completeOAuthMobile", "api.oauth.complete_oauth_mobile.request.app_error",
			map[string]interface{}{"Service": model.SERVICE_WEIXIN}, userApiErr.Error())
		c.Err = err
		return		
	}

	user, err := app.CompleteOAuthMobile(model.SERVICE_WEIXIN, resp.Body, props)
	if err != nil {
		c.Err = err
		return
	}
	
	err = app.DoLogin(w, r, user, "")
	if err != nil {
		c.Err = err
		return
	}
	       
	if user != nil && user.HeadImgUrl != "" {
        if err := app.SetOAuthProfileImage(user); err != nil {
            c.Err = err
            return
        }
    }


	user.Sanitize(map[string]bool{})
	w.Write([]byte(user.ToJson()))	
	return
}
