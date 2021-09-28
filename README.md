# go-pets

## How to run

run `gcloud auth application-default login` to set account key for terraform(note that your account must have sufficient
permission to create required resources)

export environment variable `PROJECT` set to gcp project to which service will be deployed

Run `./install.sh` to create resources and deploy service to App Engine

## How to destroy

Run `./uninstall.sh`
