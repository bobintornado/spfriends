package handlers

import (
	"net/http"

	"github.com/bobintornado/spfriends/internal/mid"
	"github.com/bobintornado/spfriends/internal/platform/db"
	"github.com/bobintornado/spfriends/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(masterDB *db.DB) http.Handler {
	app := web.New(mid.RequestLogger, mid.Metrics, mid.ErrorHandler)

	// a set of methods to operate on relationships
	rel := Relationship{
		MasterDB: masterDB,
	}
	app.Handle("POST", "/v1/createFriendship", rel.CreateFriendship)
	app.Handle("POST", "/v1/getFriendsOfUser", rel.GetFriendsOfUser)
	app.Handle("POST", "/v1/getCommonFriends", rel.GetCommonFriends)
	app.Handle("POST", "/v1/subscribe", rel.CreateSubscription)
	app.Handle("POST", "/v1/block", rel.CreateBlock)
	app.Handle("POST", "/v1/getUpdateList", rel.GetUpdateList)

	h := Health{
		MasterDB: masterDB,
	}
	app.Handle("GET", "/v1/health", h.Check)

	return app
}
