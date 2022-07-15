package cli

import (
	"strconv"
	"strings"
)

type StringSlice []string

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *StringSlice) String() string {
	return "StringSlice"
}

type IntSlice []int

func (s *IntSlice) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*s = append(*s, i)
	return nil
}

func (s *IntSlice) String() string {
	return "IntSlice"
}

type UintSlice []uint

func (s *UintSlice) Set(value string) error {
	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	*s = append(*s, uint(i))
	return nil
}

func (s *UintSlice) String() string {
	return "UintSlice"
}

type FloatSlice []float64

func (s *FloatSlice) Set(value string) error {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	*s = append(*s, f)
	return nil
}

func (s *FloatSlice) String() string {
	return "FloatSlice"
}

type StringMap map[string]string

func (s *StringMap) Set(value string) error {
	kv := strings.Split(value, "=")
	(*s)[kv[0]] = kv[1]
	return nil
}

func (s *StringMap) String() string {
	return "StringMap"
}

type Map map[string]interface{}

func (s *Map) Set(value string) error {
	kv := strings.Split(value, "=")
	(*s)[kv[0]] = kv[1]
	return nil
}

func (s *Map) String() string {
	return "Map"
}
