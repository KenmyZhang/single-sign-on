package oauthwechat

import (
	"encoding/json"
	"github.com/KenmyZhang/single-sign-on/einterfaces"
	"github.com/KenmyZhang/single-sign-on/model"
	"io"
)

type WeChatProvider struct {
}

type WeChatUser struct {
	OpenId 		string  `json:"openid"`
	NickName    string  `json:"nickname"`
	Sex			int     `json:"sex"`
	Language    string  `json:"Language"`
	City		string  `json:"city"`
	Province    string  `json:"province"`
	Country     string  `json:"country"`
	HeadImgUrl  string  `json:"headimgurl"`
	Privilege []string  `json:"privilege"`
	Unionid		string  `json:"unionid"`
}

func init() {
	provider := &WeChatProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_WECHAT, provider)
}

func userFromWeChatUser(wcu *WeChatUser) *model.User{
	user := &model.User{}
	username := "weixin" + wcu.OpenId
	user.Username = username
	user.AuthData = &wcu.OpenId
	user.AuthService = model.USER_AUTH_SERVICE_WECHAT 
	user.Nickname = wcu.NickName
	user.HeadImgUrl = wcu.HeadImgUrl
	return user
}

func weChatUserFromJson(data io.Reader) *WeChatUser {
	decoder := json.NewDecoder(data)
	var wcu WeChatUser
	err :=decoder.Decode(&wcu)
	if err == nil {
		return &wcu	
	} else {
		return nil
	}
}

func (wcu *WeChatUser) ToJson() string {
	b, err := json.Marshal(wcu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (wcu *WeChatUser) IsValid() bool {
	if len(wcu.OpenId) == 0 {
		return false
	}

	return true
}

func (wcu *WeChatUser) getAuthData() string {
	return wcu.OpenId
}

func (m *WeChatProvider) GetIdentifier() string {
	return model.USER_AUTH_SERVICE_WECHAT
}

func (m *WeChatProvider) GetUserFromJson(data io.Reader) *model.User {
	wcu := weChatUserFromJson(data)
	if wcu.IsValid() {
		return userFromWeChatUser(wcu)
	}

	return &model.User{}
}

func (m *WeChatProvider) GetAuthDataFromJson(data io.Reader) string {
	wcu := weChatUserFromJson(data)

	if wcu.IsValid() {
		return wcu.getAuthData()
	}

	return ""
}
