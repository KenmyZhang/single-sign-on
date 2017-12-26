package model

import (
	"encoding/json"
	"io"
)

const (
	ACCESS_TOKEN_GRANT_TYPE  = "authorization_code"
	ACCESS_TOKEN_TYPE        = "bearer"
	REFRESH_TOKEN_GRANT_TYPE = "refresh_token"
)

type AccessData struct {
	ClientId     string `bson:"clientId" json:"client_id"`
	UserId       string `bson:"userId" json:"user_id"`
	Token        string `bson:"token" json:"token"`
	RefreshToken string `bson:"refreshToken" json:"refresh_token"`
	RedirectUri  string `bson:"redirectUri" json:"redirect_uri"`
	ExpiresAt    int64  `bson:"expiresAt" json:"expires_at"`
	Scope        string `bson:"scope" json:"scope"`
}

type AccessResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
}

func (ad *AccessData) IsValid() *AppError {

	if len(ad.ClientId) == 0 || len(ad.ClientId) > 26 {
		return NewLocAppError("AccessData.IsValid", "model.access.is_valid.client_id.app_error", nil, "")
	}

	if len(ad.UserId) == 0 || len(ad.UserId) > 26 {
		return NewLocAppError("AccessData.IsValid", "model.access.is_valid.user_id.app_error", nil, "")
	}

	if len(ad.Token) != 26 {
		return NewLocAppError("AccessData.IsValid", "model.access.is_valid.access_token.app_error", nil, "")
	}

	if len(ad.RefreshToken) > 26 {
		return NewLocAppError("AccessData.IsValid", "model.access.is_valid.refresh_token.app_error", nil, "")
	}

	if len(ad.RedirectUri) == 0 || len(ad.RedirectUri) > 256 || !IsValidHttpUrl(ad.RedirectUri) {
		return NewLocAppError("AccessData.IsValid", "model.access.is_valid.redirect_uri.app_error", nil, "")
	}

	return nil
}

func (me *AccessData) IsExpired() bool {

	if me.ExpiresAt <= 0 {
		return false
	}

	if GetMillis() > me.ExpiresAt {
		return true
	}

	return false
}

func (ad *AccessData) ToJson() string {
	b, err := json.Marshal(ad)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func AccessDataFromJson(data io.Reader) *AccessData {
	decoder := json.NewDecoder(data)
	var ad AccessData
	err := decoder.Decode(&ad)
	if err == nil {
		return &ad
	} else {
		return nil
	}
}

func (ar *AccessResponse) ToJson() string {
	b, err := json.Marshal(ar)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func AccessResponseFromJson(data io.Reader) *AccessResponse {
	decoder := json.NewDecoder(data)
	var ar AccessResponse
	err := decoder.Decode(&ar)
	if err == nil {
		return &ar
	} else {
		return nil
	}
}
