package discovery

type Discovery interface {
	IsAvailable() (bool, error)
	Register()
	Unregister()
}
