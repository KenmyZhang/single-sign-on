package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	ME = "me"

	DEFAULT_LOCALE = "zh-CN"

	USER_EMAIL_MAX_LENGTH     = 128
	USER_NICKNAME_MAX_RUNES   = 64
	USER_FIRST_NAME_MAX_RUNES = 64
	USER_LAST_NAME_MAX_RUNES  = 64
	USER_NAME_MAX_LENGTH      = 64
	USER_NAME_MIN_LENGTH      = 3
	USER_PASSWORD_MAX_LENGTH  = 72
)

type User struct {
	Id                 string    `bson:"_id" json:"id"`
	CreateAt           int64     `bson:"createAt" json:"create_at,omitempty"`
	UpdateAt           int64     `bson:"updateAt" json:"update_at,omitempty"`
	DeleteAt           int64     `bson:"deleteAt" json:"delete_at"`
	Username           string    `bson:"username" json:"username"`
	Gender             string    `bson:"gender"   json:"gender"`
	Password           string    `bson:"password" json:"password,omitempty"`
	AuthData           *string   `bson:"authData" json:"auth_data,omitempty"`
	AuthService        string    `bson:"authService" json:"auth_service"`
	Email              string    `bson:"email" json:"email"`
	EmailVerified      bool      `bson:"emailVerified" json:"email_verified,omitempty"`
	Nickname           string    `bson:"nickname" json:"nickname"`
	FirstName          string    `bson:"firstName" json:"first_name"`
	LastName           string    `bson:"lastName" json:"last_name"`
	Position           string    `bson:"position" json:"position"`
	Roles              string    `bson:"roles" json:"roles"`
	AllowMarketing     bool      `bson:"allowMarketing" json:"allow_marketing,omitempty"`
	Props              StringMap `bson:"props" json:"props,omitempty"`
	NotifyProps        StringMap `bson:"notifyProps" json:"notify_props,omitempty"`
	LastPasswordUpdate int64     `bson:"lastPasswordUpdate" json:"last_password_update,omitempty"`
	LastPictureUpdate  int64     `bson:"lastPictureUpdate" json:"last_picture_update,omitempty"`
	FailedAttempts     int       `bson:"failedAttempts" json:"failed_attempts,omitempty"`
	Locale             string    `bson:"locale" json:"locale"`
	MfaActive          bool      `bson:"mfaActive" json:"mfa_active,omitempty"`
	MfaSecret          string    `bson:"mfaSecret" json:"mfa_secret,omitempty"`
	LastActivityAt     int64     `bson:"-" db:"-" json:"last_activity_at,omitempty"`
	Names              []string  `bson:"names,omitempty" db:"-" json:"-"` 
	HasTeams           bool      `bson:"hasTeams" db:"-" json:"-"`        
	HeadImgUrl         string    `bson:"headImgUrl" json:"-"`            
	Mobile             string    `bson:"mobile" json:"mobile"`      
}

type LoginIdAndPassword struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

func (o *LoginIdAndPassword) ToJson() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (u *User) SetNames() {
	names := []string{u.Username, u.Email}

	if len(u.FirstName) > 0 {
		names = append(names, u.FirstName)
	}

	if len(u.LastName) > 0 {
		names = append(names, u.LastName)
	}

	if len(u.Nickname) > 0 {
		names = append(names, u.Nickname)
	}

	u.Names = names
}

func (u *User) IsValid() *AppError {
	if len(u.Id) != 26 {
		return InvalidUserError("id", "")
	}

	if u.CreateAt == 0 {
		return InvalidUserError("create_at", u.Id)
	}

	if u.UpdateAt == 0 {
		return InvalidUserError("update_at", u.Id)
	}

	if !IsValidUsername(u.Username) {
		return InvalidUserError("username", u.Id)
	}

	if (len(u.Email) > USER_EMAIL_MAX_LENGTH || len(u.Email) == 0) && u.AuthService == "" && u.Mobile == "" {
		return InvalidUserError("email", u.Id)
	}

	if utf8.RuneCountInString(u.Nickname) > USER_NICKNAME_MAX_RUNES {
		return InvalidUserError("nickname", u.Id)
	}

	if utf8.RuneCountInString(u.FirstName) > USER_FIRST_NAME_MAX_RUNES {
		return InvalidUserError("first_name", u.Id)
	}

	if utf8.RuneCountInString(u.LastName) > USER_LAST_NAME_MAX_RUNES {
		return InvalidUserError("last_name", u.Id)
	}

	if len(u.Password) > USER_PASSWORD_MAX_LENGTH {
		return InvalidUserError("password_limit", u.Id)
	}

	return nil
}

