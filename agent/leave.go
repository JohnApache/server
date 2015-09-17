package agent

type Leave struct {
	funcs func(*Session)
}

func NewLeave(funcs func(*Session)) *Leave {
	return &Leave{
		funcs: funcs,
	}
}

func (le *Leave) User(s *Session, i *int) error {
	le.funcs(s)
	return nil
}
