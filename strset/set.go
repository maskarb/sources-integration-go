package strset

type Set map[string]bool

func New(strs ...string) Set {
	s := Set{}
	for _, str := range strs {
		s[str] = true
	}
	return s
}

func (s Set) Contains(value string) bool {
	_, c := s[value]
	return c
}
