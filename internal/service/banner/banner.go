package banner

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/go-playground/validator/v10"

	"avito-test-task/internal/lib/api/msg"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/models/dto/banner"
	"avito-test-task/internal/models/entity"
	"avito-test-task/internal/service"
	"avito-test-task/internal/storage/repo"
)

var (
	ErrNotSaved      = errors.New(msg.ErrSavingBanner)
	ErrAlreadyExists = errors.New(msg.BannerAlreadyExists)
	ErrNotFound      = errors.New(msg.BannerNotFound)
	ErrNotActive     = errors.New(msg.BannerNotActive)
	ErrUnknown       = errors.New(msg.ErrUnknown)
	ErrNotUnique     = errors.New(msg.BannerNotUnique)
)

var (
	validatr = validator.New()
)

type Service struct {
	reader  repo.BannerReader
	saver   repo.BannerSaver
	deleter repo.BannerDeleter
	updater repo.BannerUpdater
	logger  *slog.Logger
}

func NewService(
	reader repo.BannerReader,
	saver repo.BannerSaver,
	deleter repo.BannerDeleter,
	updater repo.BannerUpdater,
	log *slog.Logger,
) *Service {
	return &Service{
		reader,
		saver,
		deleter,
		updater,
		log.With(slog.String("comp", "service.banner.NewService")),
	}
}

// SaveBanner saves a new banner to the storage.
// It validates the input data and returns an error if the data is invalid.
func (s *Service) SaveBanner(ctx context.Context, dto banner.CreateDTO) (int64, error) {
	if err := validatr.Struct(dto); err != nil {
		var validErrs validator.ValidationErrors
		errors.As(err, &validErrs)
		s.logger.Error("request validation failed", sl.Err(err))
		return 0, service.ValidationErr(validErrs)
	}

	model := dto.ToModel()
	s.logger.Info("saving banner", slog.String("title", model.Title))
	id, err := s.saver.SaveBanner(ctx, model)
	if errors.Is(err, repo.ErrBannerAlreadyExists) {
		s.logger.Error("banner already exists", sl.Err(err))
		return 0, ErrAlreadyExists
	} else if err != nil {
		s.logger.Error("failed to save banner", sl.Err(err))
		return 0, ErrNotSaved
	}

	return id, nil
}

// BannerByFeatureTag returns a banner by the feature and tag ID.
// It also respects the limit and offset parameters.
func (s *Service) BannerByFeatureTag(
	ctx context.Context,
	featureID, tagID int64,
	useLastRevision, asUser bool,
) (*entity.Banner, error) {
	b, err := s.reader.BannerByFeatureTag(ctx, featureID, tagID, useLastRevision)
	if err != nil || b == nil {
		if errors.Is(err, repo.ErrBannerNotFound) {
			s.logger.Info("banner not found",
				slog.Int64("featureID", featureID), slog.Int64("tagID", tagID))
			return nil, ErrNotFound
		} else if errors.Is(err, repo.ErrBannerNotUnique) {
			s.logger.Info("banner not unique",
				slog.Int64("featureID", featureID), slog.Int64("tagID", tagID))
			return nil, ErrNotUnique
		}
		s.logger.Error("failed to get banner by feature and tag",
			sl.Err(err), slog.Int64("featureID", featureID), slog.Int64("tagID", tagID))
		return nil, ErrUnknown
	}

	if asUser {
		if !b.IsActive {
			s.logger.Info("banner not active, restricting user access", slog.Int64("id", b.ID))
			return nil, ErrNotActive
		}
	}

	return b, nil
}

// BannersByFeatureTag returns a list of banners by the feature and tag ID.
// It also respects the limit and offset parameters.
func (s *Service) BannersByFeatureTag(
	ctx context.Context,
	featureID, tagID *int64,
	limit, offset *int,
	useLastRevision *bool,
) ([]*entity.Banner, error) {
	banners, err := s.reader.BannersByFeatureTag(ctx, featureID, tagID, limit, offset, useLastRevision)
	if err != nil {
		s.logger.Error("failed to get banner by feature and tag", sl.Err(err))
		return nil, ErrUnknown
	}

	return banners, nil
}

// DeleteBanner deletes a banner by the ID.
// If the banner was not found, it returns an error.
func (s *Service) DeleteBanner(ctx context.Context, id int64) error {
	if err := validatr.Var(id, "required"); err != nil {
		var validErrs validator.ValidationErrors
		errors.As(err, &validErrs)
		s.logger.Error("request validation failed", sl.Err(err))
		return service.ValidationErr(validErrs)
	}

	err := s.deleter.DeleteBanner(ctx, id)
	if errors.Is(err, repo.ErrBannerNotFound) {
		s.logger.Error("banner not found", sl.Err(err))
		return ErrNotFound
	} else if err != nil {
		s.logger.Error("failed to delete banner", sl.Err(err))
		return ErrUnknown
	}

	return nil
}

// UpdateBanner updates a banner by the ID.
// If the banner was not found, it returns an error.
func (s *Service) UpdateBanner(ctx context.Context, id int64, dto banner.UpdateDTO) error {
	if err := validatr.Struct(dto); err != nil {
		var validErrs validator.ValidationErrors
		errors.As(err, &validErrs)
		s.logger.Error("request validation failed", sl.Err(err))
		return service.ValidationErr(validErrs)
	}

	model := dto.ToModel(id)
	s.logger.Info("updating banner", slog.String("id", strconv.FormatInt(id, 10)))
	err := s.updater.UpdateBanner(ctx, model)
	if errors.Is(err, repo.ErrBannerNotFound) {
		s.logger.Error("banner not found", sl.Err(err))
		return ErrNotFound
	} else if errors.Is(err, repo.ErrBannerAlreadyExists) {
		s.logger.Error("unable to update banner", sl.Err(err))
		return ErrAlreadyExists
	} else if err != nil {
		s.logger.Error("failed to update banner", sl.Err(err))
		return ErrUnknown
	}

	return nil
}

func (s *Service) DeleteBannerByFeatureTag(ctx context.Context, featureID, tagID *int64) error {
	return nil
}
