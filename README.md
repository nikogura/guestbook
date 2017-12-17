# guestbook
[![Circle CI](https://circleci.com/gh/nikogura/guestbook.svg?style=shield)](https://circleci.com/gh/nikogura/guestbook)

[![Go Report Card](https://goreportcard.com/badge/github.com/nikogura/guestbook)](https://goreportcard.com/report/github.com/nikogura/guestbook)

A really simple implementation of your classic guestbook hello world app.

Why?  I wanted to play with some database stuff in go.

I also needed a frontend/backend app to mess with as a subject for some infrastructre projects.

This isn't going to win any points, but as long as it stands up and does *something*, it's good enough.

## Running the Code

The guestbook is written in Go.  You'll need a Go SDK.

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
        
To run it locally a Postgres database must be running on localhost:5432, and you'll need to create the database 'guestbook', and user 'guestbook', with password 'guestbook'.
        
        
You can run the service with:

        guestbook run
        
You can run the snapshot tool with:

        guestbook snapshot <space separated list of instance names>
        
Snapshotting only works with AWS, not VirtualBox.  If you call ```guestbook snapshot``` without additional arguments, it will snapshot *all* your running instances.  Use with care.


## Running with Vagrant and VirtualBox

### Prereqs
 
Make sure you have Vagrant and Virtualbox installed

* Install VirtualBox  https://www.virtualbox.org/wiki/Downloads

* Install Vagrant https://www.vagrantup.com/downloads.html

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

* Install Terraform.  You can do this manually via: [Terraform Download Page](https://www.terraform.io/downloads.html)  Alternately, on a Mac, you can just run ```brew install terraform```.  Your choice.

* Make sure you have ssh keys installed in the normal location (```~/.ssh/id_rsa``` and ```~/.ssh/id_rsa.pub```)

### Run Service

To spin it all up, from the root of this checked out repository run Terraform thusly:

        cd terraform
        
        terraform apply
    
### Test

The Terraform command above will output the front end ELB's url.  Point a browser at that and enjoy!

*NOTE: It can take a few minutes for DNS to pick up the change and activate the above url.  Please be patient.*



