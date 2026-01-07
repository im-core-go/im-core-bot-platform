package http_model

type UserRegisterReq struct {
	NickName string  `json:"nick_name" binding:"required"`
	Email    *string `json:"email" binding:"omitempty"`
	Phone    string  `json:"phone" binding:"required"`
	Password string  `json:"password" binding:"required"`
}

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateReq struct {
	NickName *string `json:"nick_name" binding:"omitempty"`
	Email    *string `json:"email" binding:"omitempty"`
	Phone    *string `json:"phone" binding:"omitempty"`
	Password *string `json:"password" binding:"omitempty"`
}

type LogoutReq struct {
	Token string `json:"token" binding:"required"`
}
<<<<<<< Updated upstream
=======

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
>>>>>>> Stashed changes
