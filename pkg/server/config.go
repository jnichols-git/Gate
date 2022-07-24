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

// Read environment variables contained in the AuthServerConfig.
// Any given AuthServerConfig object has a series of variables naming environment values set in the OS, to keep them private.
// This function takes those names and loads them into the config as their non-_ENV counterparts.
//
// Calling:
//   - cfg *AuthServerConfig: Config to read environment values into.
// Output:
//   - error: Returned if any ENV value couldn't be read.
func (cfg *AuthServerConfig) readEnvs() error {
	var ok bool = true
	cfg.SMTPHost.Username, ok = os.LookupEnv(cfg.SMTPHost.Username_ENV)
	if !ok {
		return errors.Errorf("Couldn't read %s", cfg.SMTPHost.Username_ENV)
	}
	cfg.SMTPHost.Password, ok = os.LookupEnv(cfg.SMTPHost.Password_ENV)
	if !ok {
		return errors.Errorf("Couldn't read %s", cfg.SMTPHost.Password_ENV)
	}
	cfg.JWT.TokenSecret, ok = os.LookupEnv(cfg.JWT.TokenSecret_ENV)
	if !ok {
		return errors.Errorf("Couldn't read %s", cfg.JWT.TokenSecret_ENV)
	}
	return nil
}

// Read a configuration file into the calling *AuthServerConfig.
//
// Calling:
//   - cfg *AuthServerConfig: Config to read file into. If no error results, cfg should be fully populated.
// Output:
//   - error: Any error that occurs when reading, including: file doesn't exist, invalid yaml format, envs not present
func (cfg *AuthServerConfig) ReadConfig(fn string) error {
	input, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(input, cfg)
	if err != nil {
		return err
	}
	err = cfg.readEnvs()
	if err != nil {
		return err
	}

	return nil
}
