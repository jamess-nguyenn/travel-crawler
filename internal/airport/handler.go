package airport

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"travel-crawler/internal/response"
	"travel-crawler/internal/validation"
)

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

		airport, err := service.Get(
			r.Context(),
			req,
		)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				response.ErrorBadRequest(w, "Airport not found")

			default:
				response.ErrorInternalServer(w, nil)
			}
			return
		}

		response.SuccessOK(
			w,
			airport,
			"Airport retrieved successfully",
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

		airports, err := service.List(
			r.Context(),
			req,
		)
		if err != nil {
			response.ErrorInternalServer(w, nil)
			return
		}

		response.SuccessOK(
			w,
			airports,
			"Airports retrieved successfully",
		)
	}
}
