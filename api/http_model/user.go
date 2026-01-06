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
	Token string `json:"token" binding:"required"`
}
