package initialization

import (
	"github.com/minhho2511/elotusteam-test/cfg"
	"github.com/minhho2511/elotusteam-test/internal/kit/services"
	"github.com/minhho2511/elotusteam-test/internal/kit/transports"
	"github.com/minhho2511/elotusteam-test/pkgs/clog"
	"github.com/minhho2511/elotusteam-test/utils"
	"github.com/uptrace/bun"
	"net/http"
)

func Routing(db *bun.DB, logger clog.Logger, jwtSvc *utils.JWTService, c cfg.Config) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/__health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	userSvc := services.NewUserSvc(db, logger, jwtSvc)
	userHttp := transports.UserHttpHandler(userSvc, logger)

	fileSvc := services.NewFileSvc(db, c, logger)
	fileHttp := transports.FileHttpHandler(fileSvc, logger, c, jwtSvc)
	mux.Handle("/user/", userHttp)
	mux.Handle("/file/", fileHttp)

	return mux
}
