package tasks_postgres_repository

import core_postgres_pool "github.com/StanislavSizhuk/go_to_do/internal/core/repository/posgres/pool"

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
