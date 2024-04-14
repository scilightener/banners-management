package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/go-playground/validator/v10"

	"avito-test-task/internal/lib/api/msg"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/models/dto/banner"
	"avito-test-task/internal/models/entity"
	"avito-test-task/internal/storage"
	"avito-test-task/internal/storage/repo"
)

var (
	ErrBannerNotSaved      = errors.New(msg.ErrSavingBanner)
	ErrBannerAlreadyExists = errors.New(msg.BannerAlreadyExists)
	ErrBannerNotFound      = errors.New(msg.BannerNotFound)
	ErrBannerNotActive     = errors.New(msg.BannerNotActive)
	ErrUnknown             = errors.New(msg.ErrUnknown)
	ErrBannerNotUnique     = errors.New(msg.BannerNotUnique)
)

type Banner struct {
	reader  repo.BannerReader
	saver   repo.BannerSaver
	deleter repo.BannerDeleter
	updater repo.BannerUpdater
	logger  *slog.Logger
}

func NewBannerService(
	reader repo.BannerReader,
	saver repo.BannerSaver,
	deleter repo.BannerDeleter,
	updater repo.BannerUpdater,
	log *slog.Logger,
) *Banner {
	return &Banner{
		reader,
		saver,
		deleter,
		updater,
		log.With(slog.String("comp", "service.banner")),
	}
}

// SaveBanner saves a new banner to the storage.
// It validates the input data and returns an error if the data is invalid.
func (s *Banner) SaveBanner(ctx context.Context, dto banner.CreateDTO) (int64, error) {
	if err := validator.New().Struct(dto); err != nil {
		var validErrs validator.ValidationErrors
		errors.As(err, &validErrs)
		s.logger.Error("request validation failed", sl.Err(err))
		return 0, ValidationErr(validErrs)
	}

	model := dto.ToModel()
	s.logger.Info("saving banner", slog.String("title", model.Title))
	id, err := s.saver.SaveBanner(ctx, model)
	if errors.Is(err, storage.ErrBannerAlreadyExists) {
		s.logger.Error("banner already exists", sl.Err(err))
		return 0, ErrBannerAlreadyExists
	} else if err != nil {
		s.logger.Error("failed to save banner", sl.Err(err))
		return 0, ErrBannerNotSaved
	}

	return id, nil
}

// BannerByFeatureTag returns a banner by the feature and tag ID.
// It also respects the limit and offset parameters.
func (s *Banner) BannerByFeatureTag(
	ctx context.Context,
	featureID, tagID int64,
	limit, offset int,
	useLastRevision, asUser bool,
) (*entity.Banner, error) {
	b, err := s.BannersByFeatureTag(ctx, &featureID, &tagID, &limit, &offset, &useLastRevision)

	if len(b) == 0 {
		s.logger.Info("banner not found")
		return nil, ErrBannerNotFound
	}

	if len(b) > 1 {
		s.logger.Error(
			"banner not unique",
			slog.String("feature_id", strconv.FormatInt(featureID, 10)),
			slog.String("tag_id", strconv.FormatInt(tagID, 10)),
			sl.Err(err),
		)
		return nil, ErrBannerNotUnique
	}

	if asUser {
		if !b[0].IsActive {
			s.logger.Info(fmt.Sprintf("banner %d not active, restricting user access", b[0].ID))
			return nil, ErrBannerNotActive
		}
	}

	return b[0], nil
}

// BannersByFeatureTag returns a list of banners by the feature and tag ID.
// It also respects the limit and offset parameters.
func (s *Banner) BannersByFeatureTag(
	ctx context.Context,
	featureID, tagID *int64,
	limit, offset *int,
	useLastRevision *bool,
) ([]*entity.Banner, error) {
	banners, err := s.reader.BannersByFeatureTag(ctx, featureID, tagID, limit, offset, useLastRevision)
	if err != nil {
		s.logger.Error("failed to get banner by feature and tag", sl.Err(err))
		return nil, ErrBannerNotUnique
	}

	return banners, nil
}

// DeleteBanner deletes a banner by the ID.
// If the banner was not found, it returns an error.
func (s *Banner) DeleteBanner(ctx context.Context, id int64) error {
	if err := validator.New().Var(id, "required"); err != nil {
		var validErrs validator.ValidationErrors
		errors.As(err, &validErrs)
		s.logger.Error("request validation failed", sl.Err(err))
		return ValidationErr(validErrs)
	}

	err := s.deleter.DeleteBanner(ctx, id)
	if errors.Is(err, storage.ErrBannerNotFound) {
		s.logger.Error("banner not found", sl.Err(err))
		return ErrBannerNotFound
	} else if err != nil {
		s.logger.Error("failed to delete banner", sl.Err(err))
		return ErrUnknown
	}

	return nil
}

// UpdateBanner updates a banner by the ID.
// If the banner was not found, it returns an error.
func (s *Banner) UpdateBanner(ctx context.Context, id int64, dto banner.UpdateDTO) error {
	if err := validator.New().Struct(dto); err != nil {
		var validErrs validator.ValidationErrors
		errors.As(err, &validErrs)
		s.logger.Error("request validation failed", sl.Err(err))
		return ValidationErr(validErrs)
	}

	model := dto.ToModel(id)
	s.logger.Info("updating banner", slog.String("id", strconv.FormatInt(id, 10)))
	err := s.updater.UpdateBanner(ctx, model)
	if errors.Is(err, storage.ErrBannerNotFound) {
		s.logger.Error("banner not found", sl.Err(err))
		return ErrBannerNotFound
	} else if err != nil {
		s.logger.Error("failed to update banner", sl.Err(err))
		return ErrUnknown
	}

	return nil
}

func (s *Banner) DeleteBannerByFeatureTag(ctx context.Context, featureID, tagID *int64) error {
	return nil
}
