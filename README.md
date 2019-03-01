# What is `miniapp`?

`miniapp` is a synthetic workload for Kubernetes/Openshift cluster. It has been designed and developed to represent Amadeus test-case for [Kubernetes Federation-v2](https://github.com/kubernetes-sigs/federation-v2) use cases. It's also simple enough to be as a simple _smoke-test_ in e2e regression and demos.


Currently  it's a minimalist [3-tiers application](https://en.wikipedia.org/wiki/Multitier_architecture#Three-tier_architecture), with usual `presentation tier`, `application tier`, `data tier`. Three REST services implements the three tiers and each REST service is backed by a process (container):

* data tier is `storage`
* application tier is `itineraries-server`
* presentation tier is `ui`.

* `storage` is a wrapper for the real DB data (currently a MongoDB).
* `itineraries-server` is the potentially cpu intensive application which computes itineraries, graph computation.
* `ui` is the _web_ server.

All data used in `miniapp` come from [openflights.org/data.html](https://openflights.org/data.html), all credits go to them.

More documentations can be found in [docs folder](./docs)

## Deployment

At the moment `miniapp` runs in an Openshift/Kubernetes cluster. 
The manifest files are available in `.../manifests/k8s` folder.

Due to its simplicity `miniapp` can  easily deployed locally on a system.
To deploy it locally some scripts have been written in `hack` folder.

## To deploy in a Kubernetes/Openshift cluster

Manifest files are in `deployment` directory.

## To deploy it locally

`miniapp` can be deployed locally in a `minikube`/`kind` environment. One can find a script in `hack/deploy_all.sh`

## How to contribute

See documentation [here](./docs/contrinute.md)