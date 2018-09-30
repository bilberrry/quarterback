package storage

import (
	"fmt"
	"path/filepath"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
	"github.com/spf13/viper"
)

type Base struct {
	target        common.TargetConfig
	storageConfig common.SubConfig
	archivePath   string
	viper         *viper.Viper
	name          string
	keep          int
}

type Context interface {
	open() error
	close()
	send(fileName string) error
	delete(fileName string) error
}

func Run(target common.TargetConfig, archivePath string) (err error) {
	if len(target.Storages) == 0 {
		return nil
	}

	logger.Info("---------- Storing process started")

	for _, storageConfig := range target.Storages {
		err := processStorage(target, storageConfig, archivePath)
		if err != nil {
			return err
		}
	}

	logger.Info("---------- Storing process finished\n")

	return nil
}

func initBase(target common.TargetConfig, storageConfig common.SubConfig, archivePath string) (base Base) {
	base = Base{
		target:        target,
		storageConfig: storageConfig,
		viper:         storageConfig.Viper,
		name:          storageConfig.Name,
		archivePath:   archivePath,
	}

	if base.viper != nil {
		base.keep = base.viper.GetInt("keep")
	}

	return
}

func processStorage(target common.TargetConfig, storageConfig common.SubConfig, archivePath string) (err error) {
	base := initBase(target, storageConfig, archivePath)

	newFileName := filepath.Base(archivePath)

	var ctx Context

	switch storageConfig.Type {
	case "local":
		ctx = &Local{Base: base}
	case "ftp":
		ctx = &FTP{Base: base}
	case "scp":
		ctx = &SCP{Base: base}
	case "s3":
		ctx = &S3{Base: base}
	default:
		logger.Warn(fmt.Errorf("storage type `%s` is not implemented", storageConfig.Type))
		return
	}

	logger.Info("=> Storage: ", storageConfig.Name, ", type: ", storageConfig.Type)

	err = ctx.open()
	if err != nil {
		return err
	}

	defer ctx.close()

	err = ctx.send(newFileName)
	if err != nil {
		return err
	}

	cleaner := Cleaner{}
	cleaner.process(target.Name, storageConfig.Name, newFileName, base.keep, ctx.delete)

	return
}
