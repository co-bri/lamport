# Cloud Install Sctipts

## Before running

1) Setup cloud instance (for instance using Terraform scripts)

## To Instnall and run lamport in the cloud
2) Move file into this directory that has host names and IPs
- e.g. lamport_hosts.txt output from lamportaws.tf Terrafirn script
3) Run lamport-start.sh
- This will install lamport on the hosts in the "lamport_hostst.txt" file
- Only currently works for Ubuntu


