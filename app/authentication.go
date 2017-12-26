package app

import (
	"net/http"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

func checkUserPassword(user *model.User, password string) *model.AppError {
	if !model.ComparePassword(user.Password, password) {
		if result := <-Srv.SqlStore.User().UpdateFailedPasswordAttempts(user.Id, user.FailedAttempts+1); result.Err != nil {
			return result.Err
		}

		return model.NewAppError("checkUserPassword",
			"api.user.check_user_password.invalid.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized,
		)
	} else {
		if result := <-Srv.SqlStore.User().UpdateFailedPasswordAttempts(user.Id, 0); result.Err != nil {
			return result.Err
		}

		return nil
	}
}

func CheckPasswordAndAllCriteria(user *model.User, password string) *model.AppError {
	if err := CheckUserAdditionalAuthenticationCriteria(user); err != nil {
		return err
	}

	if err := checkUserPassword(user, password); err != nil {
		return err
	}

	return nil
}

func CheckUserAdditionalAuthenticationCriteria(user *model.User) *model.AppError {
	if err := checkUserNotDisabled(user); err != nil {
		return err
	}

	if err := checkUserLoginAttempts(user); err != nil {
		return err
	}

	return nil
}

func checkUserLoginAttempts(user *model.User) *model.AppError {
	if user.FailedAttempts >= utils.Cfg.ServiceSettings.MaximumLoginAttempts {
		return model.NewAppError("checkUserLoginAttempts",
			"api.user.check_user_login_attempts.too_many.app_error", nil,
			"user_id="+user.Id, http.StatusUnauthorized,
		)
	}

	return nil
}

func checkUserNotDisabled(user *model.User) *model.AppError {
	if user.DeleteAt > 0 {
		return model.NewAppError("Login",
			"api.user.login.inactive.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized,
		)
	}

	return nil
}

func authenticateUser(user *model.User, password string) (*model.User, *model.AppError) {
	if err := CheckPasswordAndAllCriteria(user, password); err != nil {
		err.StatusCode = http.StatusUnauthorized
		return user, err
	} else {
		return user, nil
	}
}
