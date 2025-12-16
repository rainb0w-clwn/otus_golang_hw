package log

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/_common"
)

func NewHandler(logger common.Logger, next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		lrw := &loggingResponseWriter{writer, http.StatusOK}

		start := time.Now()
		next.ServeHTTP(writer, request)
		end := time.Since(start)

		logJSON, err := json.Marshal(
			struct {
				IP        string
				Datetime  string
				Method    string
				Path      string
				HTTP      string
				Status    string
				Time      string
				UserAgent string
			}{
				IP:        request.RemoteAddr,
				Datetime:  time.Now().Format(time.RFC822),
				Method:    request.Method,
				Path:      request.URL.Path,
				HTTP:      request.Proto,
				Status:    strconv.Itoa(lrw.StatusCode),
				Time:      end.String(),
				UserAgent: request.UserAgent(),
			},
		)
		if err != nil {
			logger.Error(err.Error())
		}

		logger.Info(string(logJSON))
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
