package config

import (
	"os"
)

type Config struct {
	Port string
	Env  string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development" // default environment
	}

	return &Config{
		Port: port,
		Env:  env,
	}, nil
}

var DevMode bool
var DevUser = struct {
	Username string
	UserId   int
	APIKey   string
	FolderId int64
}{
	Username: "devuser",
	UserId:   1,
	APIKey:   "bbedfac8e64caff3f2aa7d4226d79e8e79f69c2ab4d5fa6a695e67d265cd23b2",
	FolderId: 1,
}
