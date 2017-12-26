package sqlStore

import (
	"database/sql"
	"fmt"
	"regexp"	
	"net/http"
	"strings"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

const (
	MISSING_ACCOUNT_ERROR               = "store.sql_user.missing_account.const"
	MISSING_AUTH_ACCOUNT_ERROR          = "store.sql_user.get_by_auth.missing_account.app_error"
	USER_SEARCH_TYPE_NAMES_NO_FULL_NAME = "Username, Nickname"
	USER_SEARCH_OPTION_ALLOW_INACTIVE   = "allow_inactive"	
	USER_SEARCH_TYPE_NAMES              = "Username, FirstName, LastName, Nickname"
	USER_SEARCH_TYPE_ALL_NO_FULL_NAME   = "Username, Nickname, Email"
	USER_SEARCH_TYPE_ALL                = "Username, FirstName, LastName, Nickname, Email, Mobile"
)

type SqlUserStore struct {
	SqlStore
}

func NewSqlUserStore(sqlStore SqlStore) UserStore {
	us := &SqlUserStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.User{}, "Users").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("Username").SetMaxSize(64).SetUnique(true)
		table.ColMap("Password").SetMaxSize(128)
		table.ColMap("AuthData").SetMaxSize(128).SetUnique(true)
		table.ColMap("AuthService").SetMaxSize(32)
		table.ColMap("Email").SetMaxSize(128).SetUnique(true)
		table.ColMap("Nickname").SetMaxSize(64)
		table.ColMap("FirstName").SetMaxSize(64)
		table.ColMap("LastName").SetMaxSize(64)
		table.ColMap("Roles").SetMaxSize(64)
		table.ColMap("Props").SetMaxSize(4000)
		table.ColMap("NotifyProps").SetMaxSize(2000)
		table.ColMap("Locale").SetMaxSize(5)
		table.ColMap("MfaSecret").SetMaxSize(128)
		table.ColMap("Position").SetMaxSize(64)
	}

	return us
}

func (us SqlUserStore) CreateIndexesIfNotExists() {
	us.CreateIndexIfNotExists("idx_users_email", "Users", "Email")
	us.CreateIndexIfNotExists("idx_users_update_at", "Users", "UpdateAt")
	us.CreateIndexIfNotExists("idx_users_create_at", "Users", "CreateAt")
	us.CreateIndexIfNotExists("idx_users_delete_at", "Users", "DeleteAt")

	us.CreateFullTextIndexIfNotExists("idx_users_all_type_txt", "Users", USER_SEARCH_TYPE_ALL)
	us.CreateFullTextIndexIfNotExists("idx_users_all_no_full_name_txt", "Users", USER_SEARCH_TYPE_ALL_NO_FULL_NAME)
	us.CreateFullTextIndexIfNotExists("idx_users_names_txt", "Users", USER_SEARCH_TYPE_NAMES)
	us.CreateFullTextIndexIfNotExists("idx_users_names_no_full_name_txt", "Users", USER_SEARCH_TYPE_NAMES_NO_FULL_NAME)
}

