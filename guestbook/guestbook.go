package guestbook

import (
	"github.com/nikogura/guestbook/guestbook/config"
	"github.com/nikogura/guestbook/guestbook/service"
	"github.com/nikogura/guestbook/guestbook/state"
	"log"
	"os"
)

// ConfigFileName default config file name
const ConfigFileName = "/etc/guestbook/guestbook.json"

// DefaultDBPort default database port
const DefaultDBPort = 5432

var logger *log.Logger

// Run runs the guestbook server
func Run() {
	var configObj config.Config

	if _, err := os.Stat(ConfigFileName); !os.IsNotExist(err) {
		configObj, err = config.ReadConfig(ConfigFileName)
		if err != nil {
			log.Printf("Error reading config file %q: %s", ConfigFileName, err)
			os.Exit(1)
		}

	} else {
		configObj, err = config.ReadConfig(config.TestConfigFileContents(DefaultDBPort))
		if err != nil {
			log.Printf("Failed to create default config: %s", err)
			os.Exit(1)
		}
	}

	logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	manager, err := state.NewGORMManager(configObj, logger)
	if err != nil {
		log.Printf("Failed to instantiate state manager: %s", err)
		os.Exit(1)
	}

	err = service.Run(configObj.GetString("server.addr", ""), &manager)

	if err != nil {
		log.Printf("error running server: %s", err)
		os.Exit(1)
	}
}
