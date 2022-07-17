package authserver

import (
	"auth/pkg/authmail"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type SMTPHostConfig struct {
	ENV_Username string `yaml:"ENV_Username"`
	ENV_Password string `yaml:"ENV_Password"`
	Host         string `yaml:"Host"`
	Port         int    `yaml:"Port"`
	Sender       string `yaml:"Sender"`
}

// ServerConfig defines configuration settings for the authentication server.
// ENV values are environment variable names
type AuthServerConfig struct {
	Domain          string         `yaml:"Domain"`
	Port            int            `yaml:"Port"`
	ENV_TokenSecret string         `yaml:"ENV_TokenSecret"`
	SMTPHost        SMTPHostConfig `yaml:"SMTPHost"`
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

func (s *AuthServerConfig) AsServer() *AuthServer {
	secret := os.Getenv(s.ENV_TokenSecret)
	smtpname := os.Getenv(s.SMTPHost.ENV_Username)
	smtppass := os.Getenv(s.SMTPHost.ENV_Password)
	if secret == "" || smtpname == "" || smtppass == "" {
		fmt.Printf("Must set ENV_TokenSecret, ENV_Username, ENV_Password\n")
		return nil
	}
	return &AuthServer{
		Address: s.Domain,
		Port:    s.Port,
		Secret:  []byte(secret),
		SESHost: authmail.Host{
			Username: smtpname,
			Password: smtppass,
			Host:     s.SMTPHost.Host,
			Port:     s.SMTPHost.Port,
			Sender:   s.SMTPHost.Sender,
		},
	}
}

func ServerFromConfig(fn string) (*AuthServer, error) {
	cfg := NewConfig()
	err := cfg.ReadConfig(fn)
	if err != nil {
		return nil, err
	}
	return cfg.AsServer(), nil
}
