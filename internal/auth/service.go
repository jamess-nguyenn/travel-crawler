package auth

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/auth0"
	"travel-crawler/internal/database"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/model"
	"travel-crawler/internal/requestcontext"
	"travel-crawler/internal/user"
)

type Service interface {
	Create(
		ctx context.Context,
		actorID uint64,
		req CreateRequest,
	) (*user.User, error)

	Get(
		ctx context.Context,
		auth0ID string,
	) (*requestcontext.User, error)

	SendVerifyEmail(
		ctx context.Context,
		auth0ID string,
	) error
}

type service struct {
	txManager   *database.TxManager
	userRepo    user.UserRepository
	auth0Client auth0.Client
	logger      *zap.SugaredLogger
}

func NewService(
	txManager *database.TxManager,
	userRepo user.UserRepository,
	auth0Client auth0.Client,
	logger *zap.SugaredLogger,
) Service {
	return &service{
		txManager:   txManager,
		userRepo:    userRepo,
		auth0Client: auth0Client,
		logger:      logger,
	}
}

func (s *service) Create(
	ctx context.Context,
	actorID uint64,
	req CreateRequest,
) (*user.User, error) {
	logger := applog.FromContext(s.logger, ctx)

	existingUsers, err := s.userRepo.GetByEmailOrAuth0ID(ctx, req.Email, req.Auth0ID)

	getLogger := logger.With(
		"email", req.Email,
		"auth0_id", req.Auth0ID,
		"data", existingUsers,
	)

	if err != nil {
		getLogger.Errorw(
			"Fail to get users",
			"error", err,
		)

		return nil, err
	}

	getLogger.Info("Get users successfully")

	for _, existingUser := range existingUsers {
		if existingUser.Email == req.Email {
			return nil, user.ErrEmailExists
		}

		if existingUser.Auth0ID == req.Auth0ID {
			return nil, user.ErrAuth0IDExists
		}
	}

	if !req.EmailVerified {
		userResp, err := s.auth0Client.GetUserByAuth0ID(ctx, req.Auth0ID)

		userLogger := logger.With(
			"auth0_id", req.Auth0ID,
			"data", userResp,
		)

		if err != nil {
			userLogger.Errorw(
				"Fail to get user by auth0 id",
				"error", err,
			)

			return nil, err
		}

		userLogger.Info("Get user by auth0 id successfully")

		if !userResp.EmailVerified {
			return nil, user.ErrEmailNotVerified
		}
	}

	now := helper.NowUTC()
	passwordHash := helper.GeneratePassword(100)

	userData := &user.User{
		Model: model.Model{
			CreatedAt: now,
			CreatedBy: actorID,
			UpdatedAt: now,
			UpdatedBy: actorID,
		},
		Email:           req.Email,
		Auth0ID:         req.Auth0ID,
		EmailVerifiedAt: &now,
		PasswordHash:    &passwordHash,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		FullName:        user.BuildFullName(req.FirstName, req.LastName),
		PhoneNumber:     req.PhoneNumber,
		Status:          user.StatusActive,
		AvatarURL:       req.AvatarURL,
	}

	createLogger := applog.WithCreate(
		logger,
		userData,
	)

	if err := s.userRepo.Create(ctx, userData); err != nil {
		createLogger.Errorw(
			"Fail to create user",
			"error", err,
		)

		return nil, err
	}

	createLogger.Info("Create user successfully")

	return userData, nil
}

func (s *service) Get(
	ctx context.Context,
	auth0ID string,
) (*requestcontext.User, error) {
	logger := applog.FromContext(s.logger, ctx)

	conditions := map[string]any{
		"auth0_id": auth0ID,
	}

	userDB, err := s.userRepo.Get(ctx, conditions)

	getLogger := applog.WithGet(
		logger,
		conditions,
		userDB,
	)

	if err != nil {
		getLogger.Errorw(
			"Fail to get user",
			"error", err,
		)

		return nil, err
	}

	getLogger.Info("Get user successfully")

	return &requestcontext.User{
		ID:       userDB.ID,
		Email:    userDB.Email,
		Auth0ID:  userDB.Auth0ID,
		FullName: userDB.FullName,
		Status:   string(userDB.Status),
	}, nil
}

func (s *service) SendVerifyEmail(
	ctx context.Context,
	auth0ID string,
) error {
	logger := applog.FromContext(s.logger, ctx)

	resp, err := s.auth0Client.SendVerificationEmail(ctx, auth0ID)
	if err != nil {
		logger.Errorw(
			"Fail to send verification email",
			"auth0_id", auth0ID,
			"error", err,
		)

		return fmt.Errorf("send verification email: %w", err)
	}

	jobLogger := logger.With(
		"auth0_id", auth0ID,
		"job_id", resp.ID,
		"status", resp.Status,
	)

	if resp.Status != auth0.VerificationEmailJobStatusPending {
		jobLogger.Warn("Unexpected verification email job status")

		return ErrVerificationEmailFailed
	}

	jobLogger.Info("Verification email job created")

	return nil
}
