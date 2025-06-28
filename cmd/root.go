package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/EdisonTantra/lemonPajak/config"
	lemonPort "github.com/EdisonTantra/lemonPajak/internal/core/port"
	lemonPsql "github.com/EdisonTantra/lemonPajak/internal/repository/postgres"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/telemetry"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// PostgreSQL driver
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "service",
	Short: "lemon service",
}

func initConfig() config.Config {
	var cfg config.Config

	f := func() error {
		// required: base config file
		viper.AddConfigPath("config")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		// optional: config file for local development
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		if err := godotenv.Load("config/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Config file changed:", e.Name)
		})

		return viper.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = "json"
		})
	}

	err := f()
	if err != nil {
		panic(err)
	}

	return cfg
}

func initLogger() (context.Context, logat.AppsLogger) {
	ctx, logger := logat.New()
	return ctx, logger
}

func initOpTel(
	ctx context.Context,
	cfgService config.Service,
	cfgTele config.Telemetry,
) (shutdown func(context.Context) error, err error) {
	shutdownCtx, err := telemetry.SetupOpTelSDK(ctx, telemetry.ConfigOptel{
		ServiceCode:    cfgService.Code,
		ServiceVersion: cfgService.Version,
		CollectorURL:   cfgTele.CollectorURL,
		SecureMode:     cfgTele.SecureMode,
	})

	return shutdownCtx, err
}

func initHTTPServer(cfg *config.ServerHTTP, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Address,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}
}

func initRepoPostgres(ctx context.Context, cfg *config.Postgresql) (lemonPort.AppsRepository, error) {
	if cfg == nil {
		return nil, errors.New("error config postgres nil")
	}

	opts := lemonPsql.NewRepoOptions{
		Username:    cfg.Username,
		Password:    cfg.Password,
		Host:        cfg.Host,
		Port:        cfg.Port,
		Database:    cfg.Database,
		SSLMode:     cfg.SSLMode,
		MaxIdleConn: cfg.MaxIdleConn,
		MaxOpenConn: cfg.MaxOpenConn,
		MaxIdleTime: cfg.MaxIdleTime,
	}

	db, err := lemonPsql.New(ctx, opts)
	if err != nil {
		return nil, err
	}

	return db, nil
}
