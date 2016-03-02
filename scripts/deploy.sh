#!/usr/bin/env bash
set -e
kubeargs="--namespace=$KUBE_ENV"

shopt -s expand_aliases
alias kubectl="sudo docker run --rm -v /home/core/.kube:/root/.kube -v /home/core/pipelines/todo-service:/root/todo-service wernight/kubectl kubectl"

rcExist=$(kubectl get -o template rc $GO_PIPELINE_NAME --template={{.kind}} $kubeargs) || true
if [ "$rcExist" != "ReplicationController" ]; then
  #sed -i -e "s/\${BUILD_TAG}/$GO_PIPELINE_LABEL/g" scripts/rc.yml
  kubectl create -f /root/todo-service/scripts/rc.yml $kubeargs
else
  kubectl rolling-update $GO_PIPELINE_NAME --image=canthefason/$GO_PIPELINE_NAME:$GO_PIPELINE_LABEL --update-period=20s $kubeargs
fi

svcExist=$(kubectl get -o template svc $GO_PIPELINE_NAME --template={{.kind}} $kubeargs) || true
if [ "$svcExist" != "Service" ]; then
  kubectl create -f /root/todo-service/scripts/svc.yml $kubeargs
fi

sleep 30

curl -k --retry 10 --retry-delay 5 -v https://${KUBE_USER}:${KUBE_PASSWORD}@${KUBE_HOST}/api/v1/proxy/namespaces/sandbox/services/todo-service/

curl -k --silent --output /dev/stderr --write-out "%{http_code}" -v https://${KUBE_USER}:${KUBE_PASSWORD}@${KUBE_HOST}/api/v1/proxy/namespaces/sandbox/services/todo-service/

if [ "$STATUSCODE" -ne "200" ]; then
  if [ "$rcExist" != "ReplicationController" ]; then
    kubectl delete -f /root/todo-service/scripts/rc.yml $kubeargs
  fi
  exit 1
fi
