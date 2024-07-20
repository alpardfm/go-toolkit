package configreader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alpardfm/go-toolkit/files"
	"github.com/spf13/viper"
)

const (
	JSONType string = "json"
	YAMLType string = "yaml"
)

type Interface interface {
	ReadConfig(cfg interface{})
	AllSettings() map[string]interface{}
}

type Options struct {
	ConfigFile            string
	ExcelParserConfigFile string
}

type configReader struct {
	viper *viper.Viper
	opt   Options
}

func Init(opt Options) Interface {

	vp := viper.New()
	vp.SetConfigFile(opt.ConfigFile)
	if err := vp.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error found during reading file. err: %w", err))
	}

	c := &configReader{
		viper: vp,
		opt:   opt,
	}

	return c
}

func (c *configReader) mergeEnvConfig() {
	enver := os.Getenv("AQUAHERO_SERVICE_VERSION")
	sm := c.viper.GetStringMap("meta")
	if enver != "" {
		sm["version"] = enver
	} else {
		sm["version"] = "dev"
	}
	c.viper.Set("meta", sm)
}

func (c *configReader) resolveJSONRef() {
	refmap := make(map[string]interface{})
	refregxp := regexp.MustCompile(`^\\$ref:#\\/(.*)$`)
	for _, k := range c.viper.AllKeys() {
		refpath := c.viper.GetString(k)
		if refregxp.MatchString(refpath) {
			v, ok := refmap[refpath]
			if !ok {
				refkey := refregxp.ReplaceAllString(refpath, "$1")
				refkey = strings.ToLower(strings.ReplaceAll(refkey, "/", "."))
				refmap[refpath] = c.viper.Get(refkey)
				c.viper.Set(k, refmap[refpath])
			} else {
				c.viper.Set(k, v)
			}
		}
	}
}

func (c *configReader) ReadConfig(cfg interface{}) {
	c.mergeEnvConfig()

	if files.GetExtension(filepath.Base(c.opt.ConfigFile)) == JSONType {
		c.resolveJSONRef()
	}
	if err := c.viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("fatal error found during unmarshaling config. err: %w", err))
	}
}

func (c *configReader) AllSettings() map[string]interface{} {
	return c.viper.AllSettings()
}
