package profile

import (
	"errors"
	"web-chat/internal/model"

	"gorm.io/gorm"
)

type profileDaoImpl struct {
	db *gorm.DB
}

func NewDao(db *gorm.DB) Dao {
	return &profileDaoImpl{db: db}
}

func (p *profileDaoImpl) GetUserProfileByUserID(userID string) (*model.UserProfile, error) {
	var entity model.UserProfile
	if err := p.db.Where("user_id = ?", userID).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (p *profileDaoImpl) UpsertUserProfile(userID string, updateMap map[string]interface{}) error {
	var entity model.UserProfile
	err := p.db.Where("user_id = ?", userID).First(&entity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		entity = model.UserProfile{UserID: userID}
		if val, ok := updateMap["system_prompt"]; ok {
			entity.SystemPrompt, _ = val.(string)
		}
		if val, ok := updateMap["model_preference"]; ok {
			entity.ModelPreference, _ = val.(string)
		}
		if val, ok := updateMap["traits"]; ok {
			entity.Traits, _ = val.(string)
		}
		return p.db.Create(&entity).Error
	}
	return p.db.Model(&model.UserProfile{}).Where("user_id = ?", userID).Updates(updateMap).Error
}
