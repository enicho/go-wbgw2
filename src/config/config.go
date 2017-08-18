package config

import gcfg "gopkg.in/gcfg.v1"

func InitConfig() Config {
	cfg := Config{}
	cfgCred := MainCfg{}
	cfgSchd := ScheduleCfg{}

	fname := "files/main.ini"
	err := gcfg.ReadFileInto(&cfgCred, fname)
	if err != nil {
		panic(err)
	}

	schedulesName := "files/schedule.ini"
	err = gcfg.ReadFileInto(&cfgSchd, schedulesName)
	if err != nil {
		panic(err)
	}

	cfg.Credentials = cfgCred
	cfg.Schedules = cfgSchd

	return cfg
}
