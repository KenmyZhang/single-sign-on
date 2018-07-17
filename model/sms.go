package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ALiYunCommunicationRequest struct {
	//system parameters
	AccessKeyId      string
	Timestamp        string
	Format           string
	SignatureMethod  string
	SignatureVersion string
	SignatureNonce   string
	Signature        string

	//business parameters
	Action          string
	Version         string
	RegionId        string
	PhoneNumbers    string
	SignName        string
	TemplateCode    string
	TemplateParam   string
	SmsUpExtendCode string
	OutId           string
}

func (req *ALiYunCommunicationRequest) SetParamsValue(accessKeyId, phoneNumbers, signName, templateCode, templateParam string) error {
	req.AccessKeyId = accessKeyId
	now := time.Now()
	local, err := time.LoadLocation("GMT")
	if err != nil {
		return err
	}
	req.Timestamp = now.In(local).Format("2006-01-02T15:04:05Z")
	fmt.Println("time:", req.Timestamp)
	req.Format = "json"
	req.SignatureMethod = "HMAC-SHA1"
	req.SignatureVersion = "1.0"
	req.SignatureNonce = NewId()
	fmt.Println("req.SignatureNonce:", req.SignatureNonce)

	req.Action = "SendSms"
	req.Version = "2017-05-25"
	req.RegionId = "cn-hangzhou"
	req.PhoneNumbers = phoneNumbers
	req.SignName = signName
	req.TemplateCode = templateCode
	req.TemplateParam = templateParam
	req.SmsUpExtendCode = "90999"
	req.OutId = "abcdefg"
	return nil
}

func (req *ALiYunCommunicationRequest) SmsParamsIsValid() error {
	if len(req.AccessKeyId) == 0 {
		return errors.New("AccessKeyId required")
	}

	if len(req.PhoneNumbers) == 0 {
		return errors.New("PhoneNumbers required")
	}

	if len(req.SignName) == 0 {
		return errors.New("SignName required")
	}

	if len(req.TemplateCode) == 0 {
		return errors.New("TemplateCode required")
	}

	if len(req.TemplateParam) == 0 {
		return errors.New("TemplateParam required")
	}

	return nil
}

func (req *ALiYunCommunicationRequest) BuildSmsRequestEndpoint(accessKeySecret, gatewayUrl string) (string, error) {
	var err error
	if err = req.SmsParamsIsValid(); err != nil {
		return "", err
	}
	// common params
	systemParams := make(map[string]string)
	systemParams["SignatureMethod"] = req.SignatureMethod
	systemParams["SignatureNonce"] = req.SignatureNonce
	systemParams["AccessKeyId"] = req.AccessKeyId
	systemParams["SignatureVersion"] = req.SignatureVersion
	systemParams["Timestamp"] = req.Timestamp
	systemParams["Format"] = req.Format

	// business params
	businessParams := make(map[string]string)
	businessParams["Action"] = req.Action
	businessParams["Version"] = req.Version
	businessParams["RegionId"] = req.RegionId
	businessParams["PhoneNumbers"] = req.PhoneNumbers
	businessParams["SignName"] = req.SignName
	businessParams["TemplateParam"] = req.TemplateParam
	businessParams["TemplateCode"] = req.TemplateCode
	businessParams["SmsUpExtendCode"] = req.SmsUpExtendCode
	businessParams["OutId"] = req.OutId
	// generate signature and sorted  query
	sortQueryString, signature := generateQueryStringAndSignature(businessParams, systemParams, accessKeySecret)
	fmt.Println("Signature:", signature)
	fmt.Println("sortQueryString:", sortQueryString)
	return gatewayUrl + "?Signature=" + signature + sortQueryString, nil
}

func generateQueryStringAndSignature(businessParams map[string]string, systemParams map[string]string, accessKeySecret string) (string, string) {
	keys := make([]string, 0)
	allParams := make(map[string]string)
	for key, value := range businessParams {
		keys = append(keys, key)
		allParams[key] = value
	}

	for key, value := range systemParams {
		keys = append(keys, key)
		allParams[key] = value
	}

	sort.Strings(keys)

	sortQueryStringTmp := ""
	for _, key := range keys {
		rstkey := specialUrlEncode(key)
		rstval := specialUrlEncode(allParams[key])
		sortQueryStringTmp = sortQueryStringTmp + "&" + rstkey + "=" + rstval
	}

	sortQueryString := strings.Replace(sortQueryStringTmp, "&", "", 1)
	stringToSign := "GET" + "&" + specialUrlEncode("/") + "&" + specialUrlEncode(sortQueryString)

	sign := sign(accessKeySecret+"&", stringToSign)
	signature := specialUrlEncode(sign)
	return sortQueryStringTmp, signature
}

func specialUrlEncode(value string) string {
	rstValue := url.QueryEscape(value)
	rstValue = strings.Replace(rstValue, "+", "%20", -1)
	rstValue = strings.Replace(rstValue, "*", "%2A", -1)
	rstValue = strings.Replace(rstValue, "%7E", "~", -1)
	return rstValue
}

func sign(accessKeySecret, sortquerystring string) string {
	h := hmac.New(sha1.New, []byte(accessKeySecret))
	h.Write([]byte(sortquerystring))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}