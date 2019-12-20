package config

type (
	Config interface {
		LoadConfig()
		Get(string) (string, bool)
	}
	config struct {
		config map[string]string
	}
)

func (c *config) LoadConfig() {

}

func (c *config) Get(k string) (v string, ok bool) {
	v, ok = c.config[k]
	return
}

func New() Config {
	return &config{
		config: make(map[string]string),
	}
}
