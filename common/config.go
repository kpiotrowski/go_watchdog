package common

import "gopkg.in/gcfg.v1"


type Config struct {
	Mail MailConf
}

type MailConf struct {
	MailFromAddress string
	MailFromPassword string
	MailServerAddress string
	MailServerPort int
	MailTo string
}

func LoadConfig(file string) (*Config, error) {
	config := new(Config)
	err := gcfg.ReadFileInto(config, file)
	if err != nil {
		return nil, err
	}
	return config, nil
}
