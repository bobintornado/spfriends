package handlers

import (
	"net/http"
	"time"

	"github.com/bobintornado/spfriends/internal/mid"
	"github.com/bobintornado/spfriends/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(zipkinHost string, apiHost string) http.Handler {
	app := web.New(mid.RequestLogger, mid.ErrorHandler)

	z := NewZipkin(zipkinHost, apiHost, time.Second)
	app.Handle("POST", "/v1/publish", z.Publish)

	h := Health{}
	app.Handle("GET", "/v1/health", h.Check)

	return app
}
