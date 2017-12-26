package model

import (
	"net/url"
)

const (
	DATABASE_DRIVER_MYSQL = "mysql"

	CONN_SECURITY_NONE     = ""
	CONN_SECURITY_PLAIN    = "PLAIN"
	CONN_SECURITY_TLS      = "TLS"
	CONN_SECURITY_STARTTLS = "STARTTLS"

	PASSWORD_MAXIMUM_LENGTH = 64
	PASSWORD_MINIMUM_LENGTH = 5

	SERVICE_WEIXIN = "weixin"

	WEBSERVER_MODE_REGULAR  = "regular"
	WEBSERVER_MODE_GZIP     = "gzip"
	WEBSERVER_MODE_DISABLED = "disabled"

	FAKE_SETTING = "********************************"

	SITENAME_MAX_LENGTH = 30

	SERVICE_SETTINGS_DEFAULT_SITE_URL        = ""
	SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE   = ""
	SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE    = ""
	SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT    = 300
	SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT   = 300
	SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM = ""
)

type SSOSettings struct {
	Enable          bool
	Secret          string
	Id              string
	Scope           string
	AuthEndpoint    string
	TokenEndpoint   string
	UserApiEndpoint string
}

type JWTSettings struct {
	Secret                 string
	Issuer                 string
	ExpireTimeLengthInDays *int
}

type ServiceSettings struct {
	SiteURL                             *string
	SiteName                            string
	ListenAddress                       string
	ConnectionSecurity                  *string
	TLSCertFile                         *string
	TLSKeyFile                          *string
	Forward80To443                      *bool
	ReadTimeout                         *int
	WriteTimeout                        *int
	MaximumLoginAttempts                int
	EnableDeveloper                     *bool
	AllowCorsFrom                       *string
	SessionLengthWebInDays              *int
	SessionCacheInMinutes               *int
	WebserverMode                       *string
	EnableInsecureOutgoingConnections   *bool
	AllowedUntrustedInternalConnections *string
}

type SqlSettings struct {
	DriverName               string
	DataSource               string
	DataSourceReplicas       []string
	DataSourceSearchReplicas []string
	MaxIdleConns             int
	MaxOpenConns             int
	Trace                    bool
	AtRestEncryptKey         string
	QueryTimeout             *int
}

type LogSettings struct {
	EnableConsole bool
	ConsoleLevel  string
	EnableFile    bool
	FileLevel     string
	FileFormat    string
	FileLocation  string
}

type EmailSettings struct {
	EnableSignUpWithEmail    bool
	EnableSignInWithEmail    *bool
	EnableSignInWithUsername *bool
	FeedbackName                      string
	FeedbackEmail                     string
	FeedbackOrganization              *string
	EnableSMTPAuth                    *bool
	SMTPUsername                      string
	SMTPPassword                      string	
	SMTPServer                        string
	SMTPPort                          string
	ConnectionSecurity                string	
	SkipServerCertificateVerification *bool		
}

type RateLimitSettings struct {
	Enable           *bool
	PerSec           int
	MaxBurst         *int
	MemoryStoreSize  int
	VaryByRemoteAddr bool
	VaryByHeader     string
}

type PrivacySettings struct {
	ShowEmailAddress bool
	ShowFullName     bool
}

type PasswordSettings struct {
	MinimumLength *int
	Lowercase     *bool
	Number        *bool
	Uppercase     *bool
	Symbol        *bool
}

type LocalizationSettings struct {
	DefaultServerLocale *string
	DefaultClientLocale *string
	AvailableLocales    *string
}

type SmsSettings struct {
	Enable     		 bool
	Provider         string
	GatewayUrl       string
	AccessKeyId    	 string
	AccessKeySecret	 string
	SignName         string
	TemplateCode     string
}

type SupportSettings struct {
	TermsOfServiceLink       *string
	PrivacyPolicyLink        *string
	AboutLink                *string
	HelpLink                 *string
	ReportAProblemLink       *string
	AdministratorsGuideLink  *string
	TroubleshootingForumLink *string
	CommercialSupportLink    *string
	SupportEmail             *string
}

type FileSettings struct {
	EnableFileAttachments *bool
	MaxFileSize           *int64
	InitialFont 		  *string
	Directory 			  *string
}

