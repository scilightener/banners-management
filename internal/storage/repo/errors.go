package repo

import (
	"errors"

	"banners-management/internal/lib/api/msg"
)

var (
	ErrBannerNotFound      = errors.New(msg.BannerNotFound)
	ErrBannerAlreadyExists = errors.New(msg.BannerAlreadyExists)
	ErrBannerNotUnique     = errors.New(msg.BannerNotUnique)
)
