package middlewares

import "github.com/apm-dev/vending-machine/domain"

type UserMiddleware struct {
	us domain.UserService
}

func InitUserMiddleware(us domain.UserService) *UserMiddleware {
	return &UserMiddleware{us}
}

