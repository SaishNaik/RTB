package config

type MongoDBConfig struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
	Uri     string `mapstructure:"uri" json:"uri"`
	Timeout int    `mapstructure:"timeout" json:"timeout"`
}
