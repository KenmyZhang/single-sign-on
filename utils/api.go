package utils

import (
	"net/http"
	"net/url"

	"github.com/KenmyZhang/single-sign-on/model"
)

func RenderWebError(err *model.AppError, w http.ResponseWriter, r *http.Request) {
	T, _ := GetTranslationsAndLocale(w, r)

	title := T("api.templates.error.title", map[string]interface{}{"SiteName": ClientCfg["SiteName"]})
	message := err.Message
	details := err.DetailedError
	link := "/"
	linkMessage := T("api.templates.error.link")

	status := http.StatusTemporaryRedirect
	if err.StatusCode != http.StatusInternalServerError {
		status = err.StatusCode
	}

	http.Redirect(
		w,
		r,
		"/error?title="+url.QueryEscape(title)+
			"&message="+url.QueryEscape(message)+
			"&details="+url.QueryEscape(details)+
			"&link="+url.QueryEscape(link)+
			"&linkmessage="+url.QueryEscape(linkMessage),
		status,
	)
}
