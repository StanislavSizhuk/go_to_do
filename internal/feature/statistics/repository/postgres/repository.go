package statistics_postgres_repository

import core_postgres_pool "github.com/StanislavSizhuk/go_to_do/internal/core/repository/posgres/pool"

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(
	pool core_postgres_pool.Pool,
) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
