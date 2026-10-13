package auth

import (
	"errors"
	"github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"net/http"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/response"
	"travel-crawler/internal/user"
	"travel-crawler/internal/validation"
)

func Verify(
	service Service,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ctx := r.Context()

		claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](ctx)
		if err != nil {
			response.ErrorUnauthorized(w, "Invalid token claims")
			return
		}

		customClaims, ok := claims.CustomClaims.(*CustomClaims)
		if !ok || customClaims == nil {
			response.ErrorUnauthorized(w, "Invalid custom claims")
			return
		}

		req := CreateRequest{
			Email:         customClaims.Email,
			Auth0ID:       claims.RegisteredClaims.Subject,
			EmailVerified: customClaims.EmailVerified,
			FirstName:     customClaims.UserMetadata.FirstName,
			LastName:      customClaims.UserMetadata.LastName,
			PhoneNumber:   customClaims.UserMetadata.PhoneNumber,
			AvatarURL:     &customClaims.Picture,
		}

		if err := validation.Validate(req); err != nil {
			response.ErrorBadRequest(w, err.Error())
			return
		}

		userData, err := service.Create(ctx, constant.SystemUserID, req)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrEmailExists) ||
				errors.Is(err, user.ErrAuth0IDExists):
				response.ErrorConflict(w, err.Error())

			case errors.Is(err, user.ErrEmailNotVerified):
				response.ErrorForbidden(w, err.Error())

			default:
				response.ErrorInternalServer(w, nil)
			}
			return
		}

		response.SuccessCreated(
			w,
			userData,
			"User created successfully",
		)
	}
}

func SendVerifyEmail(
	service Service,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ctx := r.Context()

		claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](ctx)
		if err != nil {
			response.ErrorUnauthorized(w, "Invalid token claims")
			return
		}

		auth0ID := claims.RegisteredClaims.Subject
		if auth0ID == "" {
			response.ErrorUnauthorized(w, "Missing token subject")
			return
		}

		if err := service.SendVerifyEmail(ctx, auth0ID); err != nil {
			switch {
			case errors.Is(err, ErrVerificationEmailFailed):
				response.ErrorBadRequest(w, err.Error())

			default:
				response.ErrorInternalServer(w, nil)
			}
			return
		}

		response.SuccessOK(
			w,
			nil,
			"Verification email accepted",
		)
	}
}
