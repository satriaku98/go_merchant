package manager

import "go_merchant/usecase"

type UseCaseManager interface {
	LoginUseCase() usecase.LoginUseCase
}

type useCaseManager struct {
	repo RepoManager
}

func (u *useCaseManager) LoginUseCase() usecase.LoginUseCase {
	return usecase.NewLoginUsecase(u.repo.LoginRepo())
}

func NewUseCaseManager(repo RepoManager) UseCaseManager {
	return &useCaseManager{
		repo,
	}
}
