package configs

type Config struct {
	MysqlConf      MysqlConfig      `json:"mysql_conf" yaml:"mysql_conf"`
	LLMRequestConf LLMRequestConfig `json:"llm_request_conf" yaml:"llm_request_conf"`
	RedisConf      RedisConfig      `json:"redis_conf" yaml:"redis_conf"`
}
type LLMRequestConfig struct {
	BaseURL string `json:"base_url" yaml:"base_url"`
	ImgURL  string `json:"img_url" yaml:"img_url"`
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
