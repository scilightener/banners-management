package main

import (
	"avito-test-task/migrator"
	"context"
	"os"
)

func main() {
	migrator.Migrate(context.Background(), os.Stdout, os.Args[1:], os.LookupEnv, "postgres")
}
