package transports

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/minhho2511/elotusteam-test/internal/kit/endpoints"
	"github.com/minhho2511/elotusteam-test/internal/kit/services"
	"github.com/minhho2511/elotusteam-test/internal/transforms"
	"github.com/minhho2511/elotusteam-test/pkgs/clog"
	"github.com/minhho2511/elotusteam-test/utils"
	"net/http"
)

func decodeUser(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req transforms.UserReq
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	utils.HtmlEscape(&req)
	return req, nil
}

func UserHttpHandler(userSvc services.UserSvc, logger clog.Logger) http.Handler {
	pr := mux.NewRouter()

	user := endpoints.NewUserEndpoint(userSvc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(logger),
		httptransport.ServerErrorEncoder(utils.EncodeError),
	}

	pr.Methods("POST").Path("/user/login").Handler(httptransport.NewServer(
		user.Login(),
		decodeUser,
		utils.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path("/user/register").Handler(httptransport.NewServer(
		user.Register(),
		decodeUser,
		utils.EncodeResponseHTTP,
		options...,
	))

	return pr
}
