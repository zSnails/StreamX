package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zSnails/streamx/pkg/logging"
)

var log = logging.Get().WithField("service", "middleware")

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithContext(r.Context()).WithFields(logrus.Fields{
			"request-uri": r.RequestURI,
			"method":      r.Method,
			"user":        r.URL.User,
			"user-agent":  r.Header.Get("User-Agent"),
		}).Infof("request from %s\n", r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}
