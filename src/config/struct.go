package config

type (
	Config struct {
		Credentials MainCfg
		Schedules   ScheduleCfg
	}

	MainCfg struct {
		CredentialCfg Cred `gcfg:"config"`
	}

	Cred struct {
		DiscordURL   string
		GuildID      string
		WebhookToken string
	}

	ScheduleCfg struct {
		TimeInfo   TimeCfg `gcfg:"Time"`
		HardBosses map[string]*BossesCfg
		MidBosses  map[string]*BossesCfg
		EasyBosses map[string]*BossesCfg
	}

	TimeCfg struct {
		HardStart    string
		MidStart     string
		EasyStart    string
		HardInterval int
		MidInterval  int
		EasyInterval int
	}

	BossesCfg struct {
		Name     string `gcfg:"name"`
		Location string `gcfg:"location"`
	}
)
