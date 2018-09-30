package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	logger "github.com/Sirupsen/logrus"
	"github.com/bilberrry/quarterback/common"
)

type PackageList []Package

type Package struct {
	File       string    `json:"file"`
	CreateDate time.Time `json:"create_date"`
}

var (
	cleanerPath = path.Join(common.HomeDir, ".quarterback")
)

type Cleaner struct {
	packages PackageList
	isLoaded bool
}

func (c *Cleaner) process(target string, fileName string, keep int, deletePackage func(fileName string) error) {
	cleanerFileName := path.Join(cleanerPath, target+".json")

	c.load(cleanerFileName)
	c.add(fileName)

	defer c.save(cleanerFileName)

	if keep == 0 {
		return
	}

	for {
		pkg := c.shift(keep)
		if pkg == nil {
			break
		}

		logger.Info("=> Pruning old backups...")

		err := deletePackage(pkg.File)
		if err != nil {
			logger.Warn("Can't remove: ", err)
		}
	}
}

func (c *Cleaner) add(fileName string) {
	c.packages = append(c.packages, Package{
		File:       fileName,
		CreateDate: time.Now(),
	})
}

func (c *Cleaner) shift(keep int) (first *Package) {
	total := len(c.packages)

	if total <= keep {
		return nil
	}

	first, c.packages = &c.packages[0], c.packages[1:]
	return
}

func (c *Cleaner) load(fileName string) {
	common.CreateDir(cleanerPath)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		ioutil.WriteFile(fileName, []byte("[]"), os.ModePerm)
	}

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error("Can't read backup history: ", err)
		return
	}

	err = json.Unmarshal(f, &c.packages)
	if err != nil {
		logger.Error("Can't parse backup history: ", err)
	}

	c.isLoaded = true
}

func (c *Cleaner) save(fileName string) {
	if !c.isLoaded {
		logger.Warn("Skipping backup history update, because file can't be loaded")
		return
	}

	data, err := json.Marshal(&c.packages)
	if err != nil {
		logger.Error("Can't encode backup history: ", err)
		return
	}

	err = ioutil.WriteFile(fileName, data, os.ModePerm)
	if err != nil {
		logger.Error("Can't save backup history: ", err)
		return
	}
}
