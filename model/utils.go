package model

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	goi18n "github.com/nicksnyder/go-i18n/i18n"
	"github.com/pborman/uuid"
)

const (
	LOWERCASE_LETTERS = "abcdefghijklmnopqrstuvwxyz"
	UPPERCASE_LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NUMBERS           = "0123456789"
	SYMBOLS           = " !\"\\#$%&'()*+,-./:;<=>?@[]^_`|~"

	USER_AUTH_SERVICE_LDAP = "ldap"
	LDAP_SYNC_TASK_NAME    = "LDAP Syncronization"
)

const (
	StatusBadRequest = http.StatusBadRequest
)

type StringInterface map[string]interface{}
type StringMap map[string]string
type StringArray []string
type EncryptStringMap map[string]string

type AppError struct {
	Id            string `json:"id"`
	Message       string `json:"message"`               // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"`        // Internal error string to help the developer
	RequestId     string `json:"request_id,omitempty"`  // The RequestId that's also set in the header
	StatusCode    int    `json:"status_code,omitempty"` // The http status code
	Where         string `json:"-"`                     // The function where it happened in the form of Struct.Func
	params        map[string]interface{}
}

func (er *AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

func (er *AppError) Translate(T goi18n.TranslateFunc) {
	if er.params == nil {
		er.Message = T(er.Id)
	} else {
		er.Message = T(er.Id, er.params)
	}
}

func (er *AppError) SystemMessage(T goi18n.TranslateFunc) string {
	if er.params == nil {
		return T(er.Id)
	} else {
		return T(er.Id, er.params)
	}
}

func NewAppError(where string, id string, params map[string]interface{}, details string, status int) *AppError {
	ap := &AppError{}
	ap.Id = id
	ap.params = params
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = status
	return ap
}

func (er *AppError) ToJson() string {
	b, err := json.Marshal(er)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func NewLocAppError(where string, id string, params map[string]interface{}, details string) *AppError {
	ap := &AppError{}
	ap.Id = id
	ap.params = params
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = 500
	return ap
}

// AppErrorFromJson will decode the input and return an AppError
func AppErrorFromJson(data io.Reader) *AppError {
	str := ""
	bytes, rerr := ioutil.ReadAll(data)
	if rerr != nil {
		str = rerr.Error()
	} else {
		str = string(bytes)
	}

	decoder := json.NewDecoder(strings.NewReader(str))
	var er AppError
	err := decoder.Decode(&er)
	if err == nil {
		return &er
	} else {
		return NewLocAppError("AppErrorFromJson", "model.utils.decode_json.app_error", nil, "body: "+str)
	}
}

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

func NewId() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26) // removes the '==' padding
	return b.String()
}

func NewRandomString(length int) string {
	var b bytes.Buffer
	str := make([]byte, length+8)
	rand.Read(str)
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(str)
	encoder.Close()
	b.Truncate(length) // removes the '==' padding
	return b.String()
}

func MapFromJson(data io.Reader) map[string]string {
	decoder := json.NewDecoder(data)

	var objmap map[string]string
	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]string)
	} else {
		return objmap
	}
}

// MapToJson converts a map to a json string
func MapToJson(objmap map[string]string) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func ArrayFromJson(data io.Reader) []string {
	decoder := json.NewDecoder(data)

	var objmap []string
	if err := decoder.Decode(&objmap); err != nil {
		return make([]string, 0)
	} else {
		return objmap
	}
}

func ArrayToJson(objmap []string) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func IsValidHttpUrl(rawUrl string) bool {
	if strings.Index(rawUrl, "http://") != 0 && strings.Index(rawUrl, "https://") != 0 {
		return false
	}

	if _, err := url.ParseRequestURI(rawUrl); err != nil {
		return false
	}

	return true
}

func StringInterfaceToJson(objmap map[string]interface{}) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Etag(parts ...interface{}) string {
	etag := CurrentVersion

	for _, part := range parts {
		etag += fmt.Sprintf(".%v", part)
	}

	return etag
}

var reservedName = []string{
	"signup",
	"login",
	"admin",
	"api",
	"oauth",
}

type Status struct {
	Status  string `json:"status"`
}


func StatusFromJson(data io.Reader) *Status {
	decoder := json.NewDecoder(data)

	var status Status
	if err := decoder.Decode(&status); err != nil {
		return nil
	} else {
		return &status
	}
}