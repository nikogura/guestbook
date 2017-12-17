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
func Ec2Session() (awsSession *session.Session) {
	awsSession = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return
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

// GetVolumeInfo gets relevant info about volumes currently in existance
func GetVolumeInfo(awsSession *session.Session, targets []string) (info []VolInfo, err error) {
	client := ec2.New(awsSession)
	info = make([]VolInfo, 0)

	// list volumes
	input := &ec2.DescribeVolumesInput{}

	result, err := client.DescribeVolumes(input)
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

	for _, vol := range result.Volumes {
		instanceId := *vol.Attachments[0].InstanceId
		deviceName := *vol.Attachments[0].Device

		if targets != nil { // if we're passed a target list, only append the info if it's one of the targets
			if StringInSlice(instanceId, targets) {
				i := VolInfo{
					InstanceId: instanceId,
					DeviceName: deviceName,
					VolumeId:   *vol.VolumeId,
				}

				info = append(info, i)
			}

		} else { // otherwise grab 'em all
			i := VolInfo{
				InstanceId: instanceId,
				DeviceName: deviceName,
				VolumeId:   *vol.VolumeId,
			}

			info = append(info, i)
		}
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

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

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

			if targets != nil { // if we have targets, only return info on the targets
				if StringInSlice(name, targets) {
					i := InstanceInfo{
						InstanceId:   *instance.InstanceId,
						InstanceName: name,
					}

					id2info[*instance.InstanceId] = i
					name2info[name] = i
				}
			} else { // return everything
				i := InstanceInfo{
					InstanceId:   *instance.InstanceId,
					InstanceName: name,
				}

				id2info[*instance.InstanceId] = i
				name2info[name] = i
			}
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

	infomaps, err := GetInstanceInfoMaps(awsSession, targets)
	if err != nil {
		err = errors.Wrapf(err, "failed to get instance info")
		return err
	}

	var volumes []VolInfo

	if targets != nil {
		// volumes need id's not names
		ids := make([]string, 0)

		for _, name := range targets {
			ids = append(ids, infomaps.Name2Info[name].InstanceId)
		}

		volumes, err = GetVolumeInfo(awsSession, ids)
		if err != nil {
			err = errors.Wrapf(err, "failed to get volume info")
			return err
		}

	} else {
		volumes, err = GetVolumeInfo(awsSession, nil)
		if err != nil {
			err = errors.Wrapf(err, "failed to get volume info")
			return err
		}
	}

	for _, volume := range volumes {
		instanceName := infomaps.Id2Info[volume.InstanceId].InstanceName
		deviceName := volume.DeviceName
		volumeId := volume.VolumeId
		nametag := GenerateNameTag(instanceName, deviceName)
		timestamp := time.Now().String()

		input := &ec2.CreateSnapshotInput{
			Description: aws.String(fmt.Sprintf("Snapshot for %s at %s", nametag, timestamp)),
			VolumeId:    &volumeId,
		}

		result, err := client.CreateSnapshot(input)
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

	input := &ec2.CreateTagsInput{
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

	_, err = client.CreateTags(input)
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
