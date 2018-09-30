package storage

import (
	"os"
	"path"
	"time"

	logger "github.com/Sirupsen/logrus"
	"github.com/secsy/goftp"
)

type FTP struct {
	Base
	path     string
	host     string
	port     string
	username string
	password string

	client *goftp.Client
}

func (ctx *FTP) open() (err error) {
	viper := ctx.viper

	viper.SetDefault("port", "21")
	viper.SetDefault("timeout", 300)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.path = viper.GetString("path")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")

	ftpConfig := goftp.Config{
		User:     ctx.username,
		Password: ctx.password,
		Timeout:  viper.GetDuration("timeout") * time.Second,
	}

	ctx.client, err = goftp.DialConfig(ftpConfig, ctx.host+":"+ctx.port)

	if err != nil {
		return err
	}

	return
}

func (ctx *FTP) close() {
	ctx.client.Close()
}

func (ctx *FTP) send(fileName string) (err error) {
	_, err = ctx.client.Stat(ctx.path)
	if os.IsNotExist(err) {
		if _, err := ctx.client.Mkdir(ctx.path); err != nil {
			return err
		}
	}

	file, err := os.Open(ctx.archivePath)
	if err != nil {
		return err
	}

	defer file.Close()

	logger.Info("=> Copying over FTP...")

	remotePath := path.Join(ctx.path, fileName)
	err = ctx.client.Store(remotePath, file)

	if err != nil {
		return err
	}

	logger.Info("=> Stored successfully")

	return nil
}

func (ctx *FTP) delete(fileName string) (err error) {
	remotePath := path.Join(ctx.path, fileName)

	err = ctx.client.Delete(remotePath)

	return
}
