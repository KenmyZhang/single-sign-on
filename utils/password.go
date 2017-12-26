package utils

import (
	"net/http"
	"strings"

	"github.com/KenmyZhang/single-sign-on/model"
)

func IsPasswordValid(password string) *model.AppError {
	id := "model.user.is_valid.pwd"
	isError := false

	if len(password) < *Cfg.PasswordSettings.MinimumLength || len(password) > model.PASSWORD_MAXIMUM_LENGTH {
		isError = true
	}

	if *Cfg.PasswordSettings.Lowercase {
		if !strings.ContainsAny(password, model.LOWERCASE_LETTERS) {
			isError = true
		}
		id = id + "_lowercase"
	}

	if *Cfg.PasswordSettings.Uppercase {
		if !strings.ContainsAny(password, model.UPPERCASE_LETTERS) {
			isError = true
		}
		id = id + "_uppercase"
	}

	if *Cfg.PasswordSettings.Number {
		if !strings.ContainsAny(password, model.NUMBERS) {
			isError = true
		}
		id = id + "_number"
	}

	if *Cfg.PasswordSettings.Symbol {
		if !strings.ContainsAny(password, model.SYMBOLS) {
			isError = true
		}
		id = id + "_symbol"
	}

	min := *Cfg.PasswordSettings.MinimumLength	

	if isError {
		return model.NewAppError("User.IsValid",
			id+".app_error", map[string]interface{}{"Min": min}, "", http.StatusBadRequest,
		)
	}

	return nil
}