func (us SqlUserStore) Save(user *model.User) StoreChannel {

	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		if len(user.Id) > 0 {
			result.Err = model.NewLocAppError("SqlUserStore.Save", "store.sql_user.save.existing.app_error", nil, "user_id="+user.Id)
			storeChannel <- result
			close(storeChannel)
			return
		}

		user.PreSave()
		if result.Err = user.IsValid(); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if err := us.GetMaster().Insert(user); err != nil {
			if IsUniqueConstraintError(err.Error(), []string{"Email", "users_email_key", "idx_users_email_unique"}) {
				result.Err = model.NewAppError("SqlUserStore.Save", "store.sql_user.save.email_exists.app_error", nil, "user_id="+user.Id+", "+err.Error(), http.StatusBadRequest)
			} else if IsUniqueConstraintError(err.Error(), []string{"Username", "users_username_key", "idx_users_username_unique"}) {
				result.Err = model.NewAppError("SqlUserStore.Save", "store.sql_user.save.username_exists.app_error", nil, "user_id="+user.Id+", "+err.Error(), http.StatusBadRequest)
			} else {
				result.Err = model.NewLocAppError("SqlUserStore.Save", "store.sql_user.save.app_error", nil, "user_id="+user.Id+", "+err.Error())
			}
		} else {
			result.Data = user
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) Update(user *model.User, trustedUpdateData bool) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		user.PreUpdate()

		if result.Err = user.IsValid(); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if oldUserResult, err := us.GetMaster().Get(model.User{}, user.Id); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.Update", "store.sql_user.update.finding.app_error", nil, "user_id="+user.Id+", "+err.Error())
		} else if oldUserResult == nil {
			result.Err = model.NewLocAppError("SqlUserStore.Update", "store.sql_user.update.find.app_error", nil, "user_id="+user.Id)
		} else {
			oldUser := oldUserResult.(*model.User)
			user.CreateAt = oldUser.CreateAt
			if user.AuthData == nil || user.AuthService == "" {
				user.AuthData = oldUser.AuthData
				user.AuthService = oldUser.AuthService
			}
			user.Password = oldUser.Password
			user.LastPasswordUpdate = oldUser.LastPasswordUpdate
			user.LastPictureUpdate = oldUser.LastPictureUpdate
			user.EmailVerified = oldUser.EmailVerified
			user.FailedAttempts = oldUser.FailedAttempts
			user.MfaSecret = oldUser.MfaSecret
			user.MfaActive = oldUser.MfaActive

			if !trustedUpdateData {
				user.Roles = oldUser.Roles
				user.DeleteAt = oldUser.DeleteAt
			}

			if user.IsOAuthUser() {
				if !trustedUpdateData {
					user.Email = oldUser.Email
				}
			} else if user.Email != oldUser.Email {
				user.EmailVerified = false
			}

			if user.Username != oldUser.Username {
				user.UpdateMentionKeysFromUsername(oldUser.Username)
			}

			if count, err := us.GetMaster().Update(user); err != nil {
				if IsUniqueConstraintError(err.Error(), []string{"Email", "users_email_key", "idx_users_email_unique"}) {
					result.Err = model.NewAppError("SqlUserStore.Update", "store.sql_user.update.email_taken.app_error", nil, "user_id="+user.Id+", "+err.Error(), http.StatusBadRequest)
				} else if IsUniqueConstraintError(err.Error(), []string{"Username", "users_username_key", "idx_users_username_unique"}) {
					result.Err = model.NewAppError("SqlUserStore.Update", "store.sql_user.update.username_taken.app_error", nil, "user_id="+user.Id+", "+err.Error(), http.StatusBadRequest)
				} else {
					result.Err = model.NewLocAppError("SqlUserStore.Update", "store.sql_user.update.updating.app_error", nil, "user_id="+user.Id+", "+err.Error())
				}
			} else if count != 1 {
				result.Err = model.NewLocAppError("SqlUserStore.Update", "store.sql_user.update.app_error", nil, fmt.Sprintf("user_id=%v, count=%v", user.Id, count))
			} else {
				user.Sanitize(map[string]bool{})
				oldUser.Sanitize(map[string]bool{})
				result.Data = [2]*model.User{user, oldUser}
			}
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) UpdateLastPictureUpdate(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		curTime := model.GetMillis()

		if _, err := us.GetMaster().Exec("UPDATE Users SET LastPictureUpdate = :Time, UpdateAt = :Time WHERE Id = :UserId", map[string]interface{}{"Time": curTime, "UserId": userId}); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.UpdateUpdateAt", "store.sql_user.update_last_picture_update.app_error", nil, "user_id="+userId)
		} else {
			result.Data = curTime
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) Get(id string) StoreChannel {

	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		if obj, err := us.GetReplica().Get(model.User{}, id); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.Get", "store.sql_user.get.app_error", nil, "user_id="+id+", "+err.Error())
		} else if obj == nil {
			result.Err = model.NewAppError("SqlUserStore.Get", MISSING_ACCOUNT_ERROR, nil, "user_id="+id, http.StatusNotFound)
		} else {
			result.Data = obj.(*model.User)
		}

		storeChannel <- result
		close(storeChannel)

	}()

	return storeChannel
}

