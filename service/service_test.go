package service

import (
	"fmt"
	"github.com/nikogura/go-postgres-testdb/testdb"
	"github.com/nikogura/guestbook/config"
	"github.com/nikogura/guestbook/state"
	"github.com/phayes/freeport"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"
)

var dbPort int
var servicePort int
var serverAddress string
var tmpDir string
var dbName string
var dbPid int
var logger *log.Logger
var configObj config.Config

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
	dir, err := ioutil.TempDir("", "testdb")
	if err != nil {
		fmt.Printf("Error creating temp dir %q: %s", tempDir, err)
		os.Exit(1)
	}

	tempDir = dir

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
	dbPort = port

	fmt.Printf("Db Started with pid %d on port %d\n", dbPid, port)

	fmt.Println("DB Setup Complete.")

	// service setup
	log.Printf("Attempting to find a free port on which to run the service.\n")
	randPort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatalf("Failed to get a free port: %s\n", err)
		os.Exit(1)
	}

	servicePort = randPort
	serverAddress = fmt.Sprintf("localhost:%d", servicePort)

	log.Printf("Running service on port %d", servicePort)

	log.Printf("Running http server on %s", serverAddress)

	logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	configObj, err = config.ReadConfig(config.TestConfigFileContents(port))
	if err != nil {
		log.Printf("Failed to create default config: %s", err)
		os.Exit(1)
	}

	manager, err := state.NewGORMManager(configObj, logger)
	if err != nil {
		log.Printf("Failed to instantiate state manager: %s", err)
		os.Exit(1)
	}

	go Run(serverAddress, &manager)

	log.Printf("Server is running.  Sleeping 3 seconds to let it get it's bearings before we hit it.\n")

	time.Sleep(time.Second * 3)

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

// these are really inelegant tests, but we're not showcasing my web app testing ability
func TestService(t *testing.T) {

	username := "fargle"

	// make sure we're empty
	uri := fmt.Sprintf("http://%s/guestbook", serverAddress)
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Printf("error hitting http server")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile(`.*Howdy stranger.*`)

	if resp.StatusCode != 200 {
		t.Fail()
	}

	if !re.MatchString(string(body)) {
		t.Fail()
	}

	// sign the book
	uri = fmt.Sprintf("http://%s/guestbook/sign", serverAddress)

	resp, err = http.PostForm(uri, url.Values{"visitor": {username}})
	if err != nil {
		fmt.Printf("Error posting to form: %s", err)
		t.Fail()
	}

	if resp.StatusCode != 200 {
		t.Fail()
	}

	// make sure we're recognized
	uri = fmt.Sprintf("http://%s/guestbook", serverAddress)
	resp, err = http.Get(uri)
	if err != nil {
		fmt.Printf("error hitting http server")
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	re = regexp.MustCompile(fmt.Sprintf(`.*Good to see you again %s!.*`, username))

	if resp.StatusCode != 200 {
		t.Fail()
	}

	if !re.MatchString(string(body)) {
		t.Fail()
	}

}
