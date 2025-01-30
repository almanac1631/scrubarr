package scrubarr

import (
	"context"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log/slog"
	"net/http"
	"os"
)

var (
	version, commit string
)

var k = koanf.New(".")

func StartApp() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	err := LoadConfig("./test/real_test_config.toml")
	if err != nil {
		panic(err)
	}
	slog.Info("starting scrubarr...", "version", version, "commit", commit)
	err = registerRetrievers(k)
	//retrieverRegistry.RefreshCachedEntryMapping()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = webserver.StartWebserver(ctx, k, retrieverRegistry)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func LoadConfig(configPath string) error {
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	return nil
}
