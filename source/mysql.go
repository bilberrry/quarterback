package source

import (
	"path"
	"strings"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
)

type MySQL struct {
	Base
	host     string
	port     string
	database string
	username string
	password string
	options  []string
}

func (ctx *MySQL) process() (err error) {
	viper := ctx.viper

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("username", "root")
	viper.SetDefault("port", 3306)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	ctx.database = viper.GetString("database")

	additionalOptions := viper.GetString("options")

	if len(additionalOptions) > 0 {
		ctx.options = strings.Split(additionalOptions, " ")
	}

	if len(ctx.database) == 0 {
		logger.Error("MySQL database name is required")
		return
	}

	err = ctx.dump()

	return
}

func (ctx *MySQL) getArgs() []string {
	var args []string

	if len(ctx.host) > 0 {
		args = append(args, "--host", ctx.host)
	}

	if len(ctx.port) > 0 {
		args = append(args, "--port", ctx.port)
	}

	if len(ctx.username) > 0 {
		args = append(args, "-u", ctx.username)
	}

	if len(ctx.password) > 0 {
		args = append(args, `-p`+ctx.password)
	}

	if len(ctx.options) > 0 {
		args = append(args, ctx.options...)
	}

	args = append(args, ctx.database)

	filePath := path.Join(ctx.WorkPath, ctx.database+".sql")

	args = append(args, "--result-file="+filePath)

	return args
}

func (ctx *MySQL) dump() error {
	logger.Info("=> Dumping MySQL...")

	_, err := common.Exec("mysqldump", ctx.getArgs()...)

	if err != nil {
		logger.Error("Dump error: %s", err)
		return err
	}

	return nil
}
