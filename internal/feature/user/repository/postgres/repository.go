package users_postgres_repository

import core_postgres_pool "github.com/StanislavSizhuk/go_to_do/internal/core/repository/posgres/pool"

type UsersRepository struct {
	pool core_postgres_pool.Pool
	// Add any necessary fields for the repository, such as a database connection or configuration.
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
