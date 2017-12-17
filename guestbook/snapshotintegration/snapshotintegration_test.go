package snapshotintegration

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nikogura/guestbook/guestbook/snapshot"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

// AWS Credentials, and an instance with name 'testbox' is required to run these tests.
// set ActuallyDoSnapshots = true if you actually want to snapshot the box
// You can confirm in the console, which is crude, and inelegant, but I was running up against resource limits, and couldn't spend to fully flesh out the test

var actuallyDoSnapshots bool
var awsSession *session.Session

// TestMain runs before each test case
func TestMain(m *testing.M) {
	setUp()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setUp() {
	actuallyDoSnapshots = false
	sess, err := snapshot.Ec2Session()
	if err != nil {
		log.Printf("AWS Session could not be created.  Skipping integration tests.")
		os.Exit(0)
	}

	awsSession = sess

	// In a perfect world, we'd create the test instance we're working with.
}

func tearDown() {

	// and in that perfect world, we'd remove it when we're done.

}

// TestGetInstanceInfoMaps gets instance information filtered to the test box name in fixtures and verifies all is well
func TestGetInstanceInfoMaps(t *testing.T) {
	//infomaps, err := snapshot.GetInstanceInfoMaps(awsSession, []string{testInstanceName()})
	infomaps, err := snapshot.GetInstanceInfoMaps(awsSession, nil)
	if err != nil {
		log.Printf("Error fetching instances: %s", err)
		t.Fail()
	}

	for id, instanceInfo := range infomaps.Id2Info {
		fmt.Printf("Instance:\n")
		fmt.Printf("  Instance ID: %s  Name: %s\n", id, instanceInfo.InstanceName)
	}

	testInfo := infomaps.Name2Info[testInstanceName()]

	assert.NotEqual(t, snapshot.InstanceInfo{}, testInfo, "Returned info is not the zero value for the struct")

	// volumes need id's not names
	ids := make([]string, 0)

	ids = append(ids, infomaps.Name2Info[testInstanceName()].InstanceId)

	volInfo, err := snapshot.GetVolumeInfo(awsSession, ids)
	if err != nil {
		log.Printf("Error getting volumes: %s", err)
		t.Fail()
	}

	found := false

	for _, vol := range volInfo {
		fmt.Printf("Volume:\n")
		fmt.Printf("  ID: %s, Instance: %s, Device: %s\n", vol.VolumeId, vol.InstanceId, vol.DeviceName)
		if vol.InstanceId == testInfo.InstanceId {
			found = true
			break
		}
	}

	assert.True(t, found, "Successfully found volume for test box")

}

// TestSnapshotRunningVolumes actually runs the snapshots
// This can stack up the snapshots quickly
// Ideally we would run this, and verify the snapshot was made, with the right tags, id's, etc
// But this could seriously clog up my free tier usage.
// I think I've demonstrated that it can be done.
func TestSnapshotRunningVolumes(t *testing.T) {
	if actuallyDoSnapshots {
		// ideally verify that we don't already have a snapshot and therefore confuse ourselves

		// make the snapshot
		err := snapshot.SnapshotRunningVolumes(awsSession, []string{testInstanceName()})
		if err != nil {
			log.Printf("error snapshotting volumes: %s", err)
			t.Fail()
		}

		// verify the snapshot is there

		// read the tags on the snapshot
		// verify Name and Date tags exist and contain expected values

		// remove the snapshot
	}
}
