package config

var config = make(map[string]string, 0)

type Config struct {
	config map[string]string
}

func (c Config) loadConfig() {

}

func (c Config) Get(k string) (v string, ok bool) {
	v, ok = c.config[k]
	return
}

func New() *Config {
	return &Config{
		config: make(map[string]string, 0),
	}
}
