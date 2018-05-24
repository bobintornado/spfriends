package relationship

// CreateUser contains information needed to create a friendship
// eg "{friends:['andy@example.com','john@example.com']}"
type CreateFriendshipRequest struct {
	Friends [2]string `json:"friends" validate:"required"`
}

type GetFriendsListRequest struct {
	Email string `json:"email" validate:"required"`
}

type GetCommonFriendsListRequest struct {
	Friends [2]string `json:"friends" validate:"required"`
}

type CreateSubscriptionRequest struct {
	Requestor string `json:"requestor" validate:"required"`
	Target    string `json:"target" validate:"required"`
}

type CreateBlockRequest struct {
	Requestor string `json:"requestor" validate:"required"`
	Target    string `json:"target" validate:"required"`
}

type GetUpdateListRequest struct {
	Sender string `json:"sender" validate:"required"`
	Text   string `json:"text" validate:"required"`
}
