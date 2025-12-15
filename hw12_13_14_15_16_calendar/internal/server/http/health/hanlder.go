package health

import (
	"net/http"
)

func NewHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = writer.Write([]byte(""))
	}
}
