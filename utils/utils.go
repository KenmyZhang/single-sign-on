package utils

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/KenmyZhang/single-sign-on/model"
)

func StringArrayIntersection(arr1, arr2 []string) []string {
	arrMap := map[string]bool{}
	result := []string{}

	for _, value := range arr1 {
		arrMap[value] = true
	}

	for _, value := range arr2 {
		if arrMap[value] {
			result = append(result, value)
		}
	}

	return result
}

func StringArrayContains(arr1, arr2 []string) bool {
	arrMap := map[string]bool{}

	for _, value := range arr1 {
		arrMap[value] = true
	}

	for _, value := range arr2 {
		if !arrMap[value] {
			return false
		}
	}

	return true
}

func FileExistsInConfigFolder(filename string) bool {
	if len(filename) == 0 {
		return false
	}

	if _, err := os.Stat(FindConfigFile(filename)); err == nil {
		return true
	}
	return false
}

func RemoveDuplicatesFromStringArray(arr []string) []string {
	result := make([]string, 0, len(arr))
	seen := make(map[string]bool)

	for _, item := range arr {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

func GetIpAddress(r *http.Request) string {
	address := r.Header.Get(model.HEADER_FORWARDED)

	if len(address) == 0 {
		address = r.Header.Get(model.HEADER_REAL_IP)
	}

	if len(address) == 0 {
		address, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return address
}

func GetHostnameFromSiteURL(siteURL string) string {
	u, err := url.Parse(siteURL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

func MapToJson(objmap map[string]string) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func UrlEncode(str string) string {
	strs := strings.Split(str, " ")

	for i, s := range strs {
		strs[i] = url.QueryEscape(s)
	}

	return strings.Join(strs, "%20")
}

func GetDisplayName(user *model.User) string {
	if user == nil {
		return ""
	}
	if len(strings.TrimSpace(user.Nickname)) != 0 {
		return user.Nickname
	}

	fullName := getFullName(user)
	if len(fullName) != 0 {
		return fullName
	}

	return user.Username
}

func getFullName(user *model.User) string {
	if len(user.FirstName) != 0 && len(user.LastName) != 0 {
		return user.FirstName + " " + user.LastName
	} else if len(user.FirstName) != 0 {
		return user.FirstName
	} else if len(user.LastName) != 0 {
		return user.LastName
	}
	return ""
}
