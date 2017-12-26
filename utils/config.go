package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	l4g "github.com/alecthomas/log4go"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/KenmyZhang/single-sign-on/model"
)

const (
	MODE_DEV        = "dev"
	MODE_BETA       = "beta"
	MODE_PROD       = "prod"
	LOG_ROTATE_SIZE = 10000
	LOG_FILENAME    = "sso.log"
)

var cfgMutex = &sync.Mutex{}
var watcher *fsnotify.Watcher
var Cfg *model.Config = &model.Config{}
var CfgHash = ""
var ClientCfgHash = ""
var CfgFileName string = ""
var CfgDisableConfigWatch = false
var ClientCfg map[string]string = map[string]string{}
var originalDisableDebugLvl l4g.Level = l4g.DEBUG
var siteURL = ""

func GetSiteURL() string {
	return siteURL
}

func SetSiteURL(url string) {
	siteURL = strings.TrimRight(url, "/")
}

var cfgListeners = map[string]func(*model.Config, *model.Config){}

func AddConfigListener(listener func(*model.Config, *model.Config)) string {
	id := model.NewId()
	cfgListeners[id] = listener
	return id
}

func RemoveConfigListener(id string) {
	delete(cfgListeners, id)
}

func FindConfigFile(fileName string) string {
	if _, err := os.Stat("./config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("./config/" + fileName)
	} else if _, err := os.Stat("../config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("../config/" + fileName)
	} else if _, err := os.Stat(fileName); err == nil {
		fileName, _ = filepath.Abs(fileName)
	}

	return fileName
}

func FindDir(dir string) (string, bool) {
	fileName := "."
	found := false
	if _, err := os.Stat("./" + dir + "/"); err == nil {
		fileName, _ = filepath.Abs("./" + dir + "/")
		found = true
	} else if _, err := os.Stat("../" + dir + "/"); err == nil {
		fileName, _ = filepath.Abs("../" + dir + "/")
		found = true
	} else if _, err := os.Stat("../../" + dir + "/"); err == nil {
		fileName, _ = filepath.Abs("../../" + dir + "/")
		found = true
	}

	return fileName + "/", found
}

func DisableDebugLogForTest() {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()
	if l4g.Global["stdout"] != nil {
		originalDisableDebugLvl = l4g.Global["stdout"].Level
		l4g.Global["stdout"].Level = l4g.ERROR
	}
}

func EnableDebugLogForTest() {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()
	if l4g.Global["stdout"] != nil {
		l4g.Global["stdout"].Level = originalDisableDebugLvl
	}
}

func configureLog(s *model.LogSettings) {
	l4g.Close()

	if s.EnableConsole {
		level := l4g.DEBUG
		if s.ConsoleLevel == "INFO" {
			level = l4g.INFO
		} else if s.ConsoleLevel == "WARN" {
			level = l4g.WARNING
		} else if s.ConsoleLevel == "ERROR" {
			level = l4g.ERROR
		}

		lw := l4g.NewConsoleLogWriter()
		lw.SetFormat("[%D %T] [%L] %M")
		l4g.AddFilter("stdout", level, lw)
	}

	if s.EnableFile {
		var fileFormat = s.FileFormat

		if fileFormat == "" {
			fileFormat = "[%D %T] [%L] %M"
		}

		level := l4g.DEBUG
		if s.FileLevel == "INFO" {
			level = l4g.INFO
		} else if s.FileLevel == "WARN" {
			level = l4g.WARNING
		} else if s.FileLevel == "ERROR" {
			level = l4g.ERROR
		}

		flw := l4g.NewFileLogWriter(GetLogFileLocation(s.FileLocation), false)
		flw.SetFormat(fileFormat)
		flw.SetRotate(true)
		flw.SetRotateLines(LOG_ROTATE_SIZE)
		l4g.AddFilter("file", level, flw)
	}
}

func GetLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		logDir, _ := FindDir("logs")
		return logDir + LOG_FILENAME
	} else {
		return fileLocation + LOG_FILENAME
	}
}

func SaveConfig(fileName string, config *model.Config) *model.AppError {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	b, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return model.NewLocAppError("SaveConfig", "utils.config.save_config.saving.app_error",
			map[string]interface{}{"Filename": fileName}, err.Error())
	}

	err = ioutil.WriteFile(fileName, b, 0644)
	if err != nil {
		return model.NewLocAppError("SaveConfig", "utils.config.save_config.saving.app_error",
			map[string]interface{}{"Filename": fileName}, err.Error())
	}

	return nil
}

func EnableConfigFromEnviromentVars() {
	viper.SetEnvPrefix("sso")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func InitializeConfigWatch() {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	if CfgDisableConfigWatch {
		return
	}

	if watcher == nil {
		var err error
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			l4g.Error(fmt.Sprintf("Failed to watch config file at %v with err=%v", CfgFileName, err.Error()))
		}

		go func() {
			configFile := filepath.Clean(CfgFileName)

			for {
				select {
				case event := <-watcher.Events:
					// we only care about the config file
					if filepath.Clean(event.Name) == configFile {
						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
							l4g.Info(fmt.Sprintf("Config file watcher detected a change reloading %v", CfgFileName))

							if configReadErr := viper.ReadInConfig(); configReadErr == nil {
								LoadConfig(CfgFileName)
							} else {
								l4g.Error(fmt.Sprintf("Failed to read while watching config file at %v with err=%v", CfgFileName, configReadErr.Error()))
							}
						}
					}
				case err := <-watcher.Errors:
					l4g.Error(fmt.Sprintf("Failed while watching config file at %v with err=%v", CfgFileName, err.Error()))
				}
			}
		}()
	}
}

