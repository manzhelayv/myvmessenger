package servers

type Servers interface {
	Start() error
	Stop() error
}
