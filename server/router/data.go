package router

// Data struct
type Data struct {
	value   string
	prefix  string
	vars    []string
	methods map[string]Handler
}

// NewData returns a new data instance
func NewData() *Data {
	return &Data{
		value:   "",
		prefix:  "",
		vars:    []string{},
		methods: make(map[string]Handler)}
}

// GetValue returns data value
func (d *Data) GetValue() string {
	return d.value
}

// SetValue sets data value
func (d *Data) SetValue(val string) {
	d.value = val
}
