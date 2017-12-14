package config

import (
	"encoding/json"
	"fmt"
)

func TestConfigFileName() string {
	return "guestbook.json"
}

func TestConfigFileContents(dbPort int) string {
	return fmt.Sprintf(`{
  "state": {
    "manager": {
      "type": "gorm",
	  	"dialect": "postgres",
    	"connect_string": "postgresql://guestbook:guestbook@localhost:%d/guestbook?sslmode=disable"
    }
  },
  "server": {
    "addr": "localhost:8080",
    "read_timeout": 5,
    "write_timeout": 10
  },
  "monitor": {
    "poll_interval": 30,
    "deployment_timeout": 30
  }
}`, dbPort)

}

func TestDefaultConfig() *config {
	c := config{
		d: make(map[string]interface{}),
	}

	state := make(map[string]interface{})
	manager := make(map[string]interface{})
	manager["type"] = TestDefaultManagerTypeName()
	manager["dialect"] = TestDefaultManagerDialect()
	manager["connect_string"] = TestDefaultManagerConnectString()

	state["manager"] = manager

	server := make(map[string]interface{})
	server["addr"] = TestDefaultServerAddr()
	server["read_timeout"] = TestDefaultServerReadTimeout()
	server["write_timeout"] = TestDefaultServerWriteTimeout()

	monitor := make(map[string]interface{})
	monitor["poll_interval"] = TestDefaultPollInterval()
	monitor["deployment_timeout"] = TestDefaultDeploymentTimeout()

	c.d["state"] = state
	c.d["server"] = server
	c.d["monitor"] = monitor

	return &c
}

func TestDefaultPollInterval() json.Number {
	return json.Number("30")
}

func TestDefaultDeploymentTimeout() json.Number {
	return json.Number("30")
}

func TestDefaultServerAddr() string {
	return "localhost:8080"
}

func TestDefaultServerReadTimeout() json.Number {
	return json.Number("5")
}

func TestDefaultServerWriteTimeout() json.Number {
	return json.Number("10")
}

func TestDefaultManagerTypeName() string {
	return "gorm"
}

func TestDefaultManagerConnectString() string {
	return "postgresql://guestbook:guestbook@localhost:5432/guestbook?sslmode=disable"
}

func TestDefaultManagerDialect() string {
	return "postgres"
}
