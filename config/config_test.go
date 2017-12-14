package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

var tempDir string
var configFileName string
var configObj Config

func setUp() {
	tempDir, err := ioutil.TempDir("", "eve")
	if err != nil {
		fmt.Printf("Error creating temp dir %q: %s", tempDir, err)
		os.Exit(1)
	}

	configFileName = fmt.Sprintf("%s/%s", tempDir, TestConfigFileName())

	err = ioutil.WriteFile(configFileName, []byte(TestConfigFileContents(5432)), 0644)

	if err != nil {
		fmt.Printf("Error writing config file %q: %s", configFileName, err)
		os.Exit(1)
	}

	configObj, err = ReadConfig(configFileName)
	if err != nil {
		fmt.Printf("Error reading config file %q: %s", configFileName, err)
		os.Exit(1)
	}
}

func tearDown() {
	os.RemoveAll(tempDir)
}

func TestMain(m *testing.M) {
	setUp()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func TestReadConfig(t *testing.T) {
	testConfig := TestDefaultConfig()

	assert.Equal(t, testConfig, configObj, "Config object read from disk meets expecatations.")

}

func TestReadConfigFromString(t *testing.T) {
	testConfig := TestDefaultConfig()

	configObjFromString, err := ReadConfig(TestConfigFileContents(5432))
	if err != nil {
		fmt.Printf("Error creating config object from string: %s", err)
		t.Fail()
	}

	assert.Equal(t, testConfig, configObjFromString, "Config object read from disk meets expecatations.")
}

func TestConfig_Get(t *testing.T) {
	expectedDialect := TestDefaultManagerDialect()

	actualDialect, err := configObj.Get("state.manager.dialect")
	if err != true {
		fmt.Printf("Error retrieving structure from config object.")
		t.Fail()
	}

	assert.Equal(t, expectedDialect, actualDialect, "Bom Engine read from config object meets expectations.")

}

func TestConfig_GetString(t *testing.T) {
	expectedConnectString := TestDefaultManagerConnectString()

	actualConnectString := configObj.GetString("state.manager.connect_string", "")

	assert.Equal(t, expectedConnectString, actualConnectString, "String fetched from config object matches expectations.")

}

func TestConfig_GetInt(t *testing.T) {
	expectedTimeout, err := strconv.Atoi(string(TestDefaultServerReadTimeout()))
	if err != nil {
		fmt.Printf("Error converting default memory %q to integer", err)
		t.Fail()
	}

	actualTimeout := configObj.GetInt("server.read_timeout", 0)

	assert.Equal(t, expectedTimeout, actualTimeout, "Integer fetched from GetInt matches expectations.")

}
