package cli

import "context"

type Ctx interface {
	context.Context
	// Pos returns the position arguments
	Pos() []string
	Flag(string) Flag
	Bool(string) bool
	Int(string) int
	Uint(string) uint
	Float(string) float64
	String(string) string
	StringSlice(string) []string
	IntSlice(string) []int
	UintSlice(string) []uint
	FloatSlice(string) []float64
	StringMap(string) map[string]string
	Map(string) map[string]interface{}
}

type ctx struct {
	context.Context
	cmd *Cmd
	pos []string
}

func (c *ctx) Pos() []string {
	return c.pos
}

func (c *ctx) Flag(key string) Flag {
	for _, v := range c.cmd.Flags {
		if key == v.Key() {
			return v
		}
	}
	return nil
}

func Get[T any](c Ctx, key string, def T) T {
	flag := c.Flag(key)
	if flag == nil {
		return def
	}
	value, ok := flag.Var().(T)
	if !ok {
		return def
	}
	return value
}

func (c *ctx) Bool(key string) bool {
	return Get(c, key, false)
}

func (c *ctx) Int(key string) int {
	return Get(c, key, 0)
}

func (c *ctx) Uint(key string) uint {
	return Get(c, key, uint(0))
}

func (c *ctx) Float(key string) float64 {
	return Get(c, key, .0)
}

func (c *ctx) String(key string) string {
	return Get(c, key, "")
}

func (c *ctx) StringSlice(key string) []string {
	return Get[[]string](c, key, nil)
}

func (c *ctx) IntSlice(key string) []int {
	return Get[[]int](c, key, nil)
}

func (c *ctx) UintSlice(key string) []uint {
	return Get[[]uint](c, key, nil)
}

func (c *ctx) FloatSlice(key string) []float64 {
	return Get[[]float64](c, key, nil)
}

func (c *ctx) StringMap(key string) map[string]string {
	return Get(c, key, map[string]string{})
}

func (c *ctx) Map(key string) map[string]any {
	return Get(c, key, map[string]any{})
}
