package env

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var once sync.Once

// Parse wraps caarlos0/env package with .env loading
func Parse(v interface{}) {
	once.Do(func() {
		envFile := filepath.Join(getRootDir(), ".env")
		if err := godotenv.Load(envFile); err != nil {
			log.Println(errors.Wrap(err, "can't load variables from .env file"))
		}
	})
	if err := env.Parse(v); err != nil {
		panic(errors.Wrap(err, "configuration setup failed"))
	}
}

func getRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("can't get working directory")
	}

	dirs := strings.Split(cwd, string(os.PathSeparator))
	dirs[0] = string(filepath.Separator) + dirs[0]
	for i, dir := range dirs {
		if dir == "internal" || dir == "cmd" {
			dirs = dirs[:i]
			break
		}
	}
	return filepath.Join(dirs...)
}
