package app

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"	
	"net/http"
	"net/url"
	"time"
	"errors"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

type ALiYunSmsClient struct {
	Request   *model.ALiYunCommunicationRequest
	GatewayUrl string
	Client    *http.Client
}

type IhuyiSmsClient struct {
	GatewayUrl string
	Client    *http.Client
}

func NewALiYunSmsClient(gatewayUrl string) *ALiYunSmsClient {
	smsClient := new(ALiYunSmsClient)
	smsClient.Request = &model.ALiYunCommunicationRequest{}
	smsClient.GatewayUrl = gatewayUrl
	smsClient.Client = &http.Client{}
	return smsClient
}

func (smsClient *ALiYunSmsClient) Execute(accessKeyId, accessKeySecret, mobile, signName, templateCode, templateParam string) (err error){
	var endpoint string
	if err = smsClient.Request.SetParamsValue(accessKeyId, mobile, signName, templateCode, templateParam); err != nil {
		return 
	}
	if endpoint, err = smsClient.Request.BuildSmsRequestEndpoint(accessKeySecret, smsClient.GatewayUrl); err != nil {
		return 
	}

	request, _ := http.NewRequest("GET",endpoint, nil)
	if err != nil {
		return 
	}		
	response, _ := smsClient.Client.Do(request)
	if err != nil {
		return 
	}		
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 
	}	
	defer response.Body.Close()	

	var result map[string] interface{}
	err = json.Unmarshal(body, &result)
	for key, value := range result {
		 l4g.Debug("key:", key, " value:",value)
	}

	return 
}


func NewIhuyiSmsClient(gatewayUrl string) *IhuyiSmsClient {
	smsClient := new(IhuyiSmsClient)
	smsClient.GatewayUrl = gatewayUrl
	smsClient.Client = &http.Client{}
	return smsClient
}

func (smsClient *IhuyiSmsClient) Execute(accessKeyId, accessKeySecret, mobile, signName, templateCode, templateParam string) (err error) {
	v := url.Values{}
	nowTime := strconv.FormatInt(time.Now().Unix(), 10)
	content := "您的验证码是：" + templateParam + "。请不要把验证码泄露给其他人。"  
    v.Set("account", accessKeyId)
    v.Set("password", GetMd5String(accessKeyId+accessKeySecret+mobile+content+nowTime))
    v.Set("mobile", mobile)
    v.Set("content", content)
    v.Set("time", nowTime)
    v.Set("format", "json")
    body := ioutil.NopCloser(strings.NewReader(v.Encode())) 
	req, err := http.NewRequest("POST", smsClient.GatewayUrl, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
 	smsResult := &model.IhuyiSmsResult{}
 	var response *http.Response
 	if response, err = utils.HttpClient().Do(req); err != nil {
		return
	} else {
		smsResult = model.IhuyiSmsResultFromJson(response.Body)
		defer CloseBody(response)
	}

	if smsResult.Code != 2 {
		err = errors.New("submit verification code error")
		return 
	}
	return nil
}


