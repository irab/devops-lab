#!/bin/bash

## install kubectl
gcloud components install kubectl --quiet

## get cluster credentials (assumes you have created a default zonal GKE cluster via the console)
gcloud container clusters get-credentials standard-cluster-1 --zone us-central1-a

## apply kubernetes manifests
kubectl apply -f kubernetes 