type Config struct {
	SqlSettings          SqlSettings
	ServiceSettings      ServiceSettings
	SupportSettings       SupportSettings
	LogSettings          LogSettings
	EmailSettings        EmailSettings
	EnableSignInWithMobile *bool
	SmsCodeSettings		 SmsSettings	
	FileSettings         FileSettings
	RateLimitSettings    RateLimitSettings
	PrivacySettings      PrivacySettings
	PasswordSettings     PasswordSettings
	WeixinSettings       SSOSettings
	JWTSettings          JWTSettings
	LocalizationSettings LocalizationSettings
}

func (o *Config) GetSSOService(service string) *SSOSettings {
	switch service {
	case SERVICE_WEIXIN:
		return &o.WeixinSettings
	}

	return nil
}

func (o *Config) GetJWTService() *JWTSettings {
	return &o.JWTSettings
}

func (o *Config) GetSmsService() *SmsSettings {
	return &o.SmsCodeSettings
}

func (o *Config) SetDefaults() {
	if o.ServiceSettings.SiteURL == nil {
		o.ServiceSettings.SiteURL = new(string)
		*o.ServiceSettings.SiteURL = SERVICE_SETTINGS_DEFAULT_SITE_URL
	}

	if o.ServiceSettings.EnableDeveloper == nil {
		o.ServiceSettings.EnableDeveloper = new(bool)
		*o.ServiceSettings.EnableDeveloper = false
	}

	if o.PasswordSettings.MinimumLength == nil {
		o.PasswordSettings.MinimumLength = new(int)
		*o.PasswordSettings.MinimumLength = PASSWORD_MINIMUM_LENGTH
	}

	if o.PasswordSettings.Lowercase == nil {
		o.PasswordSettings.Lowercase = new(bool)
		*o.PasswordSettings.Lowercase = false
	}

	if o.PasswordSettings.Number == nil {
		o.PasswordSettings.Number = new(bool)
		*o.PasswordSettings.Number = false
	}

	if o.PasswordSettings.Uppercase == nil {
		o.PasswordSettings.Uppercase = new(bool)
		*o.PasswordSettings.Uppercase = false
	}

	if o.PasswordSettings.Symbol == nil {
		o.PasswordSettings.Symbol = new(bool)
		*o.PasswordSettings.Symbol = false
	}

	if o.ServiceSettings.SessionLengthWebInDays == nil {
		o.ServiceSettings.SessionLengthWebInDays = new(int)
		*o.ServiceSettings.SessionLengthWebInDays = 30
	}

	if o.ServiceSettings.SessionCacheInMinutes == nil {
		o.ServiceSettings.SessionCacheInMinutes = new(int)
		*o.ServiceSettings.SessionCacheInMinutes = 10
	}

	if o.ServiceSettings.AllowCorsFrom == nil {
		o.ServiceSettings.AllowCorsFrom = new(string)
		*o.ServiceSettings.AllowCorsFrom = SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM
	}

	if o.ServiceSettings.WebserverMode == nil {
		o.ServiceSettings.WebserverMode = new(string)
		*o.ServiceSettings.WebserverMode = "gzip"
	} else if *o.ServiceSettings.WebserverMode == "regular" {
		*o.ServiceSettings.WebserverMode = "gzip"
	}

	if o.LocalizationSettings.DefaultServerLocale == nil {
		o.LocalizationSettings.DefaultServerLocale = new(string)
		*o.LocalizationSettings.DefaultServerLocale = DEFAULT_LOCALE
	}

	if o.LocalizationSettings.DefaultClientLocale == nil {
		o.LocalizationSettings.DefaultClientLocale = new(string)
		*o.LocalizationSettings.DefaultClientLocale = DEFAULT_LOCALE
	}

	if o.LocalizationSettings.AvailableLocales == nil {
		o.LocalizationSettings.AvailableLocales = new(string)
		*o.LocalizationSettings.AvailableLocales = ""
	}

	if o.EmailSettings.EnableSignInWithEmail == nil {
		o.EmailSettings.EnableSignInWithEmail = new(bool)

		if o.EmailSettings.EnableSignUpWithEmail == true {
			*o.EmailSettings.EnableSignInWithEmail = true
		} else {
			*o.EmailSettings.EnableSignInWithEmail = false
		}
	}

	if o.EnableSignInWithMobile == nil {
		o.EnableSignInWithMobile = new(bool)

		*o.EnableSignInWithMobile = false
	}

	if o.EmailSettings.EnableSignInWithUsername == nil {
		o.EmailSettings.EnableSignInWithUsername = new(bool)
		*o.EmailSettings.EnableSignInWithUsername = false
	}

	if o.RateLimitSettings.Enable == nil {
		o.RateLimitSettings.Enable = new(bool)
		*o.RateLimitSettings.Enable = false
	}

	if o.RateLimitSettings.MaxBurst == nil {
		o.RateLimitSettings.MaxBurst = new(int)
		*o.RateLimitSettings.MaxBurst = 100
	}

	if o.ServiceSettings.ConnectionSecurity == nil {
		o.ServiceSettings.ConnectionSecurity = new(string)
		*o.ServiceSettings.ConnectionSecurity = ""
	}

	if o.ServiceSettings.TLSKeyFile == nil {
		o.ServiceSettings.TLSKeyFile = new(string)
		*o.ServiceSettings.TLSKeyFile = SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE
	}

	if o.ServiceSettings.TLSCertFile == nil {
		o.ServiceSettings.TLSCertFile = new(string)
		*o.ServiceSettings.TLSCertFile = SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE
	}

	if o.ServiceSettings.ReadTimeout == nil {
		o.ServiceSettings.ReadTimeout = new(int)
		*o.ServiceSettings.ReadTimeout = SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT
	}

	if o.ServiceSettings.WriteTimeout == nil {
		o.ServiceSettings.WriteTimeout = new(int)
		*o.ServiceSettings.WriteTimeout = SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT
	}

	if o.ServiceSettings.Forward80To443 == nil {
		o.ServiceSettings.Forward80To443 = new(bool)
		*o.ServiceSettings.Forward80To443 = false
	}

}

