package rabbitmq

// Header struct to store filtering headers dynamically
type Header struct {
	Fields map[string]interface{}
}

// NewHeader creates a new Header instance
func NewHeader() *Header {
	return &Header{Fields: make(map[string]interface{})}
}

// AddField adds a new key-value pair to the Header
func (h *Header) AddField(key string, value interface{}) {
	h.Fields[key] = value
}
