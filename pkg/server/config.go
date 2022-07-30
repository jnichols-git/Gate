package server

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config_Admin struct {
	Email    string `yaml:"-"`
	Username string `yaml:"-"`
	Password string `yaml:"-"`
}
type Config_SMTP struct {
	Username  string `yaml:"-"`
	Password  string `yaml:"-"`
	Host      string `yaml:"Host"`
	Port      int    `yaml:"Port"`
	Sender    string `yaml:"Sender"`
	TestEmail string `yaml:"TestEmail"`
}
type Config_GateKey struct {
	GatekeySecret  string `yaml:"-"`
	UserValidTime  int    `yaml:"UserValidTime"`
	AdminValidTime int    `yaml:"AdminValidTime"`
}

// ServerConfig defines configuration settings for the authentication server.
type AuthServerConfig struct {
	Domain     string         `yaml:"Domain"`
	Address    string         `yaml:"Address"`
	Port       int            `yaml:"Port"`
	Local      bool           `yaml:"Local"`
	SSLKeyFile string         `yaml:"-"`
	SSLCrtFile string         `yaml:"-"`
	Admin      Config_Admin   `yaml:"-"`
	SMTPHost   Config_SMTP    `yaml:"SMTP"`
	GateKey    Config_GateKey `yaml:"GateKey"`
}

func NewConfig() *AuthServerConfig {
	return &AuthServerConfig{}
}

// Read a secret value.
// This reads secret values from Docker Swarm. See README for how to set these up.
//
// Input:
//   - name string: Secret name.
//   - fileOnly bool: If true, this function will return only the filepath and not its contents.
// Output:
//   - string: Secret value
//   - error: If the secret file can't be read, returns this.
func getSecret(name string, fileOnly bool) (string, error) {
	path := fmt.Sprintf("/run/secrets/%s", name)
	if fileOnly {
		return path, nil
	}
	val, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Couldn't read secret %s; %v", name, err)
	}
	return strings.TrimSpace(string(val)), nil
}

// Read secrets needed for AuthServerConfig from docker swarm secrets.
//
// Calling:
//   - cfg *AuthServerConfig: Config to read environment values into.
// Output:
//   - error: Returned if any ENV value couldn't be read.
func (cfg *AuthServerConfig) readSecrets() error {
	var err error
	// Admin
	cfg.Admin.Email, err = getSecret("gate-admin-email", false)
	if err != nil {
		return err
	}
	cfg.Admin.Username, err = getSecret("gate-admin-username", false)
	if err != nil {
		return err
	}
	// This one can actually just fall through if it doesn't succeed. server.go handles this failure case.
	cfg.Admin.Password, _ = getSecret("gate-admin-password", false)
	// SMTP
	cfg.SMTPHost.Username, err = getSecret("gate-smtp-username", false)
	if err != nil {
		return err
	}
	cfg.SMTPHost.Password, err = getSecret("gate-smtp-password", false)
	if err != nil {
		return err
	}
	// Generate a random 32-byte secret
	secret := make([]byte, 32)
	rand.Read(secret[:])
	cfg.GateKey.GatekeySecret = base64.RawURLEncoding.EncodeToString(secret)
	// Get SSL certs
	cfg.SSLKeyFile, err = getSecret("gate-ssl-key", true)
	cfg.SSLCrtFile, err = getSecret("gate-ssl-crt", true)

	return nil
}

// Imitate the behavior of readSecrets using environment variables.
// Used for local running of Gate outside of Docker.
func (cfg *AuthServerConfig) readEnvironment() error {
	var ok bool
	// Admin
	cfg.Admin.Email, ok = os.LookupEnv("GATE_ADMIN_EMAIL")
	if !ok {
		return fmt.Errorf("Couldn't find environment variable %s.", "GATE_ADMIN_EMAIL")
	}
	cfg.Admin.Username, ok = os.LookupEnv("GATE_ADMIN_USERNAME")
	if !ok {
		return fmt.Errorf("Couldn't find environment variable %s.", "GATE_ADMIN_USERNAME")
	}
	cfg.Admin.Password, ok = os.LookupEnv("GATE_ADMIN_PASSWORD")
	if !ok {
		return fmt.Errorf("Couldn't find environment variable %s.", "GATE_ADMIN_PASSWORD")
	}
	// SMTP
	cfg.SMTPHost.Username, ok = os.LookupEnv("GATE_SMTP_USERNAME")
	if !ok {
		return fmt.Errorf("Couldn't find environment variable %s.", "GATE_SMTP_USERNAME")
	}
	cfg.SMTPHost.Password, ok = os.LookupEnv("GATE_SMTP_PASSWORD")
	if !ok {
		return fmt.Errorf("Couldn't find environment variable %s.", "GATE_SMTP_PASSWORD")
	}
	// Gatekey
	secret := make([]byte, 32)
	rand.Read(secret[:])
	cfg.GateKey.GatekeySecret = base64.RawURLEncoding.EncodeToString(secret)
	// SSL certs
	cfg.SSLKeyFile, ok = os.LookupEnv("GATE_SSL_KEY")
	if !ok {
		return fmt.Errorf("Couldn't find environment variable %s.", "GATE_SSL_KEY")
	}
	cfg.SSLCrtFile, ok = os.LookupEnv("GATE_SSL_CRT")
	if !ok {
		return fmt.Errorf("Couldn't find environment variable %s.", "GATE_SSL_CRT")
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
	if cfg.Local {
		err = cfg.readEnvironment()
	} else {
		err = cfg.readSecrets()
	}
	if err != nil {
		return err
	}
	if cfg.Local {
		cfg.Address = "localhost"
	}
	cfg.Address = fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	return nil
}