func (o *Config) IsValid() *AppError {

	if o.ServiceSettings.MaximumLoginAttempts <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.login_attempts.app_error", nil, "")
	}

	if len(*o.ServiceSettings.SiteURL) != 0 {
		if _, err := url.ParseRequestURI(*o.ServiceSettings.SiteURL); err != nil {
			return NewLocAppError("Config.IsValid", "model.config.is_valid.site_url.app_error", nil, "")
		}
	}

	if len(o.ServiceSettings.ListenAddress) == 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.listen_address.app_error", nil, "")
	}

	if o.RateLimitSettings.MemoryStoreSize <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.rate_mem.app_error", nil, "")
	}

	if o.RateLimitSettings.PerSec <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.rate_sec.app_error", nil, "")
	}

	if *o.RateLimitSettings.MaxBurst <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.max_burst.app_error", nil, "")
	}

	if *o.PasswordSettings.MinimumLength < PASSWORD_MINIMUM_LENGTH || *o.PasswordSettings.MinimumLength > PASSWORD_MAXIMUM_LENGTH {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.password_length.app_error", map[string]interface{}{"MinLength": PASSWORD_MINIMUM_LENGTH, "MaxLength": PASSWORD_MAXIMUM_LENGTH}, "")
	}

	if len(o.ServiceSettings.SiteName) > SITENAME_MAX_LENGTH {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.sitename_length.app_error", map[string]interface{}{"MaxLength": SITENAME_MAX_LENGTH}, "")
	}

	if !(*o.ServiceSettings.ConnectionSecurity == CONN_SECURITY_NONE || *o.ServiceSettings.ConnectionSecurity == CONN_SECURITY_TLS) {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.webserver_security.app_error", nil, "")
	}

	if *o.ServiceSettings.ReadTimeout <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.read_timeout.app_error", nil, "")
	}

	if *o.ServiceSettings.WriteTimeout <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.write_timeout.app_error", nil, "")
	}

	return nil
}

func (o *Config) GetSanitizeOptions() map[string]bool {
	options := map[string]bool{}
	options["fullname"] = o.PrivacySettings.ShowFullName
	options["email"] = o.PrivacySettings.ShowEmailAddress

	return options
}

func (o *Config) Sanitize() {
	if len(o.WeixinSettings.Secret) > 0 {
		o.WeixinSettings.Secret = FAKE_SETTING
	}

	if len(o.SmsCodeSettings.AccessKeySecret) > 0 {
		o.SmsCodeSettings.AccessKeySecret = FAKE_SETTING
	}		
}
