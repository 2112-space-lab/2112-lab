package xcontext

import (
	"context"

	log "github.com/Elbujito/2112/src/app-service/pkg/log"
)

const (
	userInfoKey contextKey = "__userInfo"

	userUIDKey         string = "__userUID"
	userGroupUIDKey    string = "__userGroupUID"
	userDisplayNameKey string = "__userName"
	userEmailKey       string = "__userEmail"
)

// UserInfo struct holds the user info including User, Group UIDs, Name and Email
type UserInfo struct {
	UserUID         string
	UserGroupUID    string
	UserDisplayName string
	UserEmail       string
}

// AppendToLogFields adds the user info into the logger
func (info UserInfo) AppendToLogFields(logFields log.Fields) log.Fields {
	logFields[userUIDKey] = info.UserUID
	logFields[userGroupUIDKey] = info.UserGroupUID
	logFields[userDisplayNameKey] = info.UserDisplayName
	logFields[userEmailKey] = info.UserEmail
	return logFields
}

// WithUserInfo adds the user info into the context and logger
func WithUserInfo(userUID string, userGroupUID string, userDisplayName string, userEmail string) ContextEnhancer {
	return func(parentCtx context.Context, logFields log.Fields) (context.Context, log.Fields) {
		user := UserInfo{
			UserUID:         userUID,
			UserGroupUID:    userGroupUID,
			UserDisplayName: userDisplayName,
			UserEmail:       userEmail,
		}

		ctx := context.WithValue(parentCtx, userInfoKey, user)
		logFields = user.AppendToLogFields(logFields)
		return ctx, logFields
	}
}

// ReadUserInfo tries to read the UserInfo from the context and returns a boolean indicating whether the user info was found
func ReadUserInfo(ctx context.Context) (user UserInfo, found bool) {
	return ReadValue(ctx, userInfoKey, UserInfo{})
}
