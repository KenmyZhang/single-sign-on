package app

import (
	"bytes"
	b64 "encoding/base64"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/KenmyZhang/single-sign-on/einterfaces"
	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
	"github.com/KenmyZhang/single-sign-on/sqlStore"
)

const (
	OAUTH_COOKIE_MAX_AGE_SECONDS = 30 * 60 // 30 minutes
	COOKIE_OAUTH                 = "MMOAUTH"
)

func GetOAuthLoginEndpoint(w http.ResponseWriter, r *http.Request, service, action string) (string, *model.AppError) {
	stateProps := map[string]string{}
	stateProps["action"] = action

	if authUrl, err := GetAuthorizationCode(w, r, service, stateProps); err != nil {
		return "", err
	} else {
		return authUrl, nil
	}
}

func GetAuthorizationCode(w http.ResponseWriter, r *http.Request, service string, props map[string]string) (string, *model.AppError) {
	sso := utils.Cfg.GetSSOService(service)
	if sso != nil && !sso.Enable {
		return "", model.NewAppError("GetAuthorizationCode", "app.oauth.get_authorization_code.unsupported.app_error",
			nil, "service="+service, http.StatusNotImplemented)
	}

	secure := false
	if GetProtocol(r) == "https" {
		secure = true
	}

	cookieValue := model.NewId()
	expiresAt := time.Unix(model.GetMillis()/1000+int64(OAUTH_COOKIE_MAX_AGE_SECONDS), 0)
	oauthCookie := &http.Cookie{
		Name:     COOKIE_OAUTH,
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   OAUTH_COOKIE_MAX_AGE_SECONDS,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
	}

	http.SetCookie(w, oauthCookie)

	clientId := sso.Id
	endpoint := sso.AuthEndpoint
	scope := sso.Scope
	propsForToken := model.StringInterface{}

	tokenExtra := generateOAuthStateTokenExtra(props["action"], cookieValue)
	stateToken, err := CreateOAuthStateToken(tokenExtra, propsForToken)
	if err != nil {
		return "", err
	}

	props["token"] = stateToken.Token
	state := b64.StdEncoding.EncodeToString([]byte(utils.MapToJson(props)))

	redirectUri := utils.GetSiteURL() + "/sso/signup/" + service + "/complete"
	authUrl := endpoint + "?appid=" + clientId + "&redirect_uri=" + url.QueryEscape(redirectUri) + "&response_type=code" + "&scope=" + utils.UrlEncode(scope) + "&state=" + url.QueryEscape(state) + "#wechat_redirect"

	return authUrl, nil
}

func CreateOAuthStateToken(extra string, props model.StringInterface) (*model.Token, *model.AppError) {
	token := model.NewToken(model.TOKEN_TYPE_OAUTH, extra, props)
	if result := <-Srv.SqlStore.Token().Save(token); result.Err != nil {
		return nil, result.Err
	}

	return token, nil
}

func generateOAuthStateTokenExtra(action, cookie string) string {
	return action + ":" + cookie
}

