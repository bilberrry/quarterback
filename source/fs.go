package source

import (
	"path"
	"path/filepath"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
)

type FS struct {
	Base
	includes []string
	excludes []string
}

func (ctx *FS) process() (err error) {
	viper := ctx.viper

	common.CreateDir(ctx.WorkPath)

	include := viper.GetStringSlice("include")
	include = ctx.cleanPaths(include)

	exclude := viper.GetStringSlice("exclude")
	exclude = ctx.cleanPaths(exclude)

	if len(include) == 0 {
		logger.Error("Include path is required")
		return
	}

	logger.Info("=> Copying files...")

	args := ctx.getArgs(ctx.WorkPath, exclude, include)

	common.Exec("tar", args...)

	return nil
}

func (ctx *FS) getArgs(WorkPath string, excludes, includes []string) (args []string) {
	tarPath := path.Join(WorkPath, "archive.tar")

	args = append(args, "-cPf", tarPath)

	for _, exclude := range excludes {
		args = append(args, "--exclude='"+filepath.Clean(exclude)+"'")
	}

	args = append(args, includes...)

	return args
}

func (ctx *FS) cleanPaths(paths []string) (results []string) {
	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}
	return
}
