package api

import (
	"net/http"
	"regexp"		
	"strconv"
	"fmt"
	l4g "github.com/alecthomas/log4go"

	"github.com/KenmyZhang/single-sign-on/app"
	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
	"github.com/KenmyZhang/single-sign-on/sqlStore"
)

func InitUser() {
	l4g.Debug(utils.T("api.user.init.debug"))
	BaseRoutes.User.Handle("", ApiCustomClaimsRequired(getUser)).Methods("GET")
	BaseRoutes.User.Handle("/image", ApiHandler(getProfileImage)).Methods("GET")
	BaseRoutes.User.Handle("/image", ApiCustomClaimsRequired(setProfileImage)).Methods("POST")		
	BaseRoutes.Users.Handle("/login", ApiHandler(login)).Methods("POST")
	BaseRoutes.Users.Handle("/logout", ApiHandler(logout)).Methods("POST")
	BaseRoutes.Users.Handle("/sendsms", ApiHandler(sendSmsCode)).Methods("POST")
	BaseRoutes.Users.Handle("/phone/signup", ApiHandler(signupByMobile)).Methods("POST")
	BaseRoutes.Users.Handle("/phone/login", ApiHandler(loginByMobile)).Methods("POST")
	BaseRoutes.Users.Handle("/phone/exist", ApiHandler(isMobileExist)).Methods("POST")
	BaseRoutes.Users.Handle("/phone/reset", ApiHandler(resetPasswordByMobile)).Methods("POST")
	BaseRoutes.Users.Handle("/email/verify/code/send", ApiHandler(sendVerificationCodeEmail)).Methods("POST")
	BaseRoutes.Users.Handle("/email/signup", ApiHandler(signupByEmail)).Methods("POST")
	BaseRoutes.Users.Handle("/email/exist", ApiHandler(isEmailExist)).Methods("POST")
	BaseRoutes.Users.Handle("/email/reset", ApiHandler(resetPasswordByEmail)).Methods("POST")
	BaseRoutes.Users.Handle("/search", ApiCustomClaimsRequired(searchUsers)).Methods("POST")
}

func getUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	var user *model.User
	var err *model.AppError

	if user, err = app.GetUser(c.Params.UserId); err != nil {
		c.Err = err
		return
	}

	etag := user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress)

	if HandleEtag(etag, "Get User", w, r) {
		return
	} else {
		if c.CustomClaims.UserId == user.Id {
			user.Sanitize(map[string]bool{})
		} else {
			app.SanitizeProfile(user, c.IsSystemAdmin())
		}
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
		RenderJson(w, user)
		return
	}
}


