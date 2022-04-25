package manager

import (
	"github.com/jmoiron/sqlx"
	"go_merchant/repository"
)

type RepoManager interface {
	LoginRepo() repository.CustomerRepo
}

type repoManager struct {
	SqlxDb *sqlx.DB
}

func (r *repoManager) LoginRepo() repository.CustomerRepo {
	return repository.NewLoginRepo(r.SqlxDb)
}

func NewRepoManager(sqlxDb *sqlx.DB) RepoManager {
	return &repoManager{
		sqlxDb,
	}
}
