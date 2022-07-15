package cli

type Flag interface {
	Key() string
	Help() string
	Invalid() bool
	Var() interface{}
}

// =========================================================== //
// I've tried to use Go generics but it made my code too ugly. //
// =========================================================== //

type BoolFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    bool
}

func (o BoolFlag) Key() string      { return o.Name }
func (o BoolFlag) Help() string     { return o.Usage }
func (o BoolFlag) Var() interface{} { return o.Value }
func (o BoolFlag) Invalid() bool    { return false }

type IntFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    int
}

func (o IntFlag) Key() string      { return o.Name }
func (o IntFlag) Help() string     { return o.Usage }
func (o IntFlag) Var() interface{} { return o.Value }
func (o IntFlag) Invalid() bool    { return false }

type UintFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    uint
}

func (o UintFlag) Key() string      { return o.Name }
func (o UintFlag) Help() string     { return o.Usage }
func (o UintFlag) Var() interface{} { return o.Value }
func (o UintFlag) Invalid() bool    { return false }

type FloatFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    float64
}

func (o FloatFlag) Key() string      { return o.Name }
func (o FloatFlag) Help() string     { return o.Usage }
func (o FloatFlag) Var() interface{} { return o.Value }
func (o FloatFlag) Invalid() bool    { return false }

type StringFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    string
}

func (o StringFlag) Key() string      { return o.Name }
func (o StringFlag) Help() string     { return o.Usage }
func (o StringFlag) Var() interface{} { return o.Value }
func (o StringFlag) Invalid() bool    { return o.Required && len(o.Value) == 0 }

type StringSliceFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    []string
}

func (o StringSliceFlag) Key() string      { return o.Name }
func (o StringSliceFlag) Help() string     { return o.Usage }
func (o StringSliceFlag) Var() interface{} { return o.Value }
func (o StringSliceFlag) Invalid() bool    { return o.Required && o.Value == nil }

type IntSliceFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    []int
}

func (o IntSliceFlag) Key() string      { return o.Name }
func (o IntSliceFlag) Help() string     { return o.Usage }
func (o IntSliceFlag) Var() interface{} { return o.Value }
func (o IntSliceFlag) Invalid() bool    { return o.Required && o.Value == nil }

type UintSliceFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    []uint
}

func (o UintSliceFlag) Key() string      { return o.Name }
func (o UintSliceFlag) Help() string     { return o.Usage }
func (o UintSliceFlag) Var() interface{} { return o.Value }
func (o UintSliceFlag) Invalid() bool    { return o.Required && o.Value == nil }

type FloatSliceFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    []float64
}

func (o FloatSliceFlag) Key() string      { return o.Name }
func (o FloatSliceFlag) Help() string     { return o.Usage }
func (o FloatSliceFlag) Var() interface{} { return o.Value }
func (o FloatSliceFlag) Invalid() bool    { return o.Required && o.Value == nil }

type StringMapFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    map[string]string
}

func (o StringMapFlag) Key() string      { return o.Name }
func (o StringMapFlag) Help() string     { return o.Usage }
func (o StringMapFlag) Var() interface{} { return o.Value }
func (o StringMapFlag) Invalid() bool    { return o.Required && o.Value == nil }

type MapFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    map[string]interface{}
}

func (o MapFlag) Key() string      { return o.Name }
func (o MapFlag) Help() string     { return o.Usage }
func (o MapFlag) Var() interface{} { return o.Value }
func (o MapFlag) Invalid() bool    { return o.Required && o.Value == nil }
