package app

import (
	"runtime/debug"

	"github.com/KenmyZhang/single-sign-on/utils"
)

func ReloadConfig() {
	debug.FreeOSMemory()
	utils.LoadConfig(utils.CfgFileName)
}