func InvalidUserError(fieldName string, userId string) *AppError {
	id := fmt.Sprintf("model.user.is_valid.%s.app_error", fieldName)
	details := ""
	if userId != "" {
		details = "user_id=" + userId
	}
	return NewAppError("User.IsValid", id, nil, details, http.StatusBadRequest)
}

func (u *User) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}

	if u.Username == "" {
		u.Username = NewId()
	}

	u.Username = strings.ToLower(u.Username)
	u.Email = strings.ToLower(u.Email)

	u.CreateAt = GetMillis()
	u.UpdateAt = u.CreateAt

	if u.Locale == "" {
		u.Locale = DEFAULT_LOCALE
	}

	if len(u.Password) > 0 {
		u.Password = HashPassword(u.Password)
	}
}

func (u *User) PreUpdate() {
	u.Username = strings.ToLower(u.Username)
	u.Email = strings.ToLower(u.Email)
	u.UpdateAt = GetMillis()
}

func (u *User) Etag(showFullName, showEmail bool) string {
	return Etag(u.Id, u.UpdateAt, showFullName, showEmail)
}

func (u *User) Sanitize(options map[string]bool) {
	u.Password = ""

	if len(options) != 0 && !options["email"] {
		u.Email = ""
	}
	if len(options) != 0 && !options["fullname"] {
		u.FirstName = ""
		u.LastName = ""
	}
}

func (u *User) IsOAuthUser() bool {
	return u.AuthService == USER_AUTH_SERVICE_WECHAT
}

func (u *User) MakeNonNil() {
	if u.Props == nil {
		u.Props = make(map[string]string)
	}

	if u.NotifyProps == nil {
		u.NotifyProps = make(map[string]string)
	}	
}

func (u *User) ClearNonProfileFields() {
	u.Password = ""
	u.FailedAttempts = 0
}

func (u *User) SanitizeProfile(options map[string]bool) {
	u.ClearNonProfileFields()

	u.Sanitize(options)
}

func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	} else if u.FirstName != "" {
		return u.FirstName
	} else if u.LastName != "" {
		return u.LastName
	} else {
		return ""
	}
}

func (u *User) GetDisplayName() string {
	if u.Nickname != "" {
		return u.Nickname
	} else if fullName := u.GetFullName(); fullName != "" {
		return fullName
	} else {
		return u.Username
	}
}

func IsValidUserRoles(userRoles string) bool {
	roles := strings.Fields(userRoles)

	for _, r := range roles {
		if !isValidRole(r) {
			return false
		}
	}

	// Exclude just the system_admin role explicitly to prevent mistakes
	if len(roles) == 1 && roles[0] == "system_admin" {
		return false
	}

	return true
}

func isValidRole(roleId string) bool {
	_, ok := BuiltInRoles[roleId]
	return ok
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(hash)
}

func ComparePassword(hash string, password string) bool {
	if len(password) == 0 || len(hash) == 0 {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var validUsernameChars = regexp.MustCompile(`^[a-z0-9\.\-_]+$`)
var firstChar = regexp.MustCompile(`^[a-z]`)

var restrictedUsernames = []string{
	"example",
}

func IsValidUsername(s string) bool {
	if len(s) < USER_NAME_MIN_LENGTH || len(s) > USER_NAME_MAX_LENGTH {
		return false
	}

	if !validUsernameChars.MatchString(s) {
		return false
	}

	if !firstChar.MatchString(s) {
		return false
	}

	for _, restrictedUsername := range restrictedUsernames {
		if s == restrictedUsername {
			return false
		}
	}

	return true
}

func (u *User) ToJson() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func UserListToJson(u []*User) string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func UserFromJson(data string) *User {
	user := &User{}
	err := json.Unmarshal([]byte(data), user)
	if err == nil {
		return user
	} else {
		return nil
	}
}

func UserListFromJson(data io.Reader) []*User {
	decoder := json.NewDecoder(data)
	var users []*User
	err := decoder.Decode(&users)
	if err == nil {
		return users
	} else {
		return nil
	}
}

func DecodeUserFromJson(data io.Reader) *User {
	decoder := json.NewDecoder(data)
	var user User
	err := decoder.Decode(&user)
	if err == nil {
		return &user
	} else {
		return nil
	}
}

func (user *User) UpdateMentionKeysFromUsername(oldUsername string) {
	nonUsernameKeys := []string{}
	splitKeys := strings.Split(user.NotifyProps["mention_keys"], ",")
	for _, key := range splitKeys {
		if key != oldUsername && key != "@"+oldUsername {
			nonUsernameKeys = append(nonUsernameKeys, key)
		}
	}

	user.NotifyProps["mention_keys"] = user.Username + ",@" + user.Username
	if len(nonUsernameKeys) > 0 {
		user.NotifyProps["mention_keys"] += "," + strings.Join(nonUsernameKeys, ",")
	}
}
