package app

import (
	"github.com/KenmyZhang/single-sign-on/model"
	"net/http"
)


func GetTokenByExtra(extra string) (*model.Token, *model.AppError) {
	if result := <-Srv.SqlStore.Token().GetByExtra(extra); result.Err != nil {
		return nil, model.NewAppError("GetTokenByExtra", "api.user.get_token_by_extra.app_error", nil, result.Err.Error(), http.StatusBadRequest)
	} else {
		return result.Data.(*model.Token), nil
	}
}

func GetTokenCountByExtra(extra string) (int64, *model.AppError) {
	if result := <-Srv.SqlStore.Token().GetTokenCountByExtra(extra); result.Err != nil {
		return 0, model.NewAppError("GetTokenCountByExtra", "api.user.get_token_count_by_extra.app_error", nil, result.Err.Error(), http.StatusBadRequest)
	} else {
		return result.Data.(int64), nil
	}
}