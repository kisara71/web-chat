package profile

import "web-chat/internal/model"

type Dao interface {
	GetUserProfileByUserID(userID string) (*model.UserProfile, error)
	UpsertUserProfile(userID string, updateMap map[string]interface{}) error
}
