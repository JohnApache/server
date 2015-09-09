package server

type Shutdown struct {
	funcs func()
}

func NewShutdown(funcs func()) *Shutdown {
	return &Shutdown{
		funcs: funcs,
	}
}
func (sh *Shutdown) Now(int, *int) error {
	sh.funcs()
	return nil
}
