package services

import (
	"context"
	"errors"
	"github.com/minhho2511/elotusteam-test/cfg"
	"github.com/minhho2511/elotusteam-test/internal/middleware"
	"github.com/minhho2511/elotusteam-test/internal/models"
	"github.com/minhho2511/elotusteam-test/internal/transforms"
	"github.com/minhho2511/elotusteam-test/pkgs/clog"
	"github.com/minhho2511/elotusteam-test/utils"
	"github.com/uptrace/bun"
	"time"
)

type FileSvc interface {
	Upload(ctx context.Context, req transforms.FileReq) (models.File, error)
}

func NewFileSvc(db *bun.DB, config cfg.Config, logger clog.Logger) FileSvc {
	return fileSvc{
		db:     db,
		c:      config,
		logger: logger,
	}
}

type fileSvc struct {
	db     *bun.DB
	logger clog.Logger
	c      cfg.Config
}

func (s fileSvc) Upload(ctx context.Context, req transforms.FileReq) (models.File, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	f := models.File{
		Name:     req.FileName,
		FileSize: req.FileSize,
		FileType: req.ContentType,
		FilePath: req.FilePath,
	}
	info := map[string]interface{}{
		"original_name": req.OriginalName,
		"user_agent":    req.UserAgent,
		"ip_address":    req.IPAddress,
		"referer":       req.Referer,
	}
	infoJson, err := utils.MapToJSON(info)
	if err != nil {
		s.logger.Error(err)
		return models.File{}, errors.New("INTERNAL_ERROR")
	}
	f.Info = infoJson
	username, ok := ctx.Value(middleware.UserContextKey).(string)
	if ok {
		var user models.User
		err = s.db.NewSelect().
			Model(&user).
			Where("username = ?", username).
			Scan(ctx)
		if err == nil {
			f.UploadBy = user.ID
		}
	}
	_, err = s.db.NewInsert().
		Model(&f).
		Returning("*").
		Exec(ctx)
	if err != nil {
		s.logger.Error(err)
		return models.File{}, errors.New("INTERNAL_ERROR")
	}
	return f, nil
}
