package user

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"travel-crawler/internal/request"
	"travel-crawler/internal/requestcontext"
	"travel-crawler/internal/response"
	"travel-crawler/internal/validation"
)

func GetMe() http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		currentUser := requestcontext.GetUser(r.Context())
		if currentUser == nil {
			response.ErrorUnauthorized(w, "Unauthorized")
			return
		}

		response.SuccessOK(
			w,
			currentUser,
			"Current user retrieved successfully",
		)
	}
}

func Get(
	service Service,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		var req GetRequest

		if err := validation.BindPath(
			r,
			&req,
		); err != nil {
			response.ErrorBadRequest(w, err.Error())
			return
		}

		user, err := service.Get(
			r.Context(),
			req,
		)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				response.ErrorBadRequest(w, "User not found")

			default:
				response.ErrorInternalServer(w, nil)
			}
			return
		}

		response.SuccessOK(
			w,
			user,
			"User retrieved successfully",
		)
	}
}

func Update(
	service Service,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ID, err := request.PathUint64(r, "id")

		if err != nil {
			response.ErrorBadRequest(w, err.Error())
			return
		}

		var req UpdateRequest

		if err := validation.BindJSON(
			r,
			&req,
		); err != nil {
			response.ErrorBadRequest(w, err.Error())
			return
		}

		if err := service.UpdateByID(
			r.Context(),
			ID,
			requestcontext.GetUserID(r.Context()),
			req,
		); err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				response.ErrorBadRequest(w, "User not found")

			default:
				response.ErrorInternalServer(w, nil)
			}
			return
		}

		response.SuccessOK(
			w,
			nil,
			"User updated successfully",
		)
	}
}

func List(
	service Service,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		var req ListRequest

		if err := validation.BindQuery(
			r,
			&req,
		); err != nil {
			response.ErrorBadRequest(w, err.Error())
			return
		}

		req.NormalizeDateRange()

		users, err := service.List(
			r.Context(),
			req,
		)
		if err != nil {
			response.ErrorInternalServer(w, nil)
			return
		}

		response.SuccessOK(
			w,
			users,
			"Users retrieved successfully",
		)
	}
}
