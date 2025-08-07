package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/minhho2511/elotusteam-test/internal/kit/services"
	"github.com/minhho2511/elotusteam-test/internal/transforms"
	"github.com/minhho2511/elotusteam-test/utils"
)

type FileEndpoint struct {
	s services.FileSvc
}

func NewFileEndpoint(s services.FileSvc) FileEndpoint {
	return FileEndpoint{
		s: s,
	}
}

func (r FileEndpoint) Upload() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transforms.FileReq)
		file, err := r.s.Upload(ctx, req)
		if err != nil {
			return nil, err
		}
		return utils.SetHttpResponse(ctx, utils.Message{Code: 200}, file, nil), nil
	}
}
