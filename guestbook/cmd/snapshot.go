// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/nikogura/guestbook/guestbook/snapshot"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Snapshots the volumes attached to the guestbook instances.",
	Long: `
Snapshots the volumes attached to the guestbook instances.

Tags them with the following:

Name: <instance name>_<device name>
Date: <timestamp>
`,
	Run: func(cmd *cobra.Command, args []string) {

		awsSession := snapshot.Ec2Session()

		if len(args) == 0 {
			fmt.Printf("Snapshotting All volumes.\n")
			err := snapshot.SnapshotRunningVolumes(awsSession, nil)
			if err != nil {
				log.Printf("error snapshotting volumes: %s", err)
				os.Exit(1)
			}
			fmt.Printf("Done.\n")

		} else {
			fmt.Printf("Snapshotting volumes for: %s\n", args)
			err := snapshot.SnapshotRunningVolumes(awsSession, args)
			if err != nil {
				log.Printf("error snapshotting volumes: %s", err)
				os.Exit(1)
			}
			fmt.Printf("Done.\n")
		}
	},
}

func init() {
	RootCmd.AddCommand(snapshotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// snapshotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// snapshotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
