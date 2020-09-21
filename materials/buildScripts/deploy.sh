#!/bin/bash

set -xe
pwd
COMMIT_ID=$(git rev-parse --short $GIT_COMMIT)

if [ "$IMAGE_TAG" = "" ]
then
   TAG=$COMMIT_ID-$BUILD_NUMBER
else
   TAG=$IMAGE_TAG
fi
echo $TAG

rm -rf /home/jenkins/helm/

#clone helm charts from repo
git clone git@code.buildsupply.com:devops/helm.git /home/jenkins/helm/
ls /home/jenkins/helm/dev/mm-api/
#updating kubeconfig
kubectl config use-context arn:aws:eks:ap-south-1:317596419736:cluster/dev-ecom-im-cluster
#deploy on kubernetes using helm
#helm upgrade mm-api /home/jenkins/helm/dev/mm-api/ --set=deployment.image.tag=$TAG -n dev