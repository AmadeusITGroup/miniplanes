# What is this?

`miniapp` is a synthetic workload to test Federation-v2. The _architecture_ is (for the moment) very simple:


`storage`

`itineraries-server` is the cpu intensive computing resource which speaks with the DB (currently a MongoDB)

`ui` is the `web` server.



# Deployment
At the moment `miniApp` runs in an Openshift/Kubernetes cluster.



## How to populate DB (locally)


## To deploy it locally
A script to deploy it locally is supplied in `hack/deploy_all.sh`
