package user

import (
	"errors"
	"fmt"
	"strings"
	"web-chat/api/http_model"

	"gorm.io/gorm"
)

func (l *logicImpl) UpdateProfile(userID string, req *http_model.UserProfileUpdateReq) error {
	if userID == "" {
		return fmt.Errorf("user id is required")
	}
	if req == nil {
		return fmt.Errorf("profile request is nil")
	}
	updateMap := make(map[string]interface{})
	if req.SystemPrompt != nil {
		updateMap["system_prompt"] = strings.TrimSpace(*req.SystemPrompt)
	}
	if req.ModelPreference != nil {
		updateMap["model_preference"] = strings.TrimSpace(*req.ModelPreference)
	}
	if req.Traits != nil {
		updateMap["traits"] = strings.TrimSpace(*req.Traits)
	}
	if len(updateMap) == 0 {
		return fmt.Errorf("no fields to update")
	}
	return l.svcCtx.Dao.ProfileDao.UpsertUserProfile(userID, updateMap)
}

func (l *logicImpl) GetProfile(userID string) (*http_model.UserProfileResp, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	profile, err := l.svcCtx.Dao.ProfileDao.GetUserProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &http_model.UserProfileResp{}, nil
		}
		return nil, err
	}
	return &http_model.UserProfileResp{
		SystemPrompt:    profile.SystemPrompt,
		ModelPreference: profile.ModelPreference,
		Traits:          profile.Traits,
		CreatedAt:       profile.CreatedAt,
		UpdatedAt:       profile.UpdatedAt,
	}, nil
}
