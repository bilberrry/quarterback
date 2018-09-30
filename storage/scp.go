package storage

import (
	"os"
	"path"
	"time"

	logger "github.com/Sirupsen/logrus"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

type SCP struct {
	Base
	path       string
	host       string
	port       string
	privateKey string
	username   string
	password   string
	client     scp.Client
}

func (ctx *SCP) open() (err error) {
	viper := ctx.viper

	var clientConfig ssh.ClientConfig

	viper.SetDefault("port", "22")
	viper.SetDefault("timeout", 300)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.path = viper.GetString("path")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	ctx.privateKey = viper.GetString("private_key")

	clientConfig, err = auth.PrivateKey(
		ctx.username,
		ctx.privateKey,
		ssh.InsecureIgnoreHostKey(),
	)

	if err != nil {
		logger.Warn(err)

		clientConfig = ssh.ClientConfig{
			User:            ctx.username,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	clientConfig.Timeout = ctx.viper.GetDuration("timeout") * time.Second

	if len(ctx.password) > 0 {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(ctx.password))
	}

	ctx.client = scp.NewClient(ctx.host+":"+ctx.port, &clientConfig)

	err = ctx.client.Connect()

	if err != nil {
		return err
	}

	defer ctx.client.Session.Close()

	ctx.client.Session.Run("mkdir -p " + ctx.path)

	return
}

func (ctx *SCP) close() {}

func (ctx *SCP) send(fileName string) (err error) {
	err = ctx.client.Connect()

	if err != nil {
		return err
	}

	defer ctx.client.Session.Close()

	file, err := os.Open(ctx.archivePath)

	if err != nil {
		return err
	}

	defer file.Close()

	remotePath := path.Join(ctx.path, fileName)

	logger.Info("=> Copying over SSH...")

	ctx.client.CopyFromFile(*file, remotePath, "0644")

	logger.Info("=> Stored successfully")

	return nil
}

func (ctx *SCP) delete(fileName string) (err error) {
	err = ctx.client.Connect()

	if err != nil {
		return
	}

	defer ctx.client.Session.Close()

	remotePath := path.Join(ctx.path, fileName)

	err = ctx.client.Session.Run("rm " + remotePath)

	return
}
