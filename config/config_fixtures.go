package config

import (
	"encoding/json"
	"fmt"
)

// TestConfigFileName  name of the standard config file
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
		"port": 8080
  }
}`, dbPort)

}

// TestDefaultConfig  a default config file for testing
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
	server["port"] = TestDefaultServerPort()

	c.d["state"] = state
	c.d["server"] = server

	return &c
}

// TestDefaultServerAddr  the default server address
func TestDefaultServerAddr() string {
	return "localhost:8080"
}

// TestDefaultManagerTypeName the default state manager type
func TestDefaultManagerTypeName() string {
	return "gorm"
}

// TestDefaultmanagerConnectString default connection string
func TestDefaultManagerConnectString() string {
	return "postgresql://guestbook:guestbook@localhost:5432/guestbook?sslmode=disable"
}

// TestDefaultManagerDialect default dialect of db
func TestDefaultManagerDialect() string {
	return "postgres"
}

func TestDefaultServerPort() json.Number {
	return json.Number("8080")
}
