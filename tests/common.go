package tests

import (
	"net/url"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"avito-test-task/internal/models/dto/banner"
	"avito-test-task/tests/suit"
)

var (
	muTag         sync.Mutex
	muFeature     sync.Mutex
	lastTagID     int64 = 0
	lastFeatureID int64 = 0

	once               sync.Once
	expect             *httpexpect.Expect
	tokenUsr, tokenAdm string
)

func initTest(t *testing.T) (*httpexpect.Expect, string, string) {
	t.Helper()
	t.Parallel()
	once.Do(func() {
		t.Helper()

		s := suit.Setup(t)
		u := url.URL{
			Scheme: "http",
			Host:   s.Cfg.HTTPServer.Address,
		}

		expect = httpexpect.Default(t, u.String())

		rqr := require.New(t)
		user, err := s.JwtManager.GenerateToken("user")
		rqr.NoError(err)
		admin, err := s.JwtManager.GenerateToken("admin")
		rqr.NoError(err)
		tokenUsr, tokenAdm = user, admin
	})

	return expect, tokenUsr, tokenAdm
}

func rawToInt64(f interface{}) int64 {
	return int64(f.(float64))
}

// getNextTagIDs returns the next tag IDs.
// It is unique for each call. It is thread-safe.
// The uniqueness is needed to avoid conflicts with the database trigger.
func getNextTagIDs(count int) []int64 {
	res := make([]int64, count)
	muTag.Lock()
	defer muTag.Unlock()
	for i := range count {
		lastTagID++
		res[i] = lastTagID
	}
	return res
}

// getNextFeatureID returns the next feature ID.
// It is unique for each call. It is thread-safe.
// The uniqueness is needed to avoid conflicts with the database trigger.
func getNextFeatureID() int64 {
	muFeature.Lock()
	defer muFeature.Unlock()
	lastFeatureID++
	return lastFeatureID
}

func newCreateBannerDTO() banner.CreateDTO {
	return banner.CreateDTO{
		TagIDs:    getNextTagIDs(2),
		FeatureID: getNextFeatureID(),
		Content: banner.CreateContent{
			Title: gofakeit.Word(),
			Text:  gofakeit.Word(),
			URL:   gofakeit.URL(),
		},
		IsActive: true,
	}
}

func createBannerDTO(featureID int64, tagIDs []int64, isActive bool) banner.CreateDTO {
	return banner.CreateDTO{
		TagIDs:    tagIDs,
		FeatureID: featureID,
		Content: banner.CreateContent{
			Title: gofakeit.Word(),
			Text:  gofakeit.Word(),
			URL:   gofakeit.URL(),
		},
		IsActive: isActive,
	}
}

func newUpdateBannerDTO() banner.UpdateDTO {
	tagIDs := getNextTagIDs(2)
	title := gofakeit.Word()
	text := gofakeit.Word()
	u := gofakeit.URL()
	featureID := getNextFeatureID()
	isActive := true

	return banner.UpdateDTO{
		TagIDs:    &tagIDs,
		FeatureID: &featureID,
		Content: &banner.UpdateContent{
			Title: &title,
			Text:  &text,
			URL:   &u,
		},
		IsActive: &isActive,
	}
}

func updateBannerDTO(featureID *int64, tagIDs *[]int64, isActive *bool) banner.UpdateDTO {
	title := gofakeit.Word()
	text := gofakeit.Word()
	u := gofakeit.URL()

	return banner.UpdateDTO{
		TagIDs:    tagIDs,
		FeatureID: featureID,
		Content: &banner.UpdateContent{
			Title: &title,
			Text:  &text,
			URL:   &u,
		},
		IsActive: isActive,
	}
}
