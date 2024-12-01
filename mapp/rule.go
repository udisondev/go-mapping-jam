package mapp

type Rule struct {
	spec string
}

func (r *Rule) Value() string {
	return r.spec
}