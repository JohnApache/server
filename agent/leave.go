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
	if le.funcs != nil {
		le.funcs(s)
	}
	sessionLeave(s)
	return nil
}
