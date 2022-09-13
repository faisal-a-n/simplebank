package util

import "github.com/spf13/viper"

//This struct will hold the env data for the whole project
type Config struct {
	DB_DRIVER string `mapstructure:"DB_DRIVER"`
	DB_SOURCE string `mapstructure:"DB_SOURCE"`
	PORT      string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
