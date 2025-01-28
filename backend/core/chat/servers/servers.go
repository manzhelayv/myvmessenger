package servers

//import "context"

type Servers interface {
	Start() error
	Stop() error
}
