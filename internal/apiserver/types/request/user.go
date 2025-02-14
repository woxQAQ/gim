package request

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password" validate:"required"`
}

type GetUserInfoRequest struct {
	ID int64 `query:"userId" validate:"required"`
}
