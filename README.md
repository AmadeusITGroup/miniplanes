##

`backend` exposes the REST API for the server DB

## To create your DB


## To deploy your app in minikube
$ kubectl get pods
NAME            READY     STATUS    RESTARTS   AGE
mongo-3-rbqn9   2/2       Running   0          44m


kubectl port-forward pod/mongo-3-rbqn9 :27017
