package handlers

import (
	"context"
	"net/http"

	"github.com/bobintornado/spfriends/internal/platform/db"
	"github.com/bobintornado/spfriends/internal/platform/web"
	"github.com/bobintornado/spfriends/internal/relationship"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Relationship represents the Relationship API method handler set.
// which in charge of relationships like friendship, subscription and block
type Relationship struct {
	MasterDB *db.DB
}

type RelationshipResponse struct {
	Success    bool     `json:"success,omitempty"`
	Friends    []string `json:"friends,omitempty"`
	Count      int      `json:"count,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
}

func (rel *Relationship) CreateFriendship(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Relationship.CreateFriendship")
	defer span.End()

	var createFriendshipRequest relationship.CreateFriendshipRequest
	if err := web.Unmarshal(r.Body, &createFriendshipRequest); err != nil {
		return errors.Wrap(err, "error unmarshaling request")
	}

	dbConn, err := rel.MasterDB.OpenPool()
	if err != nil {
		return errors.Wrap(err, "can't get connection from pool")
	}
	defer (*dbConn).Close()

	err = relationship.CreateFriendship(ctx, dbConn, &createFriendshipRequest)
	if err != nil {
		return errors.Wrapf(err, "friendship request: %+v", &createFriendshipRequest)
	}

	res := RelationshipResponse{
		Success: true,
	}

	web.Respond(ctx, w, res, http.StatusCreated)
	return nil
}

func (rel *Relationship) GetFriendsOfUser(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Relationship.GetFriendsOfUser")
	defer span.End()

	var getFriendsListRequest relationship.GetFriendsListRequest
	if err := web.Unmarshal(r.Body, &getFriendsListRequest); err != nil {
		return errors.Wrap(err, "")
	}

	dbConn, err := rel.MasterDB.OpenPool()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer (*dbConn).Close()

	friendList, err := relationship.ListFriendsOfUser(ctx, dbConn, &getFriendsListRequest)
	if err != nil {
		return errors.Wrapf(err, "request: %+v", &getFriendsListRequest)
	}

	res := RelationshipResponse{
		Success: true,
		Friends: friendList,
		Count:   len(friendList),
	}

	web.Respond(ctx, w, res, http.StatusOK)
	return nil
}

func (rel *Relationship) GetCommonFriends(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Relationship.GetCommonFriends")
	defer span.End()

	var getCommonFriendsListRequest relationship.GetCommonFriendsListRequest
	if err := web.Unmarshal(r.Body, &getCommonFriendsListRequest); err != nil {
		return errors.Wrap(err, "")
	}

	dbConn, err := rel.MasterDB.OpenPool()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer (*dbConn).Close()

	friendList, err := relationship.ListCommonFriends(ctx, dbConn, &getCommonFriendsListRequest)
	if err != nil {
		return errors.Wrapf(err, "commonfriends: %+v", &getCommonFriendsListRequest)
	}

	res := RelationshipResponse{
		Success: true,
		Friends: friendList,
		Count:   len(friendList),
	}

	web.Respond(ctx, w, res, http.StatusOK)
	return nil
}

func (rel *Relationship) CreateSubscription(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Relationship.CreateSubscription")
	defer span.End()

	var createSubscriptionRequest relationship.CreateSubscriptionRequest
	if err := web.Unmarshal(r.Body, &createSubscriptionRequest); err != nil {
		return errors.Wrap(err, "")
	}

	dbConn, err := rel.MasterDB.OpenPool()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer (*dbConn).Close()

	err = relationship.CreateSubscription(ctx, dbConn, &createSubscriptionRequest)
	if err != nil {
		return errors.Wrapf(err, "subscription: %+v", &createSubscriptionRequest)
	}

	res := RelationshipResponse{
		Success: true,
	}

	web.Respond(ctx, w, res, http.StatusCreated)
	return nil
}

func (rel *Relationship) CreateBlock(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Relationship.CreateBlock")
	defer span.End()

	var createBlockRequest relationship.CreateBlockRequest
	if err := web.Unmarshal(r.Body, &createBlockRequest); err != nil {
		return errors.Wrap(err, "")
	}

	dbConn, err := rel.MasterDB.OpenPool()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer (*dbConn).Close()

	err = relationship.CreateBlock(ctx, dbConn, &createBlockRequest)
	if err != nil {
		return errors.Wrapf(err, "block: %+v", &createBlockRequest)
	}

	res := RelationshipResponse{
		Success: true,
	}

	web.Respond(ctx, w, res, http.StatusCreated)
	return nil
}

func (rel *Relationship) GetUpdateList(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Relationship.GetUpdateList")
	defer span.End()

	var getUpdateListRequest relationship.GetUpdateListRequest
	if err := web.Unmarshal(r.Body, &getUpdateListRequest); err != nil {
		return errors.Wrap(err, "")
	}

	dbConn, err := rel.MasterDB.OpenPool()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer (*dbConn).Close()

	list, err := relationship.GetUpdateList(ctx, dbConn, &getUpdateListRequest)
	if err != nil {
		return errors.Wrapf(err, "block: %+v", &getUpdateListRequest)
	}

	res := RelationshipResponse{
		Success:    true,
		Recipients: list,
	}

	web.Respond(ctx, w, res, http.StatusOK)
	return nil
}
