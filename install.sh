#!/bin/bash
echo "Setting environment for $PROJECT" &&
export TF_VAR_project_id=$PROJECT
gcloud config set project $PROJECT &&
gcloud components install app-engine-go &&
gcloud app create

gsutil mb gs://denys-klymenko-pets

return_value=$?

if [ $return_value = 0 ]; then
    echo "folder exist"
else
    echo "failed to create bucket"
fi

cd functions/images

go mod vendor

zip -r images.zip .

cd ../..

swag init --g api/pets/api.go -o ./docs

#gcloud builds submit .

cd terraform
terraform init  -backend-config="bucket=denys-klymenko-pets"  &&
terraform apply -auto-approve &&
cd ..
gcloud app deploy -q
