package compression

import (
	"os"
	"path"
	"time"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
	"github.com/spf13/viper"
)

type Base struct {
	name   string
	target common.TargetConfig
	viper  *viper.Viper
}

type Context interface {
	process() (archivePath string, err error)
}

func (ctx *Base) getFilePath(ext string) string {
	return path.Join(os.TempDir(), os.Args[0], time.Now().Format("2006-01-02-15_04_05")+ext)
}

func initBase(target common.TargetConfig) (base Base) {
	base = Base{
		name:   target.Name,
		target: target,
		viper:  target.Compression.Viper,
	}
	return
}

func Run(target common.TargetConfig) (archivePath string, err error) {
	base := initBase(target)

	var ctx Context

	switch target.Compression.Type {
	case "tgz":
		ctx = &Tgz{Base: base}
	default:
		ctx = &Tgz{Base: base}
	}

	logger.Info("---------- Compression started")
	logger.Info("=> Type: " + target.Compression.Type)

	_ = os.Chdir(path.Join(target.WorkPath, "../"))

	archivePath, err = ctx.process()

	if err != nil {
		return
	}

	logger.Info("=> Archive: ", archivePath)
	logger.Info("------------ Compression finished\n")

	return
}
