package wrapper

type StringSlice struct {
	Value []string
}

func (s *StringSlice) Add(script string) *StringSlice {
	s.Value = append(s.Value, script)
	return s
}
func (s *StringSlice) AddSeveral(script []string) *StringSlice {
	s.Value = append(s.Value, script...)
	return s
}

func (s *StringSlice) Get() []string {
	return s.Value
}

func NewStringSlice() *StringSlice {
	return &StringSlice{
		Value: []string{},
	}
}
