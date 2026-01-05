package configs

type Config struct {
	MysqlConf MysqlConfig `json:"mysql_conf" yaml:"mysql_conf"`
	RedisConf RedisConfig `json:"redis_conf" yaml:"redis_conf"`
}

type MysqlConfig struct {
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Database string `yaml:"database" json:"database"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" json:"addr"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}
