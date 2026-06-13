package config

import (
	"os"
	"strings"
)

var Env = loadEnv()

func Get(key string) string {
	return Env[key]
}

func loadEnv() map[string]string {
	env := make(map[string]string)

	for _, path := range []string{".env", "Backend/.env", "../Backend/.env"} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			value = strings.Trim(value, `"'`)
			if key != "" {
				env[key] = value
			}
		}

		return env
	}

	return env
}
