# Kubernetes Metric Server Installation

The guide assumes that the Kubernetes environment is already configured, either with minikube or a real node cluster.

1. Inside metric-server folder run ```kubectl -n kube-system apply -f .```
2. ```kubectl top pods``` to see the metrics of pods