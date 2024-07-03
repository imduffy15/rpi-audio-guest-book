package config

import (
	"github.com/spf13/viper"
)

type InputType string

const (
	Keyboard InputType = "keyboard"
	GPIO     InputType = "gpio"
)

type Config struct {
	InputListener  InputType `mapstructure:"input_listener" required:"true"`
	Volume         int       `mapstructure:"volume" required:"true"`
	GPIOBoard      string    `mapstructure:"gpio_board" required:"true"`
	GPIOPin        int       `mapstructure:"gpio_pin" required:"true"`
	RecordingsPath string    `mapstructure:"recordings_path" required:"true"`
}

func LoadViperConfig() *Config {
	viper.AutomaticEnv()

	viper.SetDefault("input_listener", Keyboard)
	viper.SetDefault("volume", 100)
	viper.SetDefault("gpio_board", "gpiochip0")
	viper.SetDefault("gpio_pin", 18)
	viper.SetDefault("recordings_path", "recordings/")

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil
	}

	return config
}
