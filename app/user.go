package app

import (
	"net/http"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"crypto/md5"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"mime/multipart"
	"github.com/golang/freetype"
	"image/png"
	"bytes"
	//b64 "encoding/base64"
	"encoding/hex"	
	"github.com/disintegration/imaging"
	l4g "github.com/alecthomas/log4go"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

const (
	TOKEN_TYPE_PASSWORD_RECOVERY  = "password_recovery"
	TOKEN_TYPE_VERIFY_EMAIL       = "verify_email"
	TOKEN_TYPE_VERIFY_MOBILE      = "verify_mobile"
	PASSWORD_RECOVER_EXPIRY_TIME  = 1000 * 60 * 60 // 1 hour
	VERIFY_EMAIL_EXPIRY_TIME      = 1000 * 60 * 60 // 1 hour
	IMAGE_PROFILE_PIXEL_DIMENSION = 128
	SMS_COOKIE_MAX_AGE_SECONDS    = 30 * 60 // 30 minutes
	VERIFICATION_CODE                      = "verified_code"	
	COOKIE_SMS_TOKEN              = "verified_token"
)

func IsFirstUserAccount() bool {
	if cr := <-Srv.SqlStore.User().GetTotalUsersCount(); cr.Err != nil {
		l4g.Error(cr.Err)
		return false
	} else {
		count := cr.Data.(int64)
		if count <= 0 {
			return true
		}
	}

	return false
}

func CreateUser(user *model.User) (*model.User, *model.AppError) {
	user.Roles = model.ROLE_USER.Id

	if result := <-Srv.SqlStore.User().GetTotalUsersCount(); result.Err != nil {
		return nil, result.Err
	} else {
		count := result.Data.(int64)
		if count <= 0 {
			user.Roles = model.ROLE_SYSTEM_ADMIN.Id
		}
	}

	if _, ok := utils.GetSupportedLocales()[user.Locale]; !ok {
		user.Locale = *utils.Cfg.LocalizationSettings.DefaultClientLocale
	}

	if ruser, err := createUser(user); err != nil {
		return nil, err
	} else {
		return ruser, nil
	}
}

func createUser(user *model.User) (*model.User, *model.AppError) {
	user.MakeNonNil()

	if err := utils.IsPasswordValid(user.Password); user.AuthService == "" && err != nil {
		return nil, err
	}

	if result := <-Srv.SqlStore.User().Save(user); result.Err != nil {
		l4g.Error(utils.T("api.user.create_user.save.error"), result.Err)
		return nil, result.Err
	} else {
		ruser := result.Data.(*model.User)
		ruser.Sanitize(map[string]bool{})
		return ruser, nil
	}
}

func IsUserSignUpAllowed() *model.AppError {
	if !utils.Cfg.EmailSettings.EnableSignUpWithEmail {
		err := model.NewLocAppError("IsUserSignUpAllowed",
			"api.user.create_user.signup_email_disabled.app_error", nil, "",
		)
		err.StatusCode = http.StatusNotImplemented
		return err
	}
	return nil
}

func makeVerificationCode() (code string) {
	code = strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(899999) + 100000)
	return
}

func GetMd5String(s string) string {
    md := md5.New()
    md.Write([]byte(s))
    return hex.EncodeToString(md.Sum(nil))
}

