# Use cases

## Data Sovereignity
A more and more asked requirements is that data should be collected, stored and handled according to a specific governance structure (laws, agreements, contracts). [Data sovereignity](https://en.wikipedia.org/wiki/Data_sovereignty) is a complex; topic for the sake of the discussion, we consider only the data location aspects.

* Federation-V2 concept: federatedService (to access the specific data)
* `miniapp` usecase: 2 or more clusters with some data (for example an airline) stored exclusively in one cluster or for example _american_ airlines data stored exclusively in US located clusters.

## Specific Workload Cluster (aka clusters as cattle)

By design Federation-V2 offers the ability to place specific workloads to specific clusters. For exampple one may reserve high intesive computing workload with clusters where GPU are available, DBs workload in clusters where  a lot of storage space is offered, keeping very short and volatile workloads for public clouds with [preemptible vms](https://cloud.google.com/preemptible-vms/) or [spot instances](https://aws.amazon.com/ec2/spot/). Having the ability

* Federation-V2 concept: placementPolicy, workload customization (overriding),
* `miniapp` use case: one cluster with all the storage (_pet_ cluster) and one or more cluster with computing resource only.

## High Availability

[Federation-V2 control plane](https://github.com/kubernetes-sigs/federation-v2#concepts) is obviously a single point of failure. A possibile mitigation could be a

* Federation-V2 concept: cluster registry to span cross federations
* `miniapp` usecase: none

## Extensibility

One of the strongest point of Federation-V2

* Federation-V2: federated controllers and Federated CRDs
* `miniapp` use case: needs to create a federated controller in the control plane which operates the custom controllers in the target clusters. For example for DB replication or backup.

## Other systems

Currently `miniapp` is focused only on Kubernetes and Openshift in a near future we may want to investigate other services as well

* [Federation Prometheus](https://prometheus.io/docs/prometheus/latest/federation/)
* [Spiffe federation](https://blog.scytale.io/federating-spiffe-7d7db8040c3)
* [Istio](https://blog.openshift.com/combining-federation-v2-and-istio-multicluster/)