func (us SqlUserStore) GetByEmail(email string) StoreChannel {

	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		email = strings.ToLower(email)

		user := model.User{}

		if err := us.GetReplica().SelectOne(&user, "SELECT * FROM Users WHERE Email = :Email", map[string]interface{}{"Email": email}); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.GetByEmail", MISSING_ACCOUNT_ERROR, nil, "email="+email+", "+err.Error())
		}

		result.Data = &user

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) GetByAuth(authData *string, authService string) StoreChannel {

	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		if authData == nil || *authData == "" {
			result.Err = model.NewLocAppError("SqlUserStore.GetByAuth", MISSING_AUTH_ACCOUNT_ERROR, nil, "authData='', authService="+authService)
			storeChannel <- result
			close(storeChannel)
			return
		}

		user := model.User{}

		if err := us.GetReplica().SelectOne(&user, "SELECT * FROM Users WHERE AuthData = :AuthData AND AuthService = :AuthService", map[string]interface{}{"AuthData": authData, "AuthService": authService}); err != nil {
			if err == sql.ErrNoRows {
				result.Err = model.NewLocAppError("SqlUserStore.GetByAuth", MISSING_AUTH_ACCOUNT_ERROR, nil, "authData="+*authData+", authService="+authService+", "+err.Error())
			} else {
				result.Err = model.NewLocAppError("SqlUserStore.GetByAuth", "store.sql_user.get_by_auth.other.app_error", nil, "authData="+*authData+", authService="+authService+", "+err.Error())
			}
		}

		result.Data = &user

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) GetByUsername(username string) StoreChannel {

	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		user := model.User{}

		if err := us.GetReplica().SelectOne(&user, "SELECT * FROM Users WHERE Username = :Username", map[string]interface{}{"Username": username}); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.GetByUsername", "store.sql_user.get_by_username.app_error",
				nil, err.Error())
		}

		result.Data = &user

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail, allowSignInWithMobile, weixinEnabled bool) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		params := map[string]interface{}{
			"LoginId":                 loginId,
			"AllowSignInWithUsername": allowSignInWithUsername,
			"AllowSignInWithEmail":    allowSignInWithEmail,
			"AllowSignInWithMobile": allowSignInWithMobile,
			"WeixinEnabled":           weixinEnabled,
		}

		users := []*model.User{}
		if _, err := us.GetReplica().Select(
			&users,
			`SELECT
				*
			FROM
				Users
			WHERE
				(:AllowSignInWithUsername AND Username = :LoginId)
				OR (:AllowSignInWithEmail AND Email = :LoginId)
				OR (:AllowSignInWithMobile AND Mobile = :LoginId)
				OR (:WeixinEnabled AND AuthService = '`+model.USER_AUTH_SERVICE_WECHAT+`' AND AuthData = :LoginId)`,
			params); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.GetForLogin", "store.sql_user.get_for_login.app_error", nil, err.Error())
		} else if len(users) == 1 {
			result.Data = users[0]
		} else if len(users) > 1 {
			result.Err = model.NewLocAppError("SqlUserStore.GetForLogin", "store.sql_user.get_for_login.multiple_users", nil, "")
		} else {
			result.Err = model.NewLocAppError("SqlUserStore.GetForLogin", "store.sql_user.get_for_login.app_error", nil, "")
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) GetTotalUsersCount() StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		if count, err := us.GetReplica().SelectInt("SELECT COUNT(Id) FROM Users"); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.GetTotalUsersCount", "store.sql_user.get_total_users_count.app_error", nil, err.Error())
		} else {
			result.Data = count
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) UpdateFailedPasswordAttempts(userId string, attempts int) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		if _, err := us.GetMaster().Exec("UPDATE Users SET FailedAttempts = :FailedAttempts WHERE Id = :UserId", map[string]interface{}{"FailedAttempts": attempts, "UserId": userId}); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.UpdateFailedPasswordAttempts", "store.sql_user.update_failed_pwd_attempts.app_error", nil, "user_id="+userId)
		} else {
			result.Data = userId
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us SqlUserStore) SearchUsers(term string, options map[string]bool) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		searchQuery := `
		SELECT
			*
		FROM
			Users
		WHERE
			SEARCH_CLAUSE
			INACTIVE_CLAUSE
			ORDER BY Username ASC
		LIMIT 100`

		storeChannel <- us.performSearch(searchQuery, term, options, map[string]interface{}{})
		close(storeChannel)

	}()

	return storeChannel
}

