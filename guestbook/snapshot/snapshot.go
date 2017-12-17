package snapshot

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"time"
)

// VolInfo contains relevant information about volumes
type VolInfo struct {
	InstanceId string
	DeviceName string
	VolumeId   string
}

// InstanceInfo contains relevant information about Instances
type InstanceInfo struct {
	InstanceId   string
	InstanceName string
}

// InstanceInfoMaps areUseful Maps of instance info.
type InstanceInfoMaps struct {
	Id2Info   map[string]InstanceInfo
	Name2Info map[string]InstanceInfo
}

// Ec2Session returns a new ec2 session with shared configs
func Ec2Session() (awsSession *session.Session, err error) {
	awsSession, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	return awsSession, err
}

// StringInSlice returns true if the given string is in the given slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetVolumeInfo gets relevant info about volumes currently in existence.  Filter is a slice of instance id's to filter on.  Filter is optional.  Without it it returns all volumes
func GetVolumeInfo(awsSession *session.Session, targets []string) (info []VolInfo, err error) {
	client := ec2.New(awsSession)
	info = make([]VolInfo, 0)

	filters := make([]*ec2.Filter, 0)

	params := &ec2.DescribeVolumesInput{}

	// process targets and massage them into aws type variables
	if targets != nil {
		awsnames := make([]*string, 0)

		for _, name := range targets {
			awsnames = append(awsnames, aws.String(name))
		}

		nameFilter := ec2.Filter{
			Name:   aws.String("attachment.instance-id"),
			Values: awsnames,
		}

		filters = append(filters, &nameFilter)
	}

	// add the filters if they exist
	if len(filters) > 0 {
		params.Filters = filters
	}

	// actually call aws for volume information
	result, err := client.DescribeVolumes(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				err = errors.Wrapf(aerr, "error searching volumes")
				return info, err
			}
		} else {
			err = errors.Wrapf(err, "error searching volumes")
			return info, err
		}
	}

	// loop through the resulting info, and set up the info we need
	for _, vol := range result.Volumes {
		instanceId := *vol.Attachments[0].InstanceId
		deviceName := *vol.Attachments[0].Device

		i := VolInfo{
			InstanceId: instanceId,
			DeviceName: deviceName,
			VolumeId:   *vol.VolumeId,
		}

		info = append(info, i)
	}

	return info, err
}

// GetInstanceInfoMaps gets relevant information about running instances
func GetInstanceInfoMaps(awsSession *session.Session, targets []string) (infomaps InstanceInfoMaps, err error) {
	client := ec2.New(awsSession)
	name2info := make(map[string]InstanceInfo)
	id2info := make(map[string]InstanceInfo)

	infomaps.Id2Info = id2info
	infomaps.Name2Info = name2info

	filters := make([]*ec2.Filter, 0)

	params := &ec2.DescribeInstancesInput{}

	if targets != nil {
		awsnames := make([]*string, 0)

		for _, name := range targets {
			awsnames = append(awsnames, aws.String(name))
		}

		nameFilter := ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: awsnames,
		}

		filters = append(filters, &nameFilter)
	}

	// add the filters if any
	if len(filters) > 0 {
		params.Filters = filters
	}

	// Actually call aws for the information
	result, err := client.DescribeInstances(params)
	if err != nil {
		err = errors.Wrapf(err, "failed to call describe instances")
		return infomaps, err
	}

	// we get back reservations
	for _, reservation := range result.Reservations {
		// which we step through to get instances
		for _, instance := range reservation.Instances {

			var name string

			// Names, however are tags, which we have to loop through
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}

			i := InstanceInfo{
				InstanceId:   *instance.InstanceId,
				InstanceName: name,
			}

			// add the info to the maps for easy and cheap lookup later
			id2info[*instance.InstanceId] = i
			name2info[name] = i
		}
	}

	return infomaps, err
}

// GenerateNameTag generates the expected value for the tag Name based on the instance name and the device name
func GenerateNameTag(instanceName string, deviceName string) (nameTag string) {
	nameTag = fmt.Sprintf("%s_%s", instanceName, deviceName)
	return nameTag
}

// SnapshotRunningVolumes takes snapshots of the volumes attached to all currently running instances
func SnapshotRunningVolumes(awsSession *session.Session, targets []string) (err error) {
	client := ec2.New(awsSession)

	// first get the information on relevant instances
	infomaps, err := GetInstanceInfoMaps(awsSession, targets)
	if err != nil {
		err = errors.Wrapf(err, "failed to get instance info")
		return err
	}

	// volume information
	var volumes []VolInfo

	// massage targets, which are instance names (tags) into instance ID'd, which is what we need to search for attached volumes
	if targets != nil {
		ids := make([]string, 0)

		for _, name := range targets {
			ids = append(ids, infomaps.Name2Info[name].InstanceId)
		}

		// get filtered volume info
		volumes, err = GetVolumeInfo(awsSession, ids)
		if err != nil {
			err = errors.Wrapf(err, "failed to get volume info")
			return err
		}

	} else {
		// get the unfiltered volume info
		volumes, err = GetVolumeInfo(awsSession, nil)
		if err != nil {
			err = errors.Wrapf(err, "failed to get volume info")
			return err
		}
	}

	// now we can step through the returned information and set up the structs we care about
	for _, volume := range volumes {
		instanceName := infomaps.Id2Info[volume.InstanceId].InstanceName
		deviceName := volume.DeviceName
		volumeId := volume.VolumeId

		// get a tag name
		nametag := GenerateNameTag(instanceName, deviceName)

		// and a timestamp
		timestamp := time.Now().String()

		// setup the params for the snapshot call
		params := &ec2.CreateSnapshotInput{
			Description: aws.String(fmt.Sprintf("Snapshot for %s at %s", nametag, timestamp)),
			VolumeId:    &volumeId,
		}

		// actually tell AWS to make a snapshot
		result, err := client.CreateSnapshot(params)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					err = errors.Wrapf(aerr, "failed to create snapshot")
					return err
				}
			} else {
				err = errors.Wrapf(aerr, "failed to create snapshot")
				return err
			}
		}

		snapshotId := *result.SnapshotId

		// Now that we have the snapshot id, we can use it to tag the resource (can't tag on creation, have to do it afterwards)
		err = TagResource(awsSession, snapshotId, nametag, timestamp)
		if err != nil {
			err = errors.Wrap(err, "failed to tag snapshot")
			return err
		}
	}

	return err
}

// TagResource tags the given resource with the supplied information
func TagResource(awsSession *session.Session, resource string, nametag string, timestamp string) (err error) {
	client := ec2.New(awsSession)

	// set up the params for tagging
	params := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(resource),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(nametag),
			},
			{
				Key:   aws.String("Date"),
				Value: aws.String(timestamp),
			},
		},
	}

	// actually call AWS with the command to tag.  Return value is an empty {} on success, so all we care about is the error.
	_, err = client.CreateTags(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				err = errors.Wrapf(aerr, "failed to create tags")
				return err
			}
		} else {
			err = errors.Wrapf(aerr, "failed to create tags")
			return err
		}
	}

	return err
}
