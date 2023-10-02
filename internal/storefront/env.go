package storefront

import (
	"fmt"

	"github.com/spf13/viper"
)

// Use viper to read .env file
// Return the value of the key
func getEnv(key string) (string, error) {
	viper.SetConfigFile("../../.env")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		return "", fmt.Errorf("error reading env")
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		return "", fmt.Errorf("invalid env type")
	}

	return value, nil
}
