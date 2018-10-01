package common

import (
	"fmt"
	"os"
	"path"
	"time"

	logger "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Targets  []TargetConfig
	HomeDir  string
	TempPath string
)

type TargetConfig struct {
	Name        string
	WorkPath    string
	Compression SubConfig
	Encryption  SubConfig
	Archive     *viper.Viper
	Sources     []SubConfig
	Storages    []SubConfig
	Viper       *viper.Viper
}

type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

func init() {
	viper.SetConfigType("yaml")

	HomeDir = os.Getenv("HOME")
	TempPath = path.Join(os.TempDir(), os.Args[0])

	viper.SetConfigName("quarterback")

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.quarterback")
	viper.AddConfigPath("/etc/quarterback/")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Can't find quarterback config: ", err)
		return
	}

	Targets = []TargetConfig{}
	for key := range viper.GetStringMap("targets") {
		Targets = append(Targets, loadTarget(key))
	}

	return
}

func loadTarget(key string) (target TargetConfig) {
	target.Name = key
	target.WorkPath = path.Join(TempPath, fmt.Sprintf("%d", time.Now().UnixNano()), key)
	target.Viper = viper.Sub("targets." + key)

	target.Compression = SubConfig{
		Type:  target.Viper.GetString("compression.type"),
		Viper: target.Viper.Sub("compression"),
	}

	target.Encryption = SubConfig{
		Type:  target.Viper.GetString("encryption.type"),
		Viper: target.Viper.Sub("encryption"),
	}

	loadSourcesConfig(&target)
	loadStoragesConfig(&target)

	return
}

func loadSourcesConfig(target *TargetConfig) {
	subViper := target.Viper.Sub("sources")

	for key := range target.Viper.GetStringMap("sources") {

		sourceViper := subViper.Sub(key)
		target.Sources = append(target.Sources, SubConfig {
			Name:  key,
			Type:  sourceViper.GetString("type"),
			Viper: sourceViper,
		})
	}
}


func loadStoragesConfig(target *TargetConfig) {
	subViper := target.Viper.Sub("storages")

	for key := range target.Viper.GetStringMap("storages") {

		storageViper := subViper.Sub(key)
		target.Storages = append(target.Storages, SubConfig {
			Name:  key,
			Type:  storageViper.GetString("type"),
			Viper: storageViper,
		})
	}
}

func (target *TargetConfig) GetSourceByName(name string) (subConfig *SubConfig) {
	for _, m := range target.Sources {
		if m.Name == name {

			subConfig = &m
			return
		}
	}
	return
}
