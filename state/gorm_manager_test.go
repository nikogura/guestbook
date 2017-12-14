package state

import (
	"fmt"
	"github.com/nikogura/guestbook/config"
	"github.com/stitchfix/go-postgres-testdb/testdb"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var tmpDir string
var dbName string
var dbPid int
var configObj config.Config
var configFileName string
var logger *log.Logger

func TestMain(m *testing.M) {
	setUp()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setUp() {
	tempDir, err := ioutil.TempDir("", "guestbook")
	if err != nil {
		fmt.Printf("Error creating temp dir %q: %s", tempDir, err)
		os.Exit(1)
	}

	// legerdemain so we populate the global variable
	tmpDir = tempDir

	// DB setup
	dbName = "guestbook"
	dbDir := fmt.Sprintf("%s/%s", tmpDir, dbName)

	err = os.Mkdir(dbDir, 0700)
	if err != nil {
		fmt.Printf("Error creating db dir: %s", err)
	}

	fmt.Printf("Starting Database Server\n")

	pid, port, err := testdb.StartTestDB(dbDir, dbName)
	if err != nil {
		fmt.Printf("Failed to start test db %q: %s\n", dbName, err)
	}

	dbPid = pid

	fmt.Printf("Db Started with pid %d on port %d\n", dbPid, port)

	fmt.Println("DB Setup Complete.")

	configFileName = fmt.Sprintf("%s/%s", tempDir, config.TestConfigFileName())

	err = ioutil.WriteFile(configFileName, []byte(config.TestConfigFileContents(port)), 0644)

	if err != nil {
		fmt.Printf("Error writing config file %q: %s", configFileName, err)
		os.Exit(1)
	}

	configObj, err = config.ReadConfig(configFileName)
	if err != nil {
		fmt.Printf("Error reading config file %q: %s", configFileName, err)
		os.Exit(1)
	}

	logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	log.Printf("Setup complete.\n")

}

func tearDown() {
	// DB Teardown
	if dbPid != 0 {
		err := testdb.StopPostgres(dbPid)
		if err != nil {
			fmt.Printf("Failed to stop postgres process %d", dbPid)
		} else {
			fmt.Printf("database pid %d shut down.\n", dbPid)
		}
	}

	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		os.RemoveAll(tmpDir)
	}

}

func TestGORMStateManager_CRUD(t *testing.T) {
	gm, err := NewGORMManager(configObj, logger)
	if err != nil {
		log.Printf("Failed to initialize gorm manager: %s", err)
		t.Fail()
	}

	// assert the visitor is not already in the db
	visitor, err := gm.GetVisitor(testVisitorIp())
	if err != nil {
		log.Printf("Error geting visitor: %s", err)
		t.Fail()
	}

	assert.Equal(t, Visitor{}, visitor, "Nonexistant visitor truly does not exist")

	// add a visitor

	visitor, err = gm.NewVisitor(testVisitorObj())
	if err != nil {
		log.Printf("Failed to add visitor")
		t.Fail()
	}

	// verify the visitor is in the db
	visitor, err = gm.GetVisitor(testVisitorIp())
	if err != nil {
		log.Printf("Failed to retrieve visitor: %s", err)
		t.Fail()
	}

	assert.Equal(t, testVisitorObj(), visitor, "Retrieved visitor object matches expectations")

	// remove the visitor

	err = gm.RemoveVisitor(visitor)
	if err != nil {
		log.Printf("Failed to remove visitor: %s", err)
		t.Fail()
	}

	// verify the visitor is truly gone
	visitor, err = gm.GetVisitor(testVisitorIp())
	if err != nil {
		log.Printf("Error geting visitor: %s", err)
		t.Fail()
	}

	assert.Equal(t, Visitor{}, visitor, "Nonexistant visitor truly does not exist")

}
