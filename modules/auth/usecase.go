package auth

func NewUseCase(authRepository Repository) UseCase {
	return useCase{authRepository: authRepository}
}

type useCase struct {
	authRepository Repository
}
