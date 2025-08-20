package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DB_URL          string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	pathToWrite := homedir + "/" + CONFIG_DB_PATH
	// let's write the json to the file
	return write(pathToWrite, *c)
}

func write(pathToWrite string, c Config) error {
	// let's export the json into the file system
	f, err := os.Create(pathToWrite)
	if err != nil {
		return err
	}

	defer f.Close()

	json_str, err := json.Marshal(c)
	f.Write(json_str)
	return nil
}

func Read() (Config, error) {
	homepath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	cfg, err := read(homepath + "/" + CONFIG_DB_PATH)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func read(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
