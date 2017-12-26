package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/spf13/cobra"

	"github.com/KenmyZhang/single-sign-on/api"
	"github.com/KenmyZhang/single-sign-on/app"
	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
	"github.com/KenmyZhang/single-sign-on/web"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "运行服务器",
	RunE:  runServerCmd,
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	config, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	utils.CfgDisableConfigWatch, _ = cmd.Flags().GetBool("disableconfigwatch")

	runServer(config)
	return nil
}

func runServer(configFileLocation string) {
	if err := utils.InitAndLoadConfig(configFileLocation); err != nil {
		l4g.Exit("Unable to load single-sign-on configuration file: ", err)
		return
	}

	if err := utils.InitTranslations(utils.Cfg.LocalizationSettings); err != nil {
		l4g.Exit("Unable to load single-sign-on translation files: %v", err)
		return
	}

	pwd, _ := os.Getwd()
	l4g.Info(utils.T("single-sign-on.current_version"), model.CurrentVersion, model.BuildNumber, model.BuildDate, model.BuildHash)
	l4g.Info(utils.T("single-sign-on.working_dir"), pwd)
	l4g.Info(utils.T("single-sign-on.config_file"), utils.FindConfigFile(configFileLocation))

	if model.BuildNumber == "dev" {
		*utils.Cfg.ServiceSettings.EnableDeveloper = true
	}

	app.NewServer()
	app.InitStores()
	api.InitRouter()
	api.InitApi()
	web.InitWeb()

	app.ReloadConfig()

	app.StartServer()

	go runTokenCleanupJob()

	utils.RegenerateClientConfig()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	app.StopServer()
}

func runTokenCleanupJob() {
	doTokenCleanup()
	model.CreateRecurringTask("Token Cleanup", doTokenCleanup, time.Hour*24)
}

func doTokenCleanup() {
	app.Srv.SqlStore.Token().Cleanup()
}
