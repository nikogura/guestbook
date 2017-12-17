# guestbook
[![Circle CI](https://circleci.com/gh/nikogura/guestbook.svg?style=shield)](https://circleci.com/gh/nikogura/guestbook)

[![Go Report Card](https://goreportcard.com/badge/github.com/nikogura/guestbook)](https://goreportcard.com/report/github.com/nikogura/guestbook)

A really simple implementation of your classic guestbook hello world app.

Why?  I wanted to play with some database stuff in go.

I also needed a frontend/backend app to mess with as a subject for some infrastructre projects.

This isn't going to win any points, but as long as it stands up and does *something*, it's good enough.


## Running with Vagrant and Virtualbox

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
    
*NOTE: Occasionally, depending on your config, the external ELB needs to be manually tweaked to include the proper availability zones.  The author has some funky config that may or may not be to blame.  If the front end ELB doesn't come up healthy in a reasonable amount of time, this is probably the cause. This warning will go away once the author pins the problem down to a root cause and corrects for it.*

### Test

The Terraform command above will output the front end ELB's url.  Point a browser at that and enjoy!

*NOTE: It can take a few minutes for DNS to pick up the change and activate the above url.  Please be patient.*



