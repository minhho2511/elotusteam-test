package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/minhho2511/elotusteam-test/internal/models"
	"github.com/minhho2511/elotusteam-test/internal/transforms"
	"github.com/minhho2511/elotusteam-test/utils"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserSvc interface {
	Register(ctx context.Context, req transforms.UserReq) error
	Login(ctx context.Context, req transforms.UserReq) (transforms.LoginResp, error)
}

func NewUserSvc(db *bun.DB, logger log.Logger, jwtSvc *utils.JWTService) UserSvc {
	return userSvc{
		db:     db,
		logger: logger,
		jwtSvc: jwtSvc,
	}
}

type userSvc struct {
	db     *bun.DB
	logger log.Logger
	jwtSvc *utils.JWTService
}

func (u userSvc) userExist(ctx context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	exists, err := u.db.NewSelect().
		Model((*models.User)(nil)).
		Where("username = ?", username).
		Exists(ctx)

	if err != nil {
		_ = u.logger.Log("err", err)
		return false, err
	}
	return exists, nil
}

func (u userSvc) Register(ctx context.Context, req transforms.UserReq) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	exists, err := u.userExist(ctx, req.Username)
	if err != nil {
		_ = u.logger.Log("err", err)
		return errors.New("INTERNAL_ERROR")
	}
	if exists {
		return errors.New("USER_WITH_USERNAME_EXISTS")
	}

	user := models.User{
		UserName: req.Username,
	}
	pass, err := utils.HashPassword(req.Password)
	if err != nil {
		_ = u.logger.Log("err", err)
		return err
	}
	user.Password = pass
	_, err = u.db.NewInsert().
		Model(&user).
		Returning("*").
		Exec(ctx)
	if err != nil {
		_ = u.logger.Log("err", err)
		return errors.New("INTERNAL_ERROR")
	}
	return nil
}

func (u userSvc) Login(ctx context.Context, req transforms.UserReq) (transforms.LoginResp, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var user models.User
	err := u.db.NewSelect().
		Model(&user).
		Where("username = ?", req.Username).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transforms.LoginResp{}, errors.New("USER_NOT_FOUND")
		}
		_ = u.logger.Log("err", err)
		return transforms.LoginResp{}, errors.New("INTERNAL_ERROR")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		_ = u.logger.Log("err", err)
		return transforms.LoginResp{}, errors.New("USERNAME_OR_PASSWORD_INVALID")
	}
	token, err := u.jwtSvc.GenerateJWT(user.UserName, 24*time.Hour)
	if err != nil {
		_ = u.logger.Log("err", err)
		return transforms.LoginResp{}, errors.New("INTERNAL_ERROR")
	}
	resp := transforms.LoginResp{
		Username: user.UserName,
		Token:    token,
	}
	return resp, nil
}