func getProfileImage(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if user, err := app.GetUser(c.Params.UserId); err != nil {
		c.Err = err
		return
	} else {
		etag := strconv.FormatInt(user.LastPictureUpdate, 10)
		if HandleEtag(etag, "Get Profile Image", w, r) {
			return
		}

		var img []byte
		img, readFailed, err := app.GetProfileImage(user)
		if err != nil {
			c.Err = err
			return
		}

		if readFailed {
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v, public", 5*60)) // 5 mins
		} else {
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v, public", 24*60*60)) // 24 hrs
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
		w.Write(img)
	}
}

func setProfileImage(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if r.ContentLength > *utils.Cfg.FileSettings.MaxFileSize {
		c.Err = model.NewLocAppError("uploadProfileImage", "api.user.upload_profile_user.too_large.app_error", nil, "")
		c.Err.StatusCode = http.StatusRequestEntityTooLarge
		return
	}

	if err := r.ParseMultipartForm(*utils.Cfg.FileSettings.MaxFileSize); err != nil {
		c.Err = model.NewLocAppError("uploadProfileImage", "api.user.upload_profile_user.parse.app_error", nil, "")
		return
	}

	m := r.MultipartForm

	imageArray, ok := m.File["image"]
	if !ok {
		c.Err = model.NewLocAppError("uploadProfileImage", "api.user.upload_profile_user.no_file.app_error", nil, "")
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	if len(imageArray) <= 0 {
		c.Err = model.NewLocAppError("uploadProfileImage", "api.user.upload_profile_user.array.app_error", nil, "")
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	imageData := imageArray[0]

	if err := app.SetProfileImage(w, c.CustomClaims.UserId, imageData); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}

func login(c *Context, w http.ResponseWriter, r *http.Request) {
	var props map[string]string
	BindJson(r.Body, &props)

	id := props["id"]
	loginId := props["login_id"]
	password := props["password"]
	deviceId := props["device_id"]

	if loginId == "" {
		c.SetInvalidParam("loginId")
		return
	}
	user, err := app.AuthenticateUserForLogin(id, loginId, password)
	if err != nil {
		c.Err = err
		return
	}

	err = app.DoLogin(w, r, user, deviceId)
	if err != nil {
		c.Err = err
		return
	}

	user.Sanitize(map[string]bool{})

	RenderJson(w, user)
}

func logout(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RemoveCustomClaimsCookie(w, r)
	ReturnStatusOK(w)
}

func sendSmsCode(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)
	mobile  := props["mobile"]	
	if len(mobile) == 0 {
		c.SetInvalidParam("mobile")
		return
	}

	if err := app.SendSmsCode(mobile, w, r); err != nil {
		c.Err = err
		return
	}
	ReturnStatusOK(w)
}


func signupByMobile(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	mobile  := props["mobile"]	
	if len(mobile) == 0 {
		c.SetInvalidParam("mobile")
		return
	}

	username := props["username"]	
	if username == "" {
		c.SetInvalidParam("username")
		return		
	}

	password := props["password"]
	if password == "" {
		c.SetInvalidParam("password")
		return	
	}

	nickname := props["nickname"]
	if nickname == "" {
		nickname = username	
	}

	if exist := app.IsMobileExist(mobile); exist == true {
		c.SetInvalidParam("mobile is already exist")
		return
	}

	verificationCode := props["verification_code"]
	if err := app.VerifiedCode(r, "", mobile, verificationCode); err != nil {
		c.Err = err
		return
	}

	user := &model.User{Username:username, Nickname:nickname, Mobile:mobile, Password:password, AllowMarketing: true}
	_, err := app.CreateUserFromEmailOrMobile(user)
	if err != nil {
		c.Err = err
		return
	}

	ruser, err := app.AuthenticateUserForLogin(user.Id, user.Username, password)
	if err != nil {
		c.Err = err
		return
	}

	err = app.DoLogin(w, r, ruser, "")
	if err != nil {
		c.Err = err
		return
	}

	ruser.Sanitize(map[string]bool{})
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ruser.ToJson()))
}

func loginByMobile(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)
	mobile  := props["mobile"]
	verificationCode := props["verification_code"]
	password := props["password"]
	deviceId := props["device_id"]

	if mobile == "" {
		c.SetInvalidParam("phone num")
		return
	}

	if verificationCode == "" && password == "" {
		c.SetInvalidParam("password  or verification code ")
		return
	}

	var err *model.AppError
	var user *model.User
	if user, err = app.GetUserByMobile(mobile); err != nil {
		c.Err = err
		return
	}

	if verificationCode != "" {
		if err = app.VerifiedCode(r, "", mobile, verificationCode); err != nil {
			c.Err = err
			return
		}
	} else {
		if err = app.CheckPasswordAndAllCriteria(user, password); err != nil {
			c.Err = err
			return
		}
	}

	err = app.DoLogin(w, r, user, deviceId)
	if err != nil {
		c.Err = err
		return
	}

	user.Sanitize(map[string]bool{})
	w.Write([]byte(user.ToJson()))
}

func isMobileExist(c *Context, w http.ResponseWriter, r *http.Request) {
	var mobileChar = regexp.MustCompile(`^[0-9]+$`)
    props := model.MapFromJson(r.Body)
    mobile  := props["mobile"]
    if mobile == "" || !mobileChar.MatchString(mobile) {
        c.SetInvalidParam("mobile")
        return
    }    	
	
	m := make(map[string]string)

	if exist := app.IsMobileExist(mobile); exist != true {
		m[model.STATUS] = "false"
	} else {
		m[model.STATUS] = "true"
	}

	w.Write([]byte(model.MapToJson(m)))
}

