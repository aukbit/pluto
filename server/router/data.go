package router

type Data struct {
	value     string
	prefix    string
	isDynamic bool
	vars      []string
	methods   map[string]Handler
}

// NewData returns a new data instance
func NewData() *Data {
	return &Data{
		value:   "",
		prefix:  "",
		vars:    []string{},
		methods: make(map[string]Handler)}
}

func (d *Data) GetValue() string {
	return d.value
}

func (d *Data) SetValue(val string) {
	d.value = val
}

func (d *Data) IsDynamic() bool {
	return d.isDynamic
}

func (d *Data) SetIsDynamic(dyn bool) {
	d.isDynamic = dyn
}

func (d *Data) GetVars() []string {
	return d.vars
}

func (d *Data) AddVar(v string) {
	d.vars = append(d.vars, v)
}

func (d *Data) Get() Handler {
	if h, ok := d.methods["GET"]; ok {
		return h
	}
	return nil
}
