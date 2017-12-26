package app  

import (
	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)


func SendVerificationCodeToEmail(userEmail, locale, siteURL, verifiedCode string) *model.AppError {
	T := utils.GetUserTranslations(locale)

	subject := T("api.templates.verify_subject",
		map[string]interface{}{"SiteName": utils.ClientCfg["SiteName"]})

	bodyPage := utils.NewHTMLTemplate("verify_body", locale)
	bodyPage.Props["SiteURL"] = siteURL
	bodyPage.Props["Title"] = T("api.templates.email_verified_code_body.title", map[string]interface{}{"SiteName": utils.ClientCfg["SiteName"]})
	bodyPage.Props["Info"] = T("api.templates.email_verified_code_body.info")
	bodyPage.Props["Button"] = verifiedCode
	bodyPage.Props["ListenAddress"] =  utils.Cfg.ServiceSettings.ListenAddress

	if err := utils.SendMail(userEmail, subject, bodyPage.Render()); err != nil {
		return model.NewLocAppError("SendVerifyEmail", "api.user.send_verify_email_and_forget.failed.error", nil, err.Error())
	}

	return nil
}