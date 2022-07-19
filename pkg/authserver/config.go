package authserver

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type SMTPHostConfig struct {
	username_ENV string `yaml:"ENV_Username"`
	Username     string `yaml:"-"`
	password_ENV string `yaml:"ENV_Password"`
	Password     string `yaml:"-"`
	Host         string `yaml:"Host"`
	Port         int    `yaml:"Port"`
	Sender       string `yaml:"Sender"`
	TestEmail    string `yaml:"TestEmail"`
}

type JWSConfig struct {
	tokenSecret_ENV string `yaml:"ENV_TokenSecret"`
	TokenSecret     string `yaml:"-"`
}

// ServerConfig defines configuration settings for the authentication server.
// ENV values are environment variable names
type AuthServerConfig struct {
	Domain   string         `yaml:"Domain"`
	Port     int            `yaml:"Port"`
	SMTPHost SMTPHostConfig `yaml:"SMTPHost"`
	JWS      JWSConfig      `yaml:"JWS"`
}

func NewConfig() *AuthServerConfig {
	return &AuthServerConfig{}
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
	return nil
}