func AuthorizeOAuthUser(w http.ResponseWriter, r *http.Request, service, code, state, redirectUri string) (string , io.ReadCloser, *model.AppError) {
	sso := utils.Cfg.GetSSOService(service)
	if sso == nil || !sso.Enable {
		return "", nil, model.NewAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.unsupported.app_error", nil, "service="+service, http.StatusNotImplemented)
	}

	stateProps := map[string]string{}
	stateStr := ""
	if b, err := b64.StdEncoding.DecodeString(state); err != nil {
		return "", nil, model.NewLocAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.invalid_state.app_error", nil, err.Error())
	} else {
		stateStr = string(b)
	}

	stateProps = model.MapFromJson(strings.NewReader(stateStr))

	expectedToken, err := GetOAuthStateToken(stateProps["token"])
	if err != nil {
		return "", nil, err
	}

	stateAction := stateProps["action"]

	cookieValue := ""
	if cookie, err := r.Cookie(COOKIE_OAUTH); err != nil {
		return "", nil, model.NewAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.invalid_state.app_error", nil, "", http.StatusBadRequest)
	} else {
		cookieValue = cookie.Value
	}

	expectedTokenExtra := generateOAuthStateTokenExtra(stateAction, cookieValue)
	if expectedTokenExtra != expectedToken.Extra {
		return "", nil, model.NewAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.invalid_state.app_error", nil, "", http.StatusBadRequest)
	}
	DeleteToken(expectedToken)

	cookie := &http.Cookie{
		Name:     COOKIE_OAUTH,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	p := url.Values{}
	p.Set("client_id", sso.Id)
	p.Set("client_secret", sso.Secret)
	p.Set("code", code)
	p.Set("grant_type", model.ACCESS_TOKEN_GRANT_TYPE)
	p.Set("redirect_uri", redirectUri)

	var req *http.Request
	tokenUrl := sso.TokenEndpoint + "?appid=" + sso.Id + "&secret=" + sso.Secret + "&code=" + code + "&grant_type=authorization_code"
	req, _ = http.NewRequest("POST", tokenUrl, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	var ar *model.AccessResponse
	var bodyBytes []byte
	if resp, err := utils.HttpClient(true).Do(req); err != nil {
		return "", nil, model.NewAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.token_failed.app_error", nil, err.Error(), http.StatusInternalServerError)
	} else {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		ar = model.AccessResponseFromJson(resp.Body)
		defer CloseBody(resp)
		if ar == nil {
			return "", nil, model.NewAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.bad_response.app_error", nil, "response_body="+string(bodyBytes), http.StatusInternalServerError)
		}
	}

	if len(ar.AccessToken) == 0 {
		return "", nil, model.NewAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.missing.app_error", nil, "response_body="+string(bodyBytes), http.StatusInternalServerError)
	}

	p = url.Values{}
	p.Set("access_token", ar.AccessToken)

	UserApiEndpoint := sso.UserApiEndpoint + "?access_token=" + ar.AccessToken + "&openid=" + ar.OpenId
	req, _ = http.NewRequest("GET", UserApiEndpoint, strings.NewReader(""))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+ar.AccessToken)

	if resp, err := utils.HttpClient(true).Do(req); err != nil {
		return "", nil, model.NewAppError("AuthorizeOAuthUser", "api.oauth.authorize_oauth_user.service.app_error", map[string]interface{}{"Service": service}, err.Error(), http.StatusInternalServerError)
	} else {
		return stateProps["redirect_uri"], resp.Body, nil
	}
}

func CloseBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}

func DeleteToken(token *model.Token) *model.AppError {
	if result := <-Srv.SqlStore.Token().Delete(token.Token); result.Err != nil {
		return result.Err
	}

	return nil
}

func GetOAuthStateToken(token string) (*model.Token, *model.AppError) {
	if result := <-Srv.SqlStore.Token().GetByToken(token); result.Err != nil {
		return nil, model.NewAppError("GetOAuthStateToken", "api.oauth.invalid_state_token.app_error", nil, result.Err.Error(), http.StatusBadRequest)
	} else {
		token := result.Data.(*model.Token)
		if token.Type != model.TOKEN_TYPE_OAUTH {
			return nil, model.NewAppError("GetOAuthStateToken", "api.oauth.invalid_state_token.app_error", nil, "", http.StatusBadRequest)
		}
		return token, nil
	}
}

func CompleteOAuth(service string, userData io.ReadCloser) (bool, *model.User, *model.AppError) {
	defer func() {
		ioutil.ReadAll(userData)
		userData.Close()
	}()

	buf := bytes.Buffer{}
	buf.ReadFrom(userData)

	authData := ""
	provider := einterfaces.GetOauthProvider(service)
	if provider == nil {
		return false, nil, model.NewAppError("LoginByOAuth", "api.user.login_by_oauth.not_available.app_error",
			map[string]interface{}{"Service": strings.Title(service)}, "", http.StatusNotImplemented)
	} else {
		authData = provider.GetAuthDataFromJson(bytes.NewReader(buf.Bytes()))
	}

	if len(authData) == 0 {
		return false, nil, model.NewAppError("LoginByOAuth", "api.user.login_by_oauth.parse.app_error",
			map[string]interface{}{"Service": service},
			"invalid credential, access_token is invalid or not latest", http.StatusBadRequest)
	}

	user, err := GetUserByAuth(&authData, service)
	if err != nil {
		if err.Id == "store.sql_user.get_by_auth.missing_account.app_error" {
			return CreateOAuthUser(service, bytes.NewReader(buf.Bytes()))
		}
		return false, nil, err
	}

	return true, user, nil
}

