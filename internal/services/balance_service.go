package services

type BalanceService interface {
}

type balanceService struct {
	userService UserService
}

func NewBalance(userService UserService) AuthService {
	return &authService{userService: userService}
}
