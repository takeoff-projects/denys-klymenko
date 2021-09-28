#!/bin/bash

gcloud config set project $PROJECT
terraform plan -chdir=terraform -destroy
terraform destroy -chdir=terraform
firebase firestore:delete --all-collections
gsutil rm -r gs://denys-klymenko-pets