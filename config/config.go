package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config interface {
	Get(path string) (interface{}, bool)
	GetString(path string, defaultVal string) string
	GetInt(path string, defaultVal int) int
}

type config struct {
	d map[string]interface{}
}

func ReadConfig(configPathOrData string) (Config, error) {
	c := &config{}

	decoder := json.NewDecoder(strings.NewReader(configPathOrData))
	decoder.UseNumber()
	err := decoder.Decode(&c.d)
	if err != nil {
		data, err := os.Open(configPathOrData)
		if err != nil {
			return nil, err
		}

		decoder = json.NewDecoder(data)
		decoder.UseNumber()
		err = decoder.Decode(&c.d)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func getFromMap(pathSpec []string, lvl int, d map[string]interface{}) (interface{}, bool) {
	key := pathSpec[lvl]
	val, ok := d[key]
	if ok && lvl < len(pathSpec)-1 {
		return getFromMap(pathSpec, lvl+1, val.(map[string]interface{}))
	}
	return val, ok
}

func (c *config) Get(path string) (interface{}, bool) {
	pathSpec := strings.Split(path, ".")
	return getFromMap(pathSpec, 0, c.d)
}

func (c *config) GetString(path string, defaultVal string) string {
	if val, ok := c.Get(path); ok {
		return val.(string)
	} else {
		return defaultVal
	}
}

func (c *config) GetInt(path string, defaultVal int) int {
	if val, ok := c.Get(path); ok {
		if num, ok := val.(json.Number); ok {
			if num, err := num.Int64(); err == nil {
				return int(num)
			}
		}
	}
	return defaultVal
}
