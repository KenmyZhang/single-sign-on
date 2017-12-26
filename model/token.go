package model

import "net/http"

const (
	TOKEN_SIZE            = 64
	MAX_TOKEN_EXIPRY_TIME = 1000 * 60 * 60 * 24 // 24 hour
	MAX_SMS_TOKEN_EXIPRY_TIME = 1000 * 60 * 1 // 1 min
	MAX_EMAIL_TOKEN_EXIPRY_TIME = 1000 * 60 * 1 // 1 min
	SEND_CODE_MAX = 60
	TOKEN_TYPE_OAUTH      = "oauth"
)

type Token struct {
	Token    string `bson:"_id"`
	CreateAt int64  `bson:"createAt"`
	Type     string
	Props    StringInterface `bson:"props" json:"props"`
	Extra    string
}

func NewToken(tokentype, extra string, props StringInterface) *Token {
	return &Token{
		Token:    NewRandomString(TOKEN_SIZE),
		CreateAt: GetMillis(),
		Type:     tokentype,
		Props: props,
		Extra:    extra,
	}
}

func (t *Token) IsValid() *AppError {
	if len(t.Token) != TOKEN_SIZE {
		return NewAppError("Token.IsValid", "model.token.is_valid.size", nil, "", http.StatusInternalServerError)
	}

	if t.CreateAt == 0 {
		return NewAppError("Token.IsValid", "model.token.is_valid.expiry", nil, "", http.StatusInternalServerError)
	}

	return nil
}
