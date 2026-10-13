package requestcontext

import (
	"context"
)

const userKey contextKey = "user"

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Auth0ID  string `json:"auth0_id"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
}

func WithUser(
	ctx context.Context,
	user *User,
) context.Context {
	return context.WithValue(
		ctx,
		userKey,
		user,
	)
}

func GetUser(
	ctx context.Context,
) *User {
	user, ok := ctx.Value(userKey).(*User)
	if !ok {
		return nil
	}

	return user
}

func GetUserID(
	ctx context.Context,
) uint64 {
	user := GetUser(ctx)
	if user == nil {
		return 0
	}

	return user.ID
}
