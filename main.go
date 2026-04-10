package main

import (
	"bufio"
	"log"
	"os"
	"sdk/router"
	"strings"
)

func main() {
	loadDotEnv(".env")

	requiredEnvVars := []string{"TBS_CLIENT_ID", "TBS_SECRET", "TBS_USER_TOKEN"}
	for _, key := range requiredEnvVars {
		if os.Getenv(key) == "" {
			log.Printf("warning: %s is not set", key)
		}
	}

	log.Println("server running on http://localhost:8080")
	log.Fatal(router.NewRouter().Run(":8080"))
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("warning: unable to open %s: %v", path, err)
		}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
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
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			log.Printf("warning: unable to set %s from %s: %v", key, path, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("warning: unable to read %s: %v", path, err)
	}
}
