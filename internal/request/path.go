package request

import (
	"github.com/gorilla/mux"
	"net/http"
	"travel-crawler/internal/helper"
)

func PathUint64(
	r *http.Request,
	name string,
) (uint64, error) {
	return helper.ParseUint64(
		mux.Vars(r)[name],
	)
}

func PathString(
	r *http.Request,
	name string,
) string {
	return mux.Vars(r)[name]
}
