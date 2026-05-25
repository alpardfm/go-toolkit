package configbuilder

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/alpardfm/go-toolkit/files"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cbroglie/mustache"
	"github.com/spf13/viper"
)

type Options struct {
	Namespace    string
	Env          string
	Key          string
	Secret       string
	Region       string
	TemplateFile string
	ConfigFile   string
}

type Interface interface {
	BuildConfig() error
}

type configBuilder struct {
	opt Options
	ssm *ssm.SSM
}

func Init(opt Options) (Interface, error) {
	var sess *session.Session

	if opt.Env == "" {
		return nil, errors.NewWithCode(codes.CodeConfigBuilderInit, "environment variable is not set")
	}

	if opt.Region == "" {
		opt.Region = "ap-southeast-1"
	}

	// if access key and secret not found in env, get credentials from local or metadata
	// this behaviour is intended to support all local, staging, and production environment
	// session credentials behaviour: https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
	if opt.Key != "" && opt.Secret != "" {
		sess = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(opt.Region),
			Credentials: credentials.NewStaticCredentials(opt.Key, opt.Secret, ""),
		}))
	} else {
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(opt.Region),
		}))
	}

	ssm := ssm.New(sess)
	return &configBuilder{
		opt: opt,
		ssm: ssm,
	}, nil
}

func (b *configBuilder) BuildConfig() error {
	if !files.IsExist(b.opt.TemplateFile) {
		return errors.NewWithCode(codes.CodeConfigBuilderBuild, "template file not found")
	}

	// Expire current credentials to refresh session in case previous session is still active
	// This is to avoid getting cached value on parameter store
	b.ssm.Config.Credentials.Expire()

	ssmparams := []*ssm.Parameter{}
	ssmres, err := b.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(fmt.Sprintf("/service/%s/%s/", b.opt.Env, b.opt.Namespace)),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return errors.NewWithCode(codes.CodeConfigBuilderBuild, "failed to get SSM parameters: %v", err)
	}
	ssmparams = append(ssmparams, ssmres.Parameters...)
	for ssmres.NextToken != nil {
		ssmres, err = b.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
			Path:           aws.String(fmt.Sprintf("/service/%s/%s/", b.opt.Env, b.opt.Namespace)),
			NextToken:      ssmres.NextToken,
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return errors.NewWithCode(codes.CodeConfigBuilderBuild, "failed to get SSM parameters (paginated): %v", err)
		}
		ssmparams = append(ssmparams, ssmres.Parameters...)
	}

	if len(ssmparams) < 1 {
		return errors.NewWithCode(codes.CodeConfigBuilderBuild, "no configuration found, please check your environment variable")
	}

	params := viper.New()
	for _, p := range ssmparams {
		val := *p.Value
		keyregxp := regexp.MustCompile(fmt.Sprintf("^\\/service\\/%s\\/%s\\/(.*)$", b.opt.Env, b.opt.Namespace))
		key := strings.ReplaceAll(keyregxp.ReplaceAllString(*p.Name, "$1"), "/", ".")
		params.Set(key, val)
	}

	body, err := os.ReadFile(b.opt.TemplateFile)
	if err != nil {
		return errors.NewWithCode(codes.CodeConfigBuilderBuild, "failed to read template file: %v", err)
	}

	conf, err := mustache.Render(string(body), true, params.AllSettings())
	if err != nil {
		return errors.NewWithCode(codes.CodeConfigBuilderBuild, "failed to render config template: %v", err)
	}

	f, err := os.Create(b.opt.ConfigFile)
	if err != nil {
		return errors.NewWithCode(codes.CodeConfigBuilderBuild, "failed to create config file: %v", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(conf))
	if err != nil {
		return errors.NewWithCode(codes.CodeConfigBuilderBuild, "failed to write config file: %v", err)
	}

	return nil
}
