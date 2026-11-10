package modules

type Module interface {
	Path() string
	Setup()
}
