package connection

import "time"

type PoolConfiguration struct {
	MaxConnections  int           `yaml:"max_connections"`
	Name            string        `yaml:"name"`
	Address         string        `yaml:"address"`
	DialTimeout     time.Duration `yaml:"dial_timeout"`
	RedialDelay     time.Duration `yaml:"redial_delay"`
	KeepAliveDelay  time.Duration `yaml:"keep_alive_delay"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ResponseTimeout time.Duration `yaml:"response_timeout"`
	TickDelay       time.Duration `yaml:"tick_delay"`
}
