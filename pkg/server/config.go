package server

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type SMTPHostConfig struct {
	Username_ENV string `yaml:"ENV_Username"`
	Username     string `yaml:"-"`
	Password_ENV string `yaml:"ENV_Password"`
	Password     string `yaml:"-"`
	Host         string `yaml:"Host"`
	Port         int    `yaml:"Port"`
	Sender       string `yaml:"Sender"`
	TestEmail    string `yaml:"TestEmail"`
}

type DBConfig struct {
	Path string `yaml:"Path"`
}

type JWTConfig struct {
	TokenSecret_ENV string `yaml:"ENV_TokenSecret"`
	TokenSecret     string `yaml:"-"`
	UserValidTime   int    `yaml:"UserValidTime"`
	AdminValidTime  int    `yaml:"AdminValidTime"`
}

// ServerConfig defines configuration settings for the authentication server.
// ENV values are environment variable names
type AuthServerConfig struct {
	Domain   string         `yaml:"Domain"`
	Port     int            `yaml:"Port"`
	SMTPHost SMTPHostConfig `yaml:"SMTPHost"`
	DB       DBConfig       `yaml:"Database"`
	JWT      JWTConfig      `yaml:"JWT"`
}

func NewConfig() *AuthServerConfig {
	return &AuthServerConfig{}
}

func (s *AuthServerConfig) ReadEnvs() error {
	var ok bool = true
	s.SMTPHost.Username, ok = os.LookupEnv(s.SMTPHost.Username_ENV)
	if !ok {
		return errors.Errorf("Couldn't read %s", s.SMTPHost.Username_ENV)
	}
	s.SMTPHost.Password, ok = os.LookupEnv(s.SMTPHost.Password_ENV)
	if !ok {
		return errors.Errorf("Couldn't read %s", s.SMTPHost.Password_ENV)
	}
	s.JWT.TokenSecret, ok = os.LookupEnv(s.JWT.TokenSecret_ENV)
	if !ok {
		return errors.Errorf("Couldn't read %s", s.JWT.TokenSecret_ENV)
	}
	return nil
}

func (s *AuthServerConfig) ReadConfig(fn string) error {
	cfg, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(cfg, s)
	if err != nil {
		return err
	}
	err = s.ReadEnvs()
	if err != nil {
		return err
	}

	return nil
}
