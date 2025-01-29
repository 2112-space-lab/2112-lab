package clients

// IClient definition
type IClient interface {
	Name() string
	Configure(v any)
}
