package config

import (
	"fmt"
	"os"
)

type Config struct {
	SatoriAppKey   string
	SatoriEndpoint string
	SatoriChannel  string
	SatoriRole     string
	SatoriSecret   string
	YoutubeAuth    string
}

func LoadConfig() (*Config, error) {
	appKey, err := getRequiredEnvVariable("SATORI_APP_KEY")
	endpoint, err := getRequiredEnvVariable("SATORI_ENDPOINT")
	channel, err := getRequiredEnvVariable("SATORI_CHANNEL")
	role, err := getRequiredEnvVariable("SATORI_ROLE")
	secret, err := getRequiredEnvVariable("SATORI_SECRET")
	authorization, err := getRequiredEnvVariable("YOUTUBE_AUTHORIZATION")

	return &Config{
		SatoriAppKey:   appKey,
		SatoriEndpoint: endpoint,
		SatoriChannel:  channel,
		SatoriRole:     role,
		SatoriSecret:   secret,
		YoutubeAuth:    authorization,
	}, err
}

func getRequiredEnvVariable(key string) (string, error) {
	envVar := os.Getenv(key)

	if envVar == "" {
		return "", fmt.Errorf("Missing or blank '%s' required environment variable", key)
	}

	return envVar, nil
}
