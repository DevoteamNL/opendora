package config

import "os"

func LookupEnvDefault(envKey string, envDefaultValue string) string {
	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}
	return envDefaultValue
}
