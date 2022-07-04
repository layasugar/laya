package data

// DataContext interface the extend of contextx.Context.Value(key string) (value interface{})
type DataContext interface {
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
}

var (
	_ DataContext = &MemoryContext{}
)

// MemoryContext gstore data tor memory
type MemoryContext struct {
	data map[string]interface{}
}

// NewMemoryContext new one
func NewMemoryContext() *MemoryContext {
	return &MemoryContext{
		data: map[string]interface{}{},
	}
}

// Set is used to gstore a new key/value pair exclusively for this contextx.
// It also lazy initializes  c.data if it was not used previously.
func (c *MemoryContext) Set(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *MemoryContext) Get(key string) (value interface{}, exists bool) {
	value, exists = c.data[key]
	return
}

// GetString returns the value associated with the key as a string.
func (c *MemoryContext) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *MemoryContext) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *MemoryContext) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *MemoryContext) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}
