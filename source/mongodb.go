package source

import (
	"strings"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
)

type MongoDB struct {
	Base
	host     string
	port     string
	database string
	username string
	password string
	authdb   string
	oplog    bool
}

func (ctx *MongoDB) process() (err error) {
	viper := ctx.viper

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("username", "root")
	viper.SetDefault("port", 27017)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	ctx.database = viper.GetString("database")
	ctx.authdb = viper.GetString("authdb")

	err = ctx.dump()

	if err != nil {
		return err
	}

	return nil
}

func (ctx *MongoDB) getArgs() []string {
	var args []string

	if len(ctx.database) > 0 {
		args = append(args, "--db="+ctx.database)
	}

	if len(ctx.username) > 0 {
		args = append(args, "--username="+ctx.username)
	}

	if len(ctx.password) > 0 {
		args = append(args, `--password=`+ctx.password)
	}

	if len(ctx.authdb) > 0 {
		args = append(args, "--authenticationDatabase="+ctx.authdb)
	}

	if len(ctx.host) > 0 {
		args = append(args, "--host="+ctx.host)
	}
	if len(ctx.port) > 0 {
		args = append(args, "--port="+ctx.port)
	}

	return args
}

func (ctx *MongoDB) dump() error {
	logger.Info("=> Dumping MongoDB...")

	_, err := common.Exec("mongodump", strings.Join(ctx.getArgs(), " "), "--out="+ctx.WorkPath)

	if err != nil {
		logger.Error("Dump error: %s", err)
		return err
	}

	return nil
}
