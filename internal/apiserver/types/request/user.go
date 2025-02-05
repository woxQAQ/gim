package request

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GetUserInfoRequest struct {
	ID int64 `query:"userId" binding:"required"`
}
