package set

func New() Strings {
	m := make(map[string]struct{})
	return m
}

type Strings map[string]struct{}

func (set Strings) Add(s string) {
	set[s] = struct{}{}
}

func (set Strings) Have(s string) bool {
	_, ok := set[s]
	return ok
}