func GetUser(userId string) (*model.User, *model.AppError) {
	if result := <-Srv.SqlStore.User().Get(userId); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func SearchUsers(props *model.UserSearch, searchOptions map[string]bool, asAdmin bool) ([]*model.User, *model.AppError) {
	if result := <-Srv.SqlStore.User().SearchUsers(props.Term, searchOptions); result.Err != nil {
		return nil, result.Err
	} else {
		users := result.Data.([]*model.User)

		for _, user := range users {
			SanitizeProfile(user, asAdmin)
		}

		return users, nil
	}
}

func GetUserForLogin(loginId string) (*model.User, *model.AppError) {
	if result := <-Srv.SqlStore.User().GetForLogin(
		loginId,
		*utils.Cfg.EmailSettings.EnableSignInWithUsername,
		*utils.Cfg.EmailSettings.EnableSignInWithEmail,
		*utils.Cfg.EnableSignInWithMobile,
		utils.Cfg.WeixinSettings.Enable,
	); result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func SanitizeProfile(user *model.User, asAdmin bool) {
	options := utils.Cfg.GetSanitizeOptions()
	if asAdmin {
		options["email"] = true
		options["fullname"] = true
		options["authservice"] = true
	}
	user.SanitizeProfile(options)
}

func GetUserByAuth(authData *string, authService string) (*model.User, *model.AppError) {
	if result := <-Srv.SqlStore.User().GetByAuth(authData, authService); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func GetUserByUsername(username string) (*model.User, *model.AppError) {
	if result := <-Srv.SqlStore.User().GetByUsername(username); result.Err != nil && result.Err.Id == "store.mgo_user.get_by_username.app_error" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func SendSmsCode(mobile string, w http.ResponseWriter, r *http.Request) *model.AppError {
	sso := utils.Cfg.GetSmsService()
	if sso != nil && !sso.Enable {
		return model.NewAppError("SendSmsCode", "api.user.send_sms_code.unsupported.app_error", nil, "Sms disabled  or not support", http.StatusNotImplemented)
	}  

	var tokenStr string
	if token, err := GetTokenByExtra(mobile); err == nil {
		if token.CreateAt > (model.GetMillis() - model.MAX_SMS_TOKEN_EXIPRY_TIME) {
			return model.NewAppError("SendSmsCode", "api.user.send_sms_code.less_than_60.app_error", nil, 
				"you should wait more than 60s,now time" + strconv.FormatInt(model.GetMillis(), 10) + "; create time:" + strconv.FormatInt(token.CreateAt, 10), http.StatusBadRequest)
		}
	}

	if count, err := GetTokenCountByExtra(mobile); err == nil {
		if count > model.SEND_CODE_MAX {
			return model.NewAppError("SendSmsCode", "api.user.send_sms_code.get_token_count_by_count.app_error", nil, "you send more than " + strconv.Itoa(model.SEND_CODE_MAX), http.StatusBadRequest)
		}
	}

    verificationCode := makeVerificationCode()
	if sso.Provider == "aliyun" {
    	//阿里云-云通信		
		templateParam := `{"code":"` + verificationCode + `"}`
 		smsClient := NewALiYunSmsClient(sso.GatewayUrl)
 		if err := smsClient.Execute(sso.AccessKeyId, sso.AccessKeySecret, mobile, sso.SignName, 
 			sso.TemplateCode, templateParam); err != nil {
 			return model.NewLocAppError("SendSmsCode", "api.user.send_sms_code.aliyun_execute.app_error", nil, "smsClient.Execute() err:" + err.Error())
 		}
	} else { 
    	verificationCode = "666666"
    }
	secure := false
	if GetProtocol(r) == "https" {
		secure = true
	}

	expiresAt := time.Unix(model.GetMillis()/1000+int64(60), 0)

	if token, err := CreateVerifyCodeToken(mobile, verificationCode); err != nil {
		return model.NewAppError("SendSmsCode", "api.user.send_sms_code.create_token.app_error", nil, "create token error" + err.Error(), http.StatusBadRequest)
	} else {
		tokenStr = token.Token
	}	

	tokenCookie := &http.Cookie{
		Name:     COOKIE_SMS_TOKEN,
		Value:    tokenStr,
		Path:     "/",
		MaxAge:   60,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
	}
	http.SetCookie(w, tokenCookie)

	return nil
}

func CreateVerifyCodeToken(mobile, verifiedCode string) (*model.Token, *model.AppError) {
	props := model.StringInterface{}
	props["verified_code"] = verifiedCode
	token := model.NewToken(TOKEN_TYPE_VERIFY_MOBILE, mobile, props)

	if result := <-Srv.SqlStore.Token().Save(token); result.Err != nil {
		return nil, result.Err
	}

	return token, nil
}

func CreateUserFromEmailOrMobile(user *model.User) (*model.User, *model.AppError) {	
	user.EmailVerified = false

	ruser, err := CreateUser(user)
	if err != nil {
		return nil, err
	}
	return ruser, nil
}

func IsMobileExist(mobile string) bool {
	if result := <-Srv.SqlStore.User().GetProfileByMobile(mobile); result.Err != nil {
		return false
	}
	return true
}

func SendEmailVerificationCode(email string) *model.AppError {
	verificationCode := makeVerificationCode()
	if token, err := GetTokenByExtra(email); err == nil {
		if token.CreateAt > (model.GetMillis() - model.MAX_EMAIL_TOKEN_EXIPRY_TIME) {
			return model.NewAppError("SendEmailVerificationCode", "api.user.send_email_verification_code.get_token.app_error", nil, "you should wait more than 60s", http.StatusBadRequest)
		}
	}

	if _, err := CreateVerifyCodeToken(email, verificationCode); err != nil {
		return model.NewAppError("SendEmailVerificationCode", "api.user.send_email_verification_code.create_token.app_error", nil, "create token error" + err.Error(), http.StatusBadRequest)
	}

	if utils.Cfg.LocalizationSettings.AvailableLocales != nil && *utils.Cfg.LocalizationSettings.AvailableLocales != "" {
		return SendVerificationCodeToEmail(email, *utils.Cfg.LocalizationSettings.AvailableLocales, utils.GetSiteURL(), verificationCode)
	}
	return model.NewAppError("SendEmailVerificationCode", "api.user.send_email_verification_code.none_locales_available", nil, "none locales available", http.StatusNotImplemented)
}

func IsEmailExist(email string) bool {
	if result := <-Srv.SqlStore.User().GetByEmail(email); result.Err != nil {
		return false
	}
	return true
}

func UpdatePassword(user *model.User, newPassword string) *model.AppError {
	if err := utils.IsPasswordValid(newPassword); err != nil {
		return err
	}

	hashedPassword := model.HashPassword(newPassword)

	if result := <-Srv.SqlStore.User().UpdatePassword(user.Id, hashedPassword); result.Err != nil {
		return model.NewLocAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, result.Err.Error())
	}

	return nil
}

func GetUserByMobile(mobile string) (*model.User, *model.AppError) {
	if result := <-Srv.SqlStore.User().GetProfileByMobile(mobile); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func GetUserByEmail(email string) (*model.User, *model.AppError) {

	if result := <-Srv.SqlStore.User().GetByEmail(email); result.Err != nil && result.Err.Id == "store.sql_user.missing_account.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func IsUsernameTaken(name string) bool {

	if !model.IsValidUsername(name) {
		return false
	}

	if result := <-Srv.SqlStore.User().GetByUsername(name); result.Err != nil {
		return false
	}

	return true
}

func IsDoctor(userId string) bool {
	if result := <-Srv.SqlStore.User().Get(userId); result.Err != nil {
		return false
	} else {
		user := result.Data.(*model.User)
		if true == strings.Contains(user.Roles, "doctor") {
			return true
		}
	}
	return false
}

func VerifiedCode(r *http.Request, email, mobile, verificationCode string) *model.AppError{
    var err *model.AppError
    var token *model.Token
    var extra, extraInToken, verificationCodeInToken string

	if email == "" {
		extra = mobile
    } else {
    	extra = email
    }

	if token, err = GetTokenByExtra(extra); err != nil {
		return model.NewLocAppError("VerifiedCode", "api.user.verified_code.get_token.app_error", nil, err.Error())
    }    

    expireTime := model.GetMillis() - model.MAX_SMS_TOKEN_EXIPRY_TIME
    if token.CreateAt < expireTime {
    	return model.NewLocAppError("VerifiedCode", "api.user.verified_code.expire.app_error", nil, "verified code expire")
    }
    
    extraInToken = token.Extra
    if value, ok := token.Props[VERIFICATION_CODE]; ok {
    	verificationCodeInToken = value.(string)
    }     	

    if extraInToken != extra {
    	return model.NewLocAppError("VerifiedCode", "api.user.verified_code.extra.app_error", nil, 
    		"extra in token is " + extraInToken + ", while input extra is " + extra)
    }

    if verificationCodeInToken != verificationCode {
    	return model.NewLocAppError("VerifiedCode", "api.user.verified_code.verification_code.app_error", nil, 
    		"verification_code error:" + verificationCode)
    }

    return nil
}

func GetProfileImage(user *model.User) ([]byte, bool, *model.AppError) {
	var img []byte
	readFailed := false

	path := "users/" + user.Id + "/profile.png"

	if data, err := ReadFile(path); err != nil {
		readFailed = true

		if img, err = CreateProfileImage(user.Username, user.Id); err != nil {
			return nil, false, err
		}

		if user.LastPictureUpdate == 0 {
			if err := writeFileLocally(img, path); err != nil {
				return nil, false, err
			}
		}

	} else {
			img = data
	}

	return img, readFailed, nil
}

func CreateProfileImage(username string, userId string) ([]byte, *model.AppError) {
	colors := []color.NRGBA{
		{197, 8, 126, 255},
		{227, 207, 18, 255},
		{28, 181, 105, 255},
		{35, 188, 224, 255},
		{116, 49, 196, 255},
		{197, 8, 126, 255},
		{197, 19, 19, 255},
		{250, 134, 6, 255},
		{227, 207, 18, 255},
		{123, 201, 71, 255},
		{28, 181, 105, 255},
		{35, 188, 224, 255},
		{116, 49, 196, 255},
		{197, 8, 126, 255},
		{197, 19, 19, 255},
		{250, 134, 6, 255},
		{227, 207, 18, 255},
		{123, 201, 71, 255},
		{28, 181, 105, 255},
		{35, 188, 224, 255},
		{116, 49, 196, 255},
		{197, 8, 126, 255},
		{197, 19, 19, 255},
		{250, 134, 6, 255},
		{227, 207, 18, 255},
		{123, 201, 71, 255},
	}

	h := fnv.New32a()
	h.Write([]byte(userId))
	seed := h.Sum32()

	initial := string(strings.ToUpper(username)[0])

	fontDir, _ := utils.FindDir("fonts")
	fontBytes, err := ioutil.ReadFile(fontDir + *utils.Cfg.FileSettings.InitialFont)
	if err != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.default_font.app_error", nil, err.Error())
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.default_font.app_error", nil, err.Error())
	}

	color := colors[int64(seed)%int64(len(colors))]
	dstImg := image.NewRGBA(image.Rect(0, 0, IMAGE_PROFILE_PIXEL_DIMENSION, IMAGE_PROFILE_PIXEL_DIMENSION))
	srcImg := image.White
	draw.Draw(dstImg, dstImg.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	size := float64(IMAGE_PROFILE_PIXEL_DIMENSION / 2)

	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(dstImg.Bounds())
	c.SetDst(dstImg)
	c.SetSrc(srcImg)

	pt := freetype.Pt(IMAGE_PROFILE_PIXEL_DIMENSION/6, IMAGE_PROFILE_PIXEL_DIMENSION*2/3)
	_, err = c.DrawString(initial, pt)
	if err != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.initial.app_error", nil, err.Error())
	}

	buf := new(bytes.Buffer)

	if imgErr := png.Encode(buf, dstImg); imgErr != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.encode.app_error", nil, imgErr.Error())
	} else {
		return buf.Bytes(), nil
	}
}

func SetProfileImage(w http.ResponseWriter, userId string, imageData *multipart.FileHeader) *model.AppError {
	file, err := imageData.Open()
	defer file.Close()
	if err != nil {
		return model.NewLocAppError("SetProfileImage", "api.user.upload_profile_user.open.app_error", nil, err.Error())
	}

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return model.NewLocAppError("SetProfileImage", "api.user.upload_profile_user.decode_config.app_error", nil, err.Error())
	} else if config.Width*config.Height > MaxImageSize {
		return model.NewLocAppError("SetProfileImage", "api.user.upload_profile_user.too_large.app_error", nil, err.Error())
	}

	file.Seek(0, 0)

	img, _, err := image.Decode(file)
	if err != nil {
		return model.NewLocAppError("SetProfileImage", "api.user.upload_profile_user.decode.app_error", nil, err.Error())
	}

	file.Seek(0, 0)

	orientation, _ := getImageOrientation(file)
	img = makeImageUpright(img, orientation)

	profileWidthAndHeight := 128
	img = imaging.Fill(img, profileWidthAndHeight, profileWidthAndHeight, imaging.Center, imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return model.NewLocAppError("SetProfileImage", "api.user.upload_profile_user.encode.app_error", nil, err.Error())
	}

	path := "users/" + userId + "/profile.png"

	if err := writeFileLocally(buf.Bytes(), *utils.Cfg.FileSettings.Directory + path); err != nil {
		return model.NewLocAppError("SetProfileImage", "api.user.upload_profile_user.upload_profile.app_error", nil, "")
	}

	result := <-Srv.SqlStore.User().UpdateLastPictureUpdate(userId)
	timestamp := strconv.FormatInt(result.Data.(int64), 10)
	w.Header().Set(model.LAST_PICTURE_UPDATE, timestamp)

	return nil
}