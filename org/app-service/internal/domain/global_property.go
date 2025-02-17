package domain

type GlobalProperty struct {
	ModelBase
	Key         string
	Value       string
	Description string
	ValueType   string
}