func CreateOAuthUser(service string, userData io.Reader) (bool, *model.User, *model.AppError) {
	var user *model.User
	provider := einterfaces.GetOauthProvider(service)
	if provider == nil {
		return false, nil, model.NewAppError("CreateOAuthUser", "api.user.create_oauth_user.not_available.app_error", map[string]interface{}{"Service": strings.Title(service)}, "", http.StatusNotImplemented)
	} else {
		user = provider.GetUserFromJson(userData)
	}
	if user == nil {
		return false, nil, model.NewAppError("CreateOAuthUser", "api.user.create_oauth_user.create.app_error", map[string]interface{}{"Service": service}, "", http.StatusInternalServerError)
	}

	found := true
	count := 0
	for found {
		user.Username = strings.ToLower(user.Username)
		if found = IsUsernameTaken(user.Username); found {
			user.Username = user.Username + strconv.Itoa(count)
			count += 1
		}
	}

	return false, user, nil
}

func SetOAuthProfileImage(user *model.User) *model.AppError {
	var body []byte
	var err error
	var resp *http.Response
	if resp, err = http.Get(user.HeadImgUrl); err != nil {
		return model.NewLocAppError("SetOAuthProfileImage", "api.user.set_oauth_profile_image.http_get.app_error", nil, err.Error())
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return model.NewLocAppError("SetOAuthProfileImage", "api.user.set_oauth_profile_image.read_all.app_error", nil, err.Error())
	}

	input := bytes.NewReader(body)

	if config, err := jpeg.DecodeConfig(input); err != nil {
		return model.NewLocAppError("SetOAuthProfileImage", "api.user.set_oauth_profile_image.decode_config.app_error", nil, err.Error())
	} else if config.Width*config.Height > MaxImageSize {
		return model.NewLocAppError("SetOAuthProfileImage", "api.user.set_oauth_profile_image.size.app_error", nil, "image is too large")
	}

	input = bytes.NewReader(body)
	orientation, _ := getImageOrientation(input)

	input = bytes.NewReader(body)
	img, err := jpeg.Decode(input)
	if err != nil {
		return model.NewLocAppError("SetOAuthProfileImage", "api.user.set_oauth_profile_image.decode.app_error", nil, err.Error())
	}

	img = makeImageUpright(img, orientation)

	profileWidthAndHeight := 128
	img = imaging.Fill(img, profileWidthAndHeight, profileWidthAndHeight, imaging.Center, imaging.Lanczos)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{100})
	if err != nil {
		return model.NewLocAppError("SetOAuthProfileImage", "api.user.set_oauth_profile_image.encode.app_error", nil, err.Error())
	}

	path := "users/" + user.Id + "/profile.png"
	if err := writeFileLocally(buf.Bytes(), *utils.Cfg.FileSettings.Directory + path); err != nil {
		return err
	}

	<-Srv.SqlStore.User().UpdateLastPictureUpdate(user.Id)

	return nil
}

func SignupWithOauth(w http.ResponseWriter, r *http.Request, service, state, loginId, password string) (*model.User, *model.AppError) {
	stateProps := map[string]string{}
	stateStr := ""

	if b, err := b64.StdEncoding.DecodeString(state); err != nil {
		return nil, model.NewLocAppError("SignupWithOauth", "api.oauth.signup_with_oauth.invalid_state.app_error", nil, err.Error())
	} else {
		stateStr = string(b)
	}

	stateProps = model.MapFromJson(strings.NewReader(stateStr))
	expectedToken, err := GetOAuthStateToken(stateProps["token"])
	if err != nil {
		return nil, err
	}

	cookieValue := ""
	if cookie, err := r.Cookie(COOKIE_OAUTH); err != nil {
		return nil, model.NewAppError("SignupWithOauth", "api.user.signup_with_oauth.invalid_cookie.app_error", nil, "", http.StatusBadRequest)
	} else {
		cookieValue = cookie.Value
	}

	expectedTokenExtra := generateOAuthStateTokenExtra(model.OAUTH_ACTION_SIGNUP, cookieValue)
	if expectedTokenExtra != expectedToken.Extra {
		return nil, model.NewAppError("SignupWithOauth", "api.user.signup_with_oauth.invalid_token_extra.app_error", nil, "", http.StatusBadRequest)
	}
	propsForToken := expectedToken.Props
	user :=  &model.User{}

	user.Username    = propsForToken["username"].(string)
	user.AuthService = propsForToken["authService"].(string)
	user.Nickname    = propsForToken["nickname"].(string)
	user.HeadImgUrl  = propsForToken["headImgUrl"].(string)
	authData        := propsForToken["authData"].(string)
	if authData != "" {
		user.AuthData = &authData
	}

	DeleteToken(expectedToken)

	var ruser *model.User
	if loginId != "" {
		ruser, err = AuthenticateUserForLogin("", loginId, password)
		if err != nil {
			return nil, err
		}
		ruser.AuthData = user.AuthData
		ruser.AuthService = user.AuthService
		ruser.HeadImgUrl = user.HeadImgUrl
		if result := <-Srv.SqlStore.User().Update(ruser, true); result.Err != nil {
			return nil, result.Err
		} else {
			ruser = result.Data.([2]*model.User)[0]
		}
	} else {
		if ruser, err = CreateUser(user); err != nil {
			return nil, err
		}
	}

	if err := SetOAuthProfileImage(user); err != nil {
		return nil, err
	}

	return ruser, nil
}