var specialUserSearchChar = []string{
	"<",
	">",
	"+",
	"-",
	"(",
	")",
	"~",
	":",
	"*",
	"\"",
	"!",
	"@",
}

func (us SqlUserStore) performSearch(searchQuery string, term string, options map[string]bool, parameters map[string]interface{}) StoreResult {
	result := StoreResult{}

	// Special handling for emails
	if strings.Contains(term, "@") && strings.Contains(term, ".") {
		if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
			lastIndex := strings.LastIndex(term, ".")
			term = term[0:lastIndex]
		}
	}

	// these chars have special meaning and can be treated as spaces
	for _, c := range specialUserSearchChar {
		term = strings.Replace(term, c, " ", -1)
	}

	var matchTerm = term

	searchType := USER_SEARCH_TYPE_ALL

	if ok := options[USER_SEARCH_OPTION_ALLOW_INACTIVE]; ok {
		searchQuery = strings.Replace(searchQuery, "INACTIVE_CLAUSE", "", 1)
	} else {
		searchQuery = strings.Replace(searchQuery, "INACTIVE_CLAUSE", "AND Users.DeleteAt = 0", 1)
	}

	if term == "" {
		searchQuery = strings.Replace(searchQuery, "SEARCH_CLAUSE", "", 1)
	} else if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
		splitTerm := strings.Fields(term)
		for i, t := range strings.Fields(term) {
			splitTerm[i] = "+" + t + "*"
		}

		term = strings.Join(splitTerm, " ")

		searchClause := fmt.Sprintf("(MATCH(%s) AGAINST (:Term IN BOOLEAN MODE)) or Nickname like :MatchTerm", searchType)
		searchQuery = strings.Replace(searchQuery, "SEARCH_CLAUSE", searchClause, 1)
	}

	var users []*model.User

	parameters["Term"] = term
	parameters["MatchTerm"] = matchTerm + "%"

	if _, err := us.GetReplica().Select(&users, searchQuery, parameters); err != nil {
		result.Err = model.NewLocAppError("SqlUserStore.Search", "store.sql_user.search.app_error", nil, "term="+term+", "+"search_type="+searchType+", "+err.Error())
	} else {
		for _, u := range users {
			u.Sanitize(map[string]bool{})
		}

		result.Data = users
	}

	return result
}


var mobileChar = regexp.MustCompile(`^[0-9]+$`)

func (us SqlUserStore) GetProfileByMobile(mobile string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	result := StoreResult{}
	
	if !mobileChar.MatchString(mobile) {
		result.Err = model.NewLocAppError("SqlUserStore.GetProfileByMobile", "store.sql_user.get_profile_by_mobile.invalid.app_error", nil, "Mobile="+mobile)
		storeChannel <- result
		close(storeChannel)
		return storeChannel	
	}

	go func() {

		user := &model.User{}
		if err := us.GetReplica().SelectOne(&user, "SELECT * FROM Users WHERE Mobile = :Mobile", map[string]interface{}{"Mobile": mobile}); err != nil {
			if err == sql.ErrNoRows {
				result.Err = model.NewLocAppError("SqlUserStore.GetProfileByMobile", "store.sql_user.missing_account.const", nil, "Mobile="+mobile+" does not exist, "+err.Error())
			} else {
				result.Err = model.NewLocAppError("SqlUserStore.GetProfileByMobile", "store.sql_user.get_profile_by_mobile.other.app_error", nil, "Mobile="+mobile+", "+err.Error())
			}
		}
		result.Data = user

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel	
}

func (us SqlUserStore) UpdatePassword(userId, hashedPassword string) StoreChannel {

	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		updateAt := model.GetMillis()
		if _, err := us.GetMaster().Exec("UPDATE Users SET Password = :Password, LastPasswordUpdate = :LastPasswordUpdate, UpdateAt = :UpdateAt, EmailVerified = true, FailedAttempts = 0 WHERE Id = :UserId", map[string]interface{}{"Password": hashedPassword, "LastPasswordUpdate": updateAt, "UpdateAt": updateAt, "UserId": userId}); err != nil {
			result.Err = model.NewLocAppError("SqlUserStore.UpdatePassword", "store.sql_user.update_password.app_error", nil, "id="+userId+", "+err.Error())
		} else {
			result.Data = userId
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}