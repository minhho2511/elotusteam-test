package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/minhho2511/elotusteam-test/cfg"
	"github.com/minhho2511/elotusteam-test/internal/initialization"
	"github.com/minhho2511/elotusteam-test/internal/middleware"
	"github.com/minhho2511/elotusteam-test/migrations"
	"github.com/minhho2511/elotusteam-test/pkgs/cache"
	"github.com/minhho2511/elotusteam-test/pkgs/clog"
	"github.com/minhho2511/elotusteam-test/pkgs/db"
	"github.com/minhho2511/elotusteam-test/utils"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	lg := clog.NewLogger(clog.NewGoKitLog())

	c := cfg.LoadConfig()
	dbConfig := db.PGConfig{
		DB:      c.DB.DBName,
		User:    c.DB.DBUser,
		Pass:    c.DB.DBPass,
		Host:    c.DB.DBHost,
		Port:    c.DB.DBPort,
		SslSkip: true,
	}
	if c.AppEnv == "dev" {
		dbConfig.QueryDebug = true
	}
	pg, err := db.MakePGConnect(dbConfig)
	if err != nil {
		panic(err)
	}
	defer pg.Close()
	lists := migrations.MigrationLists()
	migration := db.NewMigrationTool(pg)
	migration.Migrate(lists)
	_ = level.Info(lg).Log("msg", "Migrate successfully !!!")

	rdb, err := cache.NewRedis(cache.Config{
		Hosts:     strings.Split(c.Cache.CacheHost, ","),
		Pass:      c.Cache.CachePass,
		DB:        c.Cache.CacheDB,
		Debug:     true,
		IsCluster: false,
	}, lg)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	jwtSvc := utils.NewJWTService(c.JWTSecret, rdb)
	mux := initialization.Routing(pg, lg, jwtSvc, c)
	httpServer := http.Server{
		Addr:    *flag.String("listen", ":"+c.HttpPort, "Listen address."),
		Handler: middleware.AppMiddleware(mux, lg),
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		_ = lg.Log("msg", "start graceful shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if e := httpServer.Shutdown(ctx); e != nil {
			panic(e)
		}
		close(idleConnectionsClosed)
	}()

	_ = lg.Log("msg", fmt.Sprintf("Listening at port %s", c.HttpPort))
	if err = httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
	<-idleConnectionsClosed
}
