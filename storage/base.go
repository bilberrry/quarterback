package storage

import (
	"fmt"
	"path/filepath"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
	"github.com/spf13/viper"
)

type Base struct {
	target      common.TargetConfig
	archivePath string
	viper       *viper.Viper
	keep        int
}

type Context interface {
	open() error
	close()
	send(fileName string) error
	delete(fileName string) error
}

func initBase(target common.TargetConfig, archivePath string) (base Base) {
	base = Base{
		target:      target,
		archivePath: archivePath,
		viper:       target.Storage.Viper,
	}

	if base.viper != nil {
		base.keep = base.viper.GetInt("keep")
	}

	return
}

func Run(target common.TargetConfig, archivePath string) (err error) {
	logger.Info("---------- Storing process started\n")

	newFileName := filepath.Base(archivePath)
	base := initBase(target, archivePath)

	var ctx Context

	switch target.Storage.Type {
	case "local":
		ctx = &Local{Base: base}
	case "ftp":
		ctx = &FTP{Base: base}
	case "scp":
		ctx = &SCP{Base: base}
	case "s3":
		ctx = &S3{Base: base}
	default:
		logger.Warn(fmt.Errorf("storage type `%s` is not implemented", target.Storage.Type))
		return
	}

	logger.Info("=> Storage: ", ", type: ", target.Storage.Type)

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
	cleaner.process(target.Name, newFileName, base.keep, ctx.delete)

	logger.Info("---------- Storing process finished\n")

	return nil
}
