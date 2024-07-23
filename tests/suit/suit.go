package suit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"banners-management/internal/app"
	"banners-management/internal/config"
	"banners-management/internal/lib/jwt"
	slogdiscard "banners-management/internal/lib/logger/slogimpl"
	"banners-management/internal/service/banner"
	"banners-management/internal/storage/pgs"
	"banners-management/migrator"
)

var (
	once sync.Once
	suit *Suit
)

type Suit struct {
	Cfg        *config.Config
	JwtManager *jwt.Manager
}

func Setup(t *testing.T) *Suit {
	t.Helper()

	//goland:noinspection GoVetLostCancel
	once.Do(func() {
		t.Helper()
		// cancel function is discarded because this context is used
		// by the server, and should live as long as context.Background
		// but no longer than one minute (time for all tests to finish)
		ctx, _ := context.WithTimeout(context.Background(), time.Minute) //nolint:govet // why: see above
		args := []string{"-migrations-path", "../migrations", "-direction", "up"}
		getenv := func(s string) (string, bool) {
			if s == "CONFIG_PATH" {
				return "../config/local.test.json", true
			}

			return "", false
		}

		migrateDown := func() {
			args := []string{"-migrations-path", "../migrations", "-direction", "down"}
			migrator.Migrate(context.Background(), io.Discard, args, getenv, "postgres")
		}

		migrateDown()

		// migrate db up
		migrator.Migrate(context.Background(), io.Discard, args, getenv, "postgres")

		// start server
		cfg := config.MustLoad([]string{}, getenv)
		s, err := pgs.New(ctx, cfg.DB.ConnectionString())
		if err != nil {
			panic(err)
		}
		l := slogdiscard.NewDiscardLogger()
		j := jwt.NewManager(string(cfg.JwtSettings.SecretKey), time.Duration(cfg.JwtSettings.Expire))
		b := banner.NewService(s, s, s, s, l)
		a := app.New(l, j, b)
		go app.RunWithConfig(ctx, []string{}, getenv, a)

		// wait for server to be ready (GET /health)
		err = waitTilReady(ctx, 5*time.Second, fmt.Sprintf("http://%s/health", cfg.HTTPServer.Address))
		if err != nil {
			panic(err)
		}

		suit = &Suit{cfg, j}
	})

	return suit
}

// waitTilReady waits til the provided endpoint becomes available (returns http.StatusOK).
// If it wasn't become available within timeout, it returns error.
func waitTilReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			_ = resp.Body.Close()
			return nil
		}
		_ = resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return errors.New("timeout reached while waiting for endpoint")
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}
