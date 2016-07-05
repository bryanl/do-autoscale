# do-autoscale

Deploy an autoscale to automatically scale your servers based on load.

# Install

## Automated Droplet configuration

If you use [doctl](https://github.com/digitalocean/doctl), you can create a server using the following: `doctl compute droplet create <name> --region nyc1 --size 4gb --image ubuntu-16-04-x64 --user-data-file userdata.sh --ssh-keys <your key id>`.  `userdata.sh` can be retrieved from [https://s3.pifft.com/autoscale/userdata.sh](https://s3.pifft.com/autoscale/userdata.sh).

## Manual Droplet configuration

1. Create a Ubuntu 16.04 Droplet. The autoscaler was designed with a 4GB Droplet in mind.
1. Download `autoscalectl` from https://s3.pifft.com/autoscale/autscalectl, and mark it executeable

## Running setup

Run `autoscalectl setup`

The setup process requires:

* DigitalOcean access token
* Fully qualified domain name for SSL configuration
* TLS certs/key (letsencrypt can also be used)
* Password for accessing the site

If you choose to use your own key/certs, place them in the /etc/autoscale/ssl directory and name them autoscale.crt and autoscale.key. You should ensure to `chmod 600 autoscale.key` for safetey.

# Usage

Navigate to https://your.host.com and use `autoscale` as your user name, and the password you chose earlier as your password.

# Upgrading

Run `autoscalectl update` to update the autoscaler..


