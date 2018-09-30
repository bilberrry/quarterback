package encryption

import (
	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
	"github.com/spf13/viper"
)

type Base struct {
	target      common.TargetConfig
	viper       *viper.Viper
	archivePath string
}

type Context interface {
	process() (encryptPath string, err error)
}

func initBase(archivePath string, target common.TargetConfig) (base Base) {
	base = Base{
		archivePath: archivePath,
		target:      target,
		viper:       target.Encryption.Viper,
	}
	return
}

func Run(archivePath string, target common.TargetConfig) (encryptPath string, err error) {
	base := initBase(archivePath, target)

	var ctx Context

	switch target.Encryption.Type {
	case "openssl":
		ctx = &OpenSSL{Base: base}
	default:
		encryptPath = archivePath
		return
	}

	logger.Info("------------ Encryption started")
	logger.Info("=> Type: " + target.Encryption.Type)

	encryptPath, err = ctx.process()

	if err != nil {
		return
	}
	logger.Info("=> Archive: ", encryptPath)
	logger.Info("------------ Encryption finished\n")

	return
}
