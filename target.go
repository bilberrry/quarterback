package main

import (
	"os"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
	"github.com/bilberrry/quarterback/compression"
	"github.com/bilberrry/quarterback/encryption"
	"github.com/bilberrry/quarterback/source"
	"github.com/bilberrry/quarterback/storage"
)

type Target struct {
	Config common.TargetConfig
}

func (ctx Target) run() {
	logger.Info("########## Running target: " + ctx.Config.Name + " ##########")
	logger.Info("Working directory: ", ctx.Config.WorkPath+"\n")

	defer ctx.cleanup()

	err := source.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	archivePath, err := compression.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	archivePath, err = encryption.Run(archivePath, ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	err = storage.Run(ctx.Config, archivePath)
	if err != nil {
		logger.Error(err)
		return
	}

}

func (ctx Target) cleanup() {
	logger.Info("Cleaning up working directory...\n")

	err := os.RemoveAll(common.TempPath)

	if err != nil {
		logger.Error("Cleanup error: ", err)
	}

	logger.Info("########## Target finished: " + ctx.Config.Name + " ##########\n")
}
