package statemachine

// CompleteOptions .
type CompleteOptions struct {
	options map[string]interface{}
}

// NewCompleteOptions .
func NewCompleteOptions() CompleteOptions {
	return CompleteOptions{
		options: make(map[string]interface{}),
	}
}

// SetString .
func (c *CompleteOptions) SetString(name string, value string) {
	c.options[name] = value
}

// SetInt .
func (c *CompleteOptions) SetInt(name string, value int64) {
	c.options[name] = value
}

// SetUInt32 .
func (c *CompleteOptions) SetUInt32(name string, value uint32) {
	c.options[name] = value
}

// SetFloat .
func (c *CompleteOptions) SetFloat(name string, value float32) {
	c.options[name] = value
}

// SetBool .
func (c *CompleteOptions) SetBool(name string, value bool) {
	c.options[name] = value
}

// SetObject .
func (c *CompleteOptions) SetObject(name string, value interface{}) {
	c.options[name] = value
}

// GetString .
func (c *CompleteOptions) GetString(name string) (string, bool) {
	return getOption[string](name, c.options)
}

// GetInt .
func (c *CompleteOptions) GetInt(name string) (int64, bool) {
	return getOption[int64](name, c.options)
}

// GetFloat .
func (c *CompleteOptions) GetFloat(name string) (float32, bool) {
	return getOption[float32](name, c.options)
}

// GetBool .
func (c *CompleteOptions) GetBool(name string) (bool, bool) {
	return getOption[bool](name, c.options)
}

// GetObject .
func (c *CompleteOptions) GetObject(name string) (interface{}, bool) {
	return getOption[interface{}](name, c.options)
}

func getOption[T any](name string, options map[string]interface{}) (T, bool) {
	var result T
	v, find := options[name]
	if !find {
		return result, false
	}
	r, ok := v.(T)
	if !ok {
		return result, false
	}
	return r, true
}
