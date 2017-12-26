package sqlStore

import (
	"os"
	"time"

	l4g "github.com/alecthomas/log4go"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

const (
	VERSION_0_0_0            = "0.0.0"
	VERSION_4_0_0            = "4.0.0"
	OLDEST_SUPPORTED_VERSION = "0.0.0"
)

const (
	EXIT_VERSION_SAVE_MISSING = 1001
	EXIT_TOO_OLD              = 1002
	EXIT_VERSION_SAVE         = 1003
	EXIT_THEME_MIGRATION      = 1004
)

func UpgradeDatabase(sqlStore SqlStore) {

	UpgradeDatabaseToVersion1(sqlStore)

	if sqlStore.GetCurrentSchemaVersion() == "" {
		if result := <-sqlStore.System().SaveOrUpdate(&model.System{Name: "Version", Value: model.CurrentVersion}); result.Err != nil {
			l4g.Critical(result.Err.Error())
			time.Sleep(time.Second)
			os.Exit(EXIT_VERSION_SAVE_MISSING)
		}

		l4g.Info(utils.T("store.sql.schema_set.info"), model.CurrentVersion)
	}
	
	if sqlStore.GetCurrentSchemaVersion() != model.CurrentVersion {
		l4g.Critical(utils.T("store.sql.schema_version.critical"), sqlStore.GetCurrentSchemaVersion(), OLDEST_SUPPORTED_VERSION, model.CurrentVersion, OLDEST_SUPPORTED_VERSION)
		time.Sleep(time.Second)
		os.Exit(EXIT_TOO_OLD)
	}
}

func saveSchemaVersion(sqlStore SqlStore, version string) {
	if result := <-sqlStore.System().Update(&model.System{Name: "Version", Value: version}); result.Err != nil {
		l4g.Critical(result.Err.Error())
		time.Sleep(time.Second)
		os.Exit(EXIT_VERSION_SAVE)
	}

	l4g.Warn(utils.T("store.sql.upgraded.warn"), version)
}

func shouldPerformUpgrade(sqlStore SqlStore, currentSchemaVersion string, expectedSchemaVersion string) bool {
	if sqlStore.GetCurrentSchemaVersion() == currentSchemaVersion {
		l4g.Warn(utils.T("store.sql.schema_out_of_date.warn"), currentSchemaVersion)
		l4g.Warn(utils.T("store.sql.schema_upgrade_attempt.warn"), expectedSchemaVersion)

		return true
	}

	return false
}

func UpgradeDatabaseToVersion1(sqlStore SqlStore) {
	if shouldPerformUpgrade(sqlStore, VERSION_0_0_0, VERSION_4_0_0) {
		sqlStore.CreateColumnIfNotExists("OutgoingWebhooks", "ContentType", "varchar(128)", "varchar(128)", "")
		saveSchemaVersion(sqlStore, VERSION_4_0_0)
	}
}
