package config

//todo set up fields in main config
type MainConfig struct {
	Application       `json:"app" mapstructure:"application"`
	GinMode           string        `json:"gin_mode" mapstructure:"gin_mode"`
	MongoDB           MongoDBConfig `json:"mongo_db" mapstructure:"mongo_db"`
	IP2LocationDBPATH string        `json:"ip2location_data_path" mapstructure:"ip2location_data_path"`
}

type Application struct {
	//Name string `json:"name" mapstructure:"name"`
	Port string `json:"port" mapstructure:"port"`
	//Env  string `json:"env" mapstructure:"env"`
}