func EnableConfigWatch() {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	if watcher != nil {
		configFile := filepath.Clean(CfgFileName)
		configDir, _ := filepath.Split(configFile)

		if watcher != nil {
			watcher.Add(configDir)
		}
	}
}

func DisableConfigWatch() {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	if watcher != nil {
		configFile := filepath.Clean(CfgFileName)
		configDir, _ := filepath.Split(configFile)
		watcher.Remove(configDir)
	}
}

func InitAndLoadConfig(filename string) error {
	if err := TranslationsPreInit(); err != nil {
		return err
	}

	EnableConfigFromEnviromentVars()
	LoadConfig(filename)
	InitializeConfigWatch()
	EnableConfigWatch()

	return nil
}

func LoadConfig(fileName string) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	oldConfig := *Cfg

	fileNameWithExtension := filepath.Base(fileName)
	fileExtension := filepath.Ext(fileNameWithExtension)
	fileDir := filepath.Dir(fileName)

	if len(fileNameWithExtension) > 0 {
		fileNameOnly := fileNameWithExtension[:len(fileNameWithExtension)-len(fileExtension)]
		viper.SetConfigName(fileNameOnly)
	} else {
		viper.SetConfigName("config")
	}

	if len(fileDir) > 0 {
		viper.AddConfigPath(fileDir)
	}

	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath(".")

	configReadErr := viper.ReadInConfig()
	if configReadErr != nil {
		errMsg := T("utils.config.load_config.opening.panic", map[string]interface{}{"Filename": fileName, "Error": configReadErr.Error()})
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(1)
	}

	var config model.Config
	unmarshalErr := viper.Unmarshal(&config)
	if unmarshalErr != nil {
		errMsg := T("utils.config.load_config.decoding.panic", map[string]interface{}{"Filename": fileName, "Error": unmarshalErr.Error()})
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(1)
	}

	CfgFileName = viper.ConfigFileUsed()

	needSave := false

	config.SetDefaults()

	if err := config.IsValid(); err != nil {
		panic(T(err.Id))
	}

	if needSave {
		cfgMutex.Unlock()
		if err := SaveConfig(CfgFileName, &config); err != nil {
			l4g.Warn(T(err.Id))
		}
		cfgMutex.Lock()
	}

	if err := ValidateLocales(&config); err != nil {
		panic(T(err.Id))
	}

	configureLog(&config.LogSettings)

	Cfg = &config
	cfgJson, _ := json.Marshal(Cfg)
	CfgHash = fmt.Sprintf("%x", md5.Sum(cfgJson))
	ClientCfg = getClientConfig(Cfg)
	clientCfgJson, _ := json.Marshal(ClientCfg)
	ClientCfgHash = fmt.Sprintf("%x", md5.Sum(clientCfgJson))

	SetDefaultRolesBasedOnConfig()
	SetSiteURL(*Cfg.ServiceSettings.SiteURL)

	for _, listener := range cfgListeners {
		listener(&oldConfig, &config)
	}
}

func RegenerateClientConfig() {
	ClientCfg = getClientConfig(Cfg)
}

func getClientConfig(c *model.Config) map[string]string {
	props := make(map[string]string)

	props["Version"] = model.CurrentVersion
	props["BuildNumber"] = model.BuildNumber
	props["BuildDate"] = model.BuildDate
	props["BuildHash"] = model.BuildHash

	props["SiteURL"] = strings.TrimRight(*c.ServiceSettings.SiteURL, "/")
	props["SiteName"] = c.ServiceSettings.SiteName

	props["EnableSignInWithMobile"] = strconv.FormatBool(*c.EnableSignInWithMobile)

	props["EnableSignUpWithWeixin"] = strconv.FormatBool(c.WeixinSettings.Enable)

	props["DefaultClientLocale"] = *c.LocalizationSettings.DefaultClientLocale
	props["AvailableLocales"] = *c.LocalizationSettings.AvailableLocales

	return props
}

func ValidateLocales(cfg *model.Config) *model.AppError {
	locales := GetSupportedLocales()
	if _, ok := locales[*cfg.LocalizationSettings.DefaultServerLocale]; !ok {
		return model.NewLocAppError("ValidateLocales", "utils.config.supported_server_locale.app_error", nil, "")
	}

	if _, ok := locales[*cfg.LocalizationSettings.DefaultClientLocale]; !ok {
		return model.NewLocAppError("ValidateLocales", "utils.config.supported_client_locale.app_error", nil, "")
	}

	if len(*cfg.LocalizationSettings.AvailableLocales) > 0 {
		for _, word := range strings.Split(*cfg.LocalizationSettings.AvailableLocales, ",") {
			if word == *cfg.LocalizationSettings.DefaultClientLocale {
				return nil
			}
		}

		return model.NewLocAppError("ValidateLocales", "utils.config.validate_locale.app_error", nil, "")
	}

	return nil
}
