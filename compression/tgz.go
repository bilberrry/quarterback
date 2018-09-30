package compression

import (
	"github.com/bilberrry/quarterback/common"
)

type Tgz struct {
	Base
}

func (ctx *Tgz) process() (archivePath string, err error) {
	var options []string

	filePath := ctx.getFilePath(".tar.gz")

	options = append(options, "-zcf")
	options = append(options, filePath)
	options = append(options, ctx.name)

	_, err = common.Exec("tar", options...)

	if err == nil {
		archivePath = filePath
		return
	}
	return
}
