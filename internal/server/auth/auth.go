package auth

type AuthenticateReq struct {
	UserName string
	Token    string
}

func IsAuthenticated(req AuthenticateReq) bool {
	// todo 访问数据库
	return true
}
