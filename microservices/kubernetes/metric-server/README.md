# Kubernetes Metric Server Installation

The guide assumes that the Kubernetes environment is already configured, either with minikube or a real node cluster.

1. Inside metric-server folder run ```kubectl -n kube-system apply -f .```
2. 



   of [Kong API Ingress Controller](https://docs.konghq.com/kubernetes-ingress-controller/latest/get-started/)
2. To assign an external IP address to KIC (if minikube is used),
   run: ```kubectl patch svc kong-gateway-proxy -n kong -p '{"spec":{"externalIPs":["192.168.58.2"]}}'```
3. Apply all the files to the Kubernetes cluster by running (from this folder): ```kubectl apply -f .```
4. To allow image download from minio, add to the bucket ```images``` this [policy](../utils/download_policy.json)
5. To access the various consoles, get the public IP address by using ```minikube service <service-name> --url```