func resetPasswordByMobile(c *Context, w http.ResponseWriter, r *http.Request) {
    props := model.MapFromJson(r.Body)
    mobile  := props["mobile"]
    verificationCode := props["verification_code"]
    newPassword := props["new_password"]
    if mobile == "" {
            c.SetInvalidParam("mobile")
            return
    }
    if verificationCode == "" {
            c.SetInvalidParam("verification code")
            return
    }

    var err *model.AppError
    var user *model.User
    if user, err = app.GetUserByMobile(mobile); err != nil {
            c.Err = err
            return
    }
    if err = app.VerifiedCode(r, "", mobile, verificationCode); err != nil {
            c.Err = err
            return
    }
    if err := app.UpdatePassword(user, newPassword); err != nil {
            c.Err = err
            return
    }

	ReturnStatusOK(w)
}

func sendVerificationCodeEmail(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	email := props["email"]
	if len(email) == 0 {
		c.SetInvalidParam("email")
		return
	}

	if err := app.SendEmailVerificationCode(email); err != nil {
		// Don't want to leak whether the email is valid or not
		l4g.Error(err.Error())
		ReturnStatusOK(w)
		return
	}

	ReturnStatusOK(w)
}

func signupByEmail(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	email := props["email"]	
	if len(email) == 0 {
		c.SetInvalidParam("email")
		return
	}

	username := props["username"]	
	if username == "" {
		c.SetInvalidParam("username")
		return		
	}

	password := props["password"]
	if password == "" {
		c.SetInvalidParam("password")
		return	
	}

	nickname := props["nickname"]
	if nickname == "" {
		nickname = username	
	}

	if exist := app.IsEmailExist(email); exist == true {
		c.SetInvalidParam("email is already exist")
		return
	}

	verificationCode := props["verification_code"]
	if err := app.VerifiedCode(r, email, "", verificationCode); err != nil {
		c.Err = err
		return
	}

	user := &model.User{Username:username, Nickname:nickname, Email:email, Password:password, AllowMarketing: true}
	_, err := app.CreateUserFromEmailOrMobile(user)
	if err != nil {
		c.Err = err
		return
	}

	ruser, err := app.AuthenticateUserForLogin(user.Id, user.Username, password)
	if err != nil {
		c.Err = err
		return
	}

	err = app.DoLogin(w, r, ruser, "")
	if err != nil {
		c.Err = err
		return
	}

	ruser.Sanitize(map[string]bool{})
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ruser.ToJson()))
}

func isEmailExist(c *Context, w http.ResponseWriter, r *http.Request) {
	var emailChar = regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
    props := model.MapFromJson(r.Body)
    email  := props["email"]
    if email == "" || !emailChar.MatchString(email) {
            c.SetInvalidParam("email")
            return
    }    	

	m := make(map[string]string)

	if exist := app.IsEmailExist(email); exist != true {
		m[model.STATUS] = "false"
	} else {
		m[model.STATUS] = "true"
	}

	w.Write([]byte(model.MapToJson(m)))
}

func resetPasswordByEmail(c *Context, w http.ResponseWriter, r *http.Request) {
    props := model.MapFromJson(r.Body)
    email  := props["email"]
    verificationCode := props["verification_code"]
    newPassword := props["new_password"]
    if email == "" {
        c.SetInvalidParam("email")
        return
    }
    if verificationCode == "" {
        c.SetInvalidParam("verification code")
        return
    }

    var err *model.AppError
    var user *model.User
    if user, err = app.GetUserByEmail(email); err != nil {
        c.Err = err
        return
    }
    if err = app.VerifiedCode(r, email, "", verificationCode); err != nil {
        c.Err = err
        return
    }
    if err := app.UpdatePassword(user, newPassword); err != nil {
        c.Err = err
        return
    }

	ReturnStatusOK(w)
}

func searchUsers(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.UserSearchFromJson(r.Body)
	if props == nil {
		c.SetInvalidParam("")
		return
	}

	if len(props.Term) == 0 {
		c.SetInvalidParam("term")
		return
	}

	searchOptions := map[string]bool{}
	searchOptions[sqlStore.USER_SEARCH_OPTION_ALLOW_INACTIVE] = props.AllowInactive

	if profiles, err := app.SearchUsers(props, searchOptions, c.IsSystemAdmin()); err != nil {
		c.Err = err
		return
	} else {
		w.Write([]byte(model.UserListToJson(profiles)))
	}
}

