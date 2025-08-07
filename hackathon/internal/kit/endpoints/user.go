package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/minhho2511/elotusteam-test/internal/kit/services"
	"github.com/minhho2511/elotusteam-test/internal/transforms"
	"github.com/minhho2511/elotusteam-test/utils"
)

type UserEndpoint struct {
	s services.UserSvc
}

func NewUserEndpoint(s services.UserSvc) UserEndpoint {
	return UserEndpoint{
		s: s,
	}
}

func (r UserEndpoint) Register() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transforms.UserReq)
		err = r.s.Register(ctx, req)
		if err != nil {
			return nil, err
		}
		return utils.SetDefaultResponse(ctx, utils.Message{Code: 200}), nil
	}
}

func (r UserEndpoint) Login() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transforms.UserReq)
		resp, err := r.s.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return utils.SetHttpResponse(ctx, utils.Message{Code: 200}, resp, nil), nil
	}
}
