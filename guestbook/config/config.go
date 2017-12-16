package config

import (
	"encoding/json"
	"os"
	"strings"
)

// Config  configuration object
type Config interface {
	Get(path string) (interface{}, bool)
	GetString(path string, defaultVal string) string
	GetInt(path string, defaultVal int) int
}

type config struct {
	d map[string]interface{}
}

// ReadConfig produces a config object from a file name or a raw string
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

// Get gets a config value from the config object
// values are specified hierarchically by a string path of the form 'foo.bar.baz' where each value between the dots is a level in the config json
func (c *config) Get(path string) (interface{}, bool) {
	pathSpec := strings.Split(path, ".")
	return getFromMap(pathSpec, 0, c.d)
}

// GetString gets a string from the config object
// values are specified hierarchically by a string path of the form 'foo.bar.baz' where each value between the dots is a level in the config json
func (c *config) GetString(path string, defaultVal string) string {
	if val, ok := c.Get(path); ok {
		return val.(string)
	}

	return defaultVal
}

// GetInt gets an integer from the config object
// values are specified hierarchically by a string path of the form 'foo.bar.baz' where each value between the dots is a level in the config json
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