func SetCookieAndToken(w http.ResponseWriter, r *http.Request, user *model.User) (string, *model.AppError) {
	secure := false
	if GetProtocol(r) == "https" {
		secure = true
	}

	cookieValue := model.NewId()
	expiresAt := time.Unix(model.GetMillis()/1000+int64(OAUTH_COOKIE_MAX_AGE_SECONDS), 0)
	oauthCookie := &http.Cookie{
		Name:     COOKIE_OAUTH,
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   OAUTH_COOKIE_MAX_AGE_SECONDS,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
	}

	http.SetCookie(w, oauthCookie)

	tokenExtra := generateOAuthStateTokenExtra(model.OAUTH_ACTION_SIGNUP, cookieValue)
	propsForToken := model.StringInterface{}
	propsForToken["username"]    = user.Username
	propsForToken["authData"]    = *(user.AuthData) 
	propsForToken["authService"] = user.AuthService  
	propsForToken["nickname"]    = user.Nickname     
	propsForToken["headImgUrl"]  = user.HeadImgUrl  

	stateToken, err := CreateOAuthStateToken(tokenExtra, propsForToken)
	if err != nil {
		return "", err
	}
	props := map[string]string{}
	props["token"] = stateToken.Token
	state := b64.StdEncoding.EncodeToString([]byte(utils.MapToJson(props)))
	return state, nil
}

func CompleteOAuthMobile(service string, userData io.ReadCloser, props map[string]string) (*model.User, *model.AppError) {
	defer func() {
		ioutil.ReadAll(userData)
		userData.Close()
	}()

	buf := bytes.Buffer{}
	buf.ReadFrom(userData)

	authData := ""
	provider := einterfaces.GetOauthProvider(service)
	if provider == nil {
		return nil, model.NewAppError("CompleteOAuthMobile", "api.user.complete_oauth_mobile.not_available.app_error",
			map[string]interface{}{"Service": strings.Title(service)}, "", http.StatusNotImplemented)
	} else {
		authData = provider.GetAuthDataFromJson(bytes.NewReader(buf.Bytes()))
	}

	if len(authData) == 0 {
		return nil, model.NewAppError("CompleteOAuthMobile", "api.user.complete_oauth_mobile.parse.app_error",
			map[string]interface{}{"Service": service}, 
			"invalid credential, access_token is invalid or not latest", http.StatusBadRequest)
	}

	user, err := GetUserByAuth(&authData, service)
	if err != nil {
		if err.Id == sqlStore.MISSING_AUTH_ACCOUNT_ERROR {
			return CreateOAuthMobileUser(service, bytes.NewReader(buf.Bytes()))
		}
		return nil, err
	}

	return user, nil
}

func CreateOAuthMobileUser(service string, userData io.Reader) (*model.User, *model.AppError) {
	var user *model.User
	provider := einterfaces.GetOauthProvider(service)
	if provider == nil {
		return nil, model.NewLocAppError("CreateOAuthMobileUser", "api.user.create_oauth_mobile_user.not_available.app_error", map[string]interface{}{"Service": strings.Title(service)}, "")
	} else {
		user = provider.GetUserFromJson(userData)
	}

	if user == nil {
		return nil, model.NewLocAppError("CreateOAuthUser", "api.user.create_oauth_mobile_user.create.app_error", map[string]interface{}{"Service": service}, "")
	}

	found := true
	count := 0
	for found {
		if found = IsUsernameTaken(user.Username); found {
			user.Username = user.Username + strconv.Itoa(count)
			count += 1
		}
	}

	user.EmailVerified = false

	ruser, err := CreateUser(user)
	if err != nil {
		return nil, err
	}
	
	return ruser, nil
}
