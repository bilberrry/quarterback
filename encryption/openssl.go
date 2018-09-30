package encryption

import (
	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
)

type OpenSSL struct {
	Base
	salt     bool
	password string
}

func (ctx *OpenSSL) process() (encryptPath string, err error) {
	var options []string
	viper := ctx.viper

	viper.SetDefault("salt", true)

	ctx.salt = viper.GetBool("salt")
	ctx.password = viper.GetString("password")

	if len(ctx.password) == 0 {
		logger.Error("Password is required")
		return
	}

	encryptPath = ctx.archivePath + ".enc"

	options = append(options, "aes-256-cbc")

	if ctx.salt {
		options = append(options, "-salt")
	}

	options = append(options, `-k`, ctx.password)
	options = append(options, "-in", ctx.archivePath, "-out", encryptPath)

	_, err = common.Exec("openssl", options...)

	return
}
