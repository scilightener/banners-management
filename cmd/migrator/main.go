package main

import (
	"banners-management/migrator"
	"context"
	"os"
)

func main() {
	migrator.Migrate(context.Background(), os.Stdout, os.Args[1:], os.LookupEnv, "postgres")
}
