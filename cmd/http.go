package cmd

import (
	"fmt"
	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	svcApp "github.com/EdisonTantra/lemonPajak/internal/core/service/application"
	svcUser "github.com/EdisonTantra/lemonPajak/internal/core/service/user"
	"github.com/EdisonTantra/lemonPajak/internal/repository/externalapi/djp"
	lemonHTTP "github.com/EdisonTantra/lemonPajak/internal/transport/http"
	httpHandlerApp "github.com/EdisonTantra/lemonPajak/internal/transport/http/handlers/application"
	httpHandlerUser "github.com/EdisonTantra/lemonPajak/internal/transport/http/handlers/user"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP Server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := initConfig()
		ctx, logger := initLogger()

		// repos
		repoPsql, err := initRepoPostgres(ctx, cfg.Postgres)
		if err != nil {
			logger.Fatal(
				ctx, fmt.Sprintf("error init psql postgres: %v\n", err),
				cons.EventLogNameRoot, err,
			)
		}

		// api clients
		djpClient := djp.New(cfg.DJPClient.BaseURL)

		// register stores to repo
		repoPsql.RegisterStore()

		// services
		userSvc := svcUser.New(repoPsql.GetUserStore())
		appSvc := svcApp.New(djpClient)

		// handlers
		handlerHTTPUser := httpHandlerUser.New(&httpHandlerUser.HandlerOpts{
			SvcUser: userSvc,
		})

		handlerHTTPApplication := httpHandlerApp.New(&httpHandlerApp.HandlerOpts{
			SvcApp: appSvc,
		})

		// routing handlers
		r := lemonHTTP.NewRouter(
			logger,
			handlerHTTPApplication,
			handlerHTTPUser,
		)
		s := initHTTPServer(cfg.ServerHTTP, r.Handlers(cfg.Service.Code))

		// running http server
		lock := make(chan error)
		go func(lock chan error) {
			lock <- s.ListenAndServe()
		}(lock)

		logger.Info(ctx, fmt.Sprintf("running at %s", s.Addr), cons.EventLogNameRoot, nil)
		err = <-lock
		if err != nil {
			_ = repoPsql.Close()
			_ = s.Close()

			logger.Fatal(ctx, "error graceful shutdown repo: %v\n", cons.EventLogNameRoot, err)
		}
	},
}
