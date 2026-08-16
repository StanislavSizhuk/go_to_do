package users_postgres_repository

import "github.com/StanislavSizhuk/go_to_do/internal/core/domain"

type UserModel struct {
	ID          int
	FullName    string
	PhoneNumber *string
	Version     int
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, userModel := range users {
		userDomains[i] = domain.NewUser(
			userModel.ID,
			userModel.Version,
			userModel.FullName,
			userModel.PhoneNumber,
		)
	}

	return userDomains
}
