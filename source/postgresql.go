package source

import (
	"os"
	"path"
	"strings"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
)

type PostgreSQL struct {
	Base
	host        string
	port        string
	database    string
	username    string
	password    string
	dumpCommand string
}

func (ctx PostgreSQL) process() (err error) {
	viper := ctx.viper

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 5432)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	ctx.database = viper.GetString("source")

	if len(ctx.database) == 0 {
		logger.Error("PostgreSQL database name is required")
		return
	}

	err = ctx.dump()

	return
}

func (ctx *PostgreSQL) getArgs() []string {
	var args []string

	if len(ctx.host) > 0 {
		args = append(args, "--host="+ctx.host)
	}

	if len(ctx.port) > 0 {
		args = append(args, "--port="+ctx.port)
	}

	if len(ctx.username) > 0 {
		args = append(args, "--username="+ctx.username)
	}

	return args
}

func (ctx *PostgreSQL) dump() error {
	filePath := path.Join(ctx.WorkPath, ctx.database+".sql")

	logger.Info("=> Dumping PostgreSQL...")

	if len(ctx.password) > 0 {
		os.Setenv("PGPASSWORD", ctx.password)
	}

	_, err := common.Exec("pg_dump "+strings.Join(ctx.getArgs(), " ")+" "+ctx.database, "-f", filePath)

	if err != nil {
		logger.Error("Dump error: %s", err)
		return err
	}

	return nil
}
