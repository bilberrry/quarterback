package source

import (
	"fmt"
	"path"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
	"github.com/spf13/viper"
)

type Base struct {
	target       common.TargetConfig
	sourceConfig common.SubConfig
	viper        *viper.Viper
	name         string
	WorkPath     string
}

type Context interface {
	process() error
}

func Run(target common.TargetConfig) error {
	if len(target.Sources) == 0 {
		return nil
	}

	logger.Info("---------- Sources processing started")

	for _, sourceConfig := range target.Sources {
		err := processSource(target, sourceConfig)
		if err != nil {
			return err
		}
	}

	logger.Info("---------- Sources processing finished\n")

	return nil
}

func initBase(target common.TargetConfig, sourceConfig common.SubConfig) (base Base) {
	base = Base{
		target:       target,
		sourceConfig: sourceConfig,
		viper:        sourceConfig.Viper,
		name:         sourceConfig.Name,
	}

	base.WorkPath = path.Join(target.WorkPath, sourceConfig.Type, base.name)
	common.CreateDir(base.WorkPath)

	return
}

func processSource(target common.TargetConfig, sourceConfig common.SubConfig) (err error) {
	base := initBase(target, sourceConfig)

	var ctx Context

	switch sourceConfig.Type {
	case "mysql":
		ctx = &MySQL{Base: base}
	case "postgresql":
		ctx = &PostgreSQL{Base: base}
	case "mongodb":
		ctx = &MongoDB{Base: base}
	case "fs":
		ctx = &FS{Base: base}
	default:
		logger.Warn(fmt.Errorf("source type `%s` is not implemented", sourceConfig.Type))
		return
	}

	logger.Info("=> Source: ", base.name, ", type: ", sourceConfig.Type)

	err = ctx.process()

	if err != nil {
		return err
	}

	return
}
