package storage

import (
	"avito-test-task/internal/lib/api/msg"
	"errors"
)

var (
	ErrBannerNotFound      = errors.New(msg.BannerNotFound)
	ErrBannerAlreadyExists = errors.New(msg.BannerAlreadyExists)
)
