package cfg

import (
	"os"
	"strconv"
)

type Config interface {
	StringWithDefaults(key string, defaultValue string) string
	IntWithDefaults(key string, defaultValue int) int
}

var EnvCfg Config = &envCfg{}

type envCfg struct {

}

func (*envCfg) StringWithDefaults(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok{
		return val
	}
	return defaultValue
}

func (*envCfg) IntWithDefaults(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok{
		result, err := strconv.Atoi(val)
		if err != nil{
			return defaultValue
		}
		return result
	}
	return defaultValue
}

func MapCfg(m map[string]string) Config{
	if m == nil{
		m = make(map[string]string)
	}
	return &mapCfg{
		m:	m,
	}
}

type mapCfg struct {
	m		map[string]string
}

func (c *mapCfg) StringWithDefaults(key string, defaultValue string) string {
	if val, ok := c.m[key]; ok{
		return val
	}
	return defaultValue
}

func (c *mapCfg) IntWithDefaults(key string, defaultValue int) int {
	if val, ok := c.m[key]; ok{
		result, err := strconv.Atoi(val)
		if err != nil{
			return defaultValue
		}
		return result
	}
	return defaultValue
}
