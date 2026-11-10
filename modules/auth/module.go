package auth

import "github.com/vitor-franz/book-log/modules"

func NewModule(authUseCase UseCase) modules.Module {
	return module{authUseCase: authUseCase}
}

type module struct {
	authUseCase UseCase
}

func (m module) Path() string {
	//TODO implement me
	panic("implement me")
}

func (m module) Setup() {
	//TODO implement me
	panic("implement me")
}
