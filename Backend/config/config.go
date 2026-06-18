package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var Env = loadEnv()

func Get(key string) string {
	return Env[key]
}

func init() {
	logfilname := fmt.Sprintf("./%s.log", "weblog")
	logFile, err := os.OpenFile(logfilname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
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
