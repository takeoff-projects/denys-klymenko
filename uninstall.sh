#!/bin/bash

export TF_VAR_project_id=$PROJECT
gcloud config set project $PROJECT
cd terraform
terraform plan  -destroy
terraform destroy
firebase firestore:delete --all-collections
gsutil rm -r gs://denys-klymenko-pets