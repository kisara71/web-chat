package http_model

type UserRegisterReq struct {
	NickName  string  `json:"nick_name" binding:"required"`
	Email     string  `json:"email" binding:"required"`
	EmailCode string  `json:"email_code" binding:"required"`
	Phone     *string `json:"phone" binding:"omitempty"`
	Password  string  `json:"password" binding:"required"`
}

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthTokenResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginCodeReq struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type SendEmailCodeReq struct {
	Email string `json:"email" binding:"required"`
}

type UserUpdateReq struct {
	NickName *string `json:"nick_name" binding:"omitempty"`
	Email    *string `json:"email" binding:"omitempty"`
	Phone    *string `json:"phone" binding:"omitempty"`
	Password *string `json:"password" binding:"omitempty"`
}

type LogoutReq struct {
	Token        string `json:"token" binding:"omitempty"`
	RefreshToken string `json:"refresh_token" binding:"omitempty"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserInfoReq struct {
	UserID string `form:"user_id" binding:"required"`
}

type UserInfoResp struct {
	UUID      string  `json:"uuid"`
	NickName  string  `json:"nick_name"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone,omitempty"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}

type UserProfileUpdateReq struct {
	SystemPrompt    *string `json:"system_prompt" binding:"omitempty"`
	ModelPreference *string `json:"model_preference" binding:"omitempty"`
	Traits          *string `json:"traits" binding:"omitempty"`
}

type UserProfileResp struct {
	SystemPrompt    string `json:"system_prompt"`
	ModelPreference string `json:"model_preference"`
	Traits          string `json:"traits"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}
