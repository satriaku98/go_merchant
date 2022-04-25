package usecase

import (
	"go_merchant/model"
	"go_merchant/repository"
)

type LoginUseCase interface {
	LoginCustomer(username string, passcode string) (model.Customer, int, error)
	LogoutCustomer(username string, passcode string) (int, error)
}

type loginUseCase struct {
	repo repository.CustomerRepo
}

func (c *loginUseCase) LoginCustomer(username string, passcode string) (model.Customer, int, error) {
	return c.repo.LoginCustomer(username, passcode)
}

func (c *loginUseCase) LogoutCustomer(username string, passcode string) (int, error) {
	return c.repo.LogoutCustomer(username, passcode)
}

func NewLoginUsecase(repo repository.CustomerRepo) LoginUseCase {
	return &loginUseCase{repo}
}
