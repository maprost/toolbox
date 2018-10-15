package tbtype

type StringSet struct {
	m map[string]struct{}
}

func NewStringSet(item ...string) StringSet {
	s := StringSet{m: make(map[string]struct{})}
	s.AddList(item)
	return s
}

func (s *StringSet) Add(item string) {
	if len(item) == 0 {
		return
	}

	s.m[item] = struct{}{}
}

func (s *StringSet) AddList(items []string) {
	for _, item := range items {
		s.Add(item)
	}
}

func (s *StringSet) AddSet(items StringSet) {
	for item := range items.m {
		s.Add(item)
	}
}

func (s *StringSet) Contains(item string) bool {
	_, ok := s.m[item]
	return ok
}

func (s *StringSet) Keys() []string {
	result := make([]string, 0, len(s.m))
	for k := range s.m {
		result = append(result, k)
	}
	return result
}

func (s *StringSet) Len() int {
	return len(s.m)
}

func (s *StringSet) IsNil() bool {
	return s.m == nil
}
