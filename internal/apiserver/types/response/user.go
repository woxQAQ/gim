package response

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Status   string `json:"status"`
}
