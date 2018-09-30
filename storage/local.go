package storage

import (
	"path"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
)

type Local struct {
	Base
	destPath string
}

func (ctx *Local) open() (err error) {
	ctx.destPath = ctx.viper.GetString("path")
	common.CreateDir(ctx.destPath)
	return
}

func (ctx *Local) close() {}

func (ctx *Local) send(fileName string) (err error) {
	logger.Info("=> Copying locally...")

	_, err = common.Exec("cp", ctx.archivePath, ctx.destPath)

	if err != nil {
		return err
	}

	logger.Info("=> Stored successfully: ", ctx.destPath)

	return nil
}

func (ctx *Local) delete(fileName string) (err error) {
	_, err = common.Exec("rm", path.Join(ctx.destPath, fileName))
	return
}
