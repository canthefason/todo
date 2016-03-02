#!/usr/bin/env bash
set -e
kubeargs="--namespace=$KUBE_ENV"

shopt -s expand_aliases
alias kubectl="sudo docker run --rm -v /home/core/.kube:/root/.kube -v /home/core/pipelines/todo-service:/root/todo-service wernight/kubectl kubectl"

rcExist=$(kubectl get -o template rc $GO_PIPELINE_NAME --template={{.kind}} $kubeargs) || true
if [ "$rcExist" != "ReplicationController" ]; then
  envsubst < scripts/rc.yml > kubetemp.yml
  cat kubetemp.yaml
  kubectl create -f /root/todo-service/kubetemp.yml $kubeargs
else
  kubectl rolling-update $GO_PIPELINE_NAME --image=canthefason/$GO_PIPELINE_NAME:$GO_PIPELINE_LABEL --update-period=20s $kubeargs
fi

svcExist=$(kubectl get -o template svc $GO_PIPELINE_NAME --template={{.kind}} $kubeargs) || true
if [ "$svcExist" != "Service" ]; then
  kubectl create -f /root/todo-service/scripts/svc.yml $kubeargs
fi
