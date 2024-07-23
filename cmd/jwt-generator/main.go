package main

import (
	"banners-management/internal/config"
	"banners-management/internal/lib/jwt"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	roleUser  = "user"
	roleAdmin = "admin"
)

func main() {
	flagSet := flag.NewFlagSet("jwt-generator-role", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	roleFlag := flagSet.String("role", "", "")

	_ = flagSet.Parse(os.Args[1:])

	if roleFlag == nil || *roleFlag == "" {
		fmt.Println("role is required")
		os.Exit(1)
	}
	role := *roleFlag
	if role != roleUser && role != roleAdmin {
		fmt.Println("role should be either 'user' or 'admin'")
		os.Exit(1)
	}
	cfg := config.MustLoad(os.Args[1:], os.LookupEnv)
	manager := jwt.NewManager(string(cfg.JwtSettings.SecretKey), time.Duration(cfg.JwtSettings.Expire))
	token, err := manager.GenerateToken(role)
	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
