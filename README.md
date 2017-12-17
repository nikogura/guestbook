# guestbook
[![Circle CI](https://circleci.com/gh/nikogura/guestbook.svg?style=shield)](https://circleci.com/gh/nikogura/guestbook)

[![Go Report Card](https://goreportcard.com/badge/github.com/nikogura/guestbook)](https://goreportcard.com/report/github.com/nikogura/guestbook)

A really simple implementation of your classic guestbook hello world app.

Why?  I wanted to play with some database stuff in go.

I also needed a frontend/backend app to mess with as a subject for some infrastructre projects.

This isn't going to win any points, but as long as it stands up and does *something*, it's good enough.

## Running the Code Locally

### Prerequisites

* A properly set up Go SDK

* Git, and access to GitHub

* Postgres running locally on port 5432

* A Database called 'guestbook' with a user called 'guestbook' and a password 'guestbook'.  Imaginative, I know


### Installation
Be sure to set $GOPATH, and put $GOPATH/bin into your $PATH

Once you have one, install the guestbook app with:

        go get github.com/nikogura/guestbook
        
        
The service expects a config file at ```/etc/guestbook/guestbook.json```.  It's contents should look like:

        {
          "state": {
            "manager": {
              "type": "gorm",
        	  "dialect": "postgres",
              "connect_string": "postgresql://guestbook:guestbook@localhost:5432/guestbook?sslmode=disable"
            }
          },
          "server": {
            "addr": "0.0.0.0:8080"
          }
        }
        
        
### Running
        
You can run the service with:

        guestbook run
        
You can run the snapshot tool with:

        guestbook snapshot <space separated list of instance names>
        
Snapshotting only works with AWS, not VirtualBox.  If you call ```guestbook snapshot``` without additional arguments, it will snapshot *all* your running instances.  Use with care.

You can get help by running:

        guestbook help


## Running the Stack with Vagrant and VirtualBox

### Prereqs
 
* VirtualBox  https://www.virtualbox.org/wiki/Downloads

* Vagrant https://www.vagrantup.com/downloads.html

### Run Service

From the root of this cloned repository run the following:

        cd vagrant
        
        vagrant up
        
### Test

Point a browser at http://localhost:8080/guestbook/

Enjoy!


## Running in AWS

### Prereqs

* Make sure you have your appropriate AWS creds positioned in ```~/.aws/credentials``` .

* Region config in ```~/.aws/config```.  This test is hardcoded for the region ```us-east-1```.

* Terraform.  You can install manually via: [Terraform Download Page](https://www.terraform.io/downloads.html)  Alternately, on a Mac, you can just run ```brew install terraform```.  Your choice.

* SSH keys installed in the normal location (```~/.ssh/id_rsa``` and ```~/.ssh/id_rsa.pub```)

### Configuration

The terraform config as writen is set to allow *my* ip access to the backend, not yours.  You'll need to modify the file ```terraform/variables.tf``` with your IP information.

*NOTE: There is a 'bug' with the external ELB config in that Terraform does not like setting an elb in multiple subnets in the same availability zone.  Likewise, it doesn't like you setting availability zones and subnets at the same time.  Depending on your user config, there is a chance that the front end ELB will not spring up with all the relevant availability zones enabled for the frontend servers.*

*If this happens, you just need to add the availability zones to the ELB via the console.  There are ways to fix this, but so far I have not found one that entirely satisfies me, or works every time in every situation.  This message will disappear if I find an acceptable solution*

### Run Service

To spin it all up, from the root of this checked out repository run Terraform thusly:

        cd terraform
        
        terraform apply
    
### Test

The Terraform command above will output the front end ELB's url.  Point a browser at that and enjoy!

*NOTE: It can take a few minutes for DNS to pick up the change and activate the above url.  Please be patient.*



