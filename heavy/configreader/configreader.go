package configreader

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/alpardfm/go-toolkit/files"
	"github.com/spf13/viper"
)

const (
	JSONType string = "json"
	YAMLType string = "yaml"
)

type Interface interface {
	ReadConfig(cfg any) error
	AllSettings() map[string]any
}

type Options struct {
	ConfigFile            string
	ExcelParserConfigFile string
}

type configReader struct {
	viper *viper.Viper
	opt   Options
}

func Init(opt Options) (Interface, error) {
	vp := viper.New()
	vp.SetConfigFile(opt.ConfigFile)
	if err := vp.ReadInConfig(); err != nil {
		return nil, errors.NewWithCode(codes.CodeConfigReaderInit, "failed to read config file: %v", err)
	}

	c := &configReader{
		viper: vp,
		opt:   opt,
	}

	return c, nil
}

func (c *configReader) mergeEnvConfig() {
	enver := os.Getenv("SERVICE_VERSION")
	sm := c.viper.GetStringMap("meta")
	if enver != "" {
		sm["version"] = enver
	} else {
		sm["version"] = "dev"
	}
	c.viper.Set("meta", sm)
}

func (c *configReader) resolveJSONRef() {
	refmap := make(map[string]any)
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

func (c *configReader) ReadConfig(cfg any) error {
	c.mergeEnvConfig()

	if files.GetExtension(filepath.Base(c.opt.ConfigFile)) == JSONType {
		c.resolveJSONRef()
	}
	if err := c.viper.Unmarshal(&cfg); err != nil {
		return errors.NewWithCode(codes.CodeConfigReaderRead, "failed to unmarshal config: %v", err)
	}
	return nil
}

func (c *configReader) AllSettings() map[string]any {
	return c.viper.AllSettings()
}
