package app

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/KenmyZhang/single-sign-on/model"	
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type ALiYunSmsClient struct {
	Request   *model.ALiYunCommunicationRequest
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



