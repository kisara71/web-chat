package model

type UserProfile struct {
	UserID          string `gorm:"primaryKey;column:user_id;type:varchar(36)"`
	SystemPrompt    string `gorm:"column:system_prompt;type:text"`
	ModelPreference string `gorm:"column:model_preference;type:text"`
	Traits          string `gorm:"column:traits;type:text"`
	CommonPartNoUnique
}

func (UserProfile) TableName() string { return "user_profile" }
