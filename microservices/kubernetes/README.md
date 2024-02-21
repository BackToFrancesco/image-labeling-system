# Kubernetes Cluster Installation Process

The guide assumes that the Kubernetes environment is already configured, either with minikube or a real node cluster.

1. Run
   the installation
   of [Kong API Ingress Controller](https://docs.konghq.com/kubernetes-ingress-controller/latest/get-started/)
2. Follow the instructions to install MongoDB with autoscaling:
3. To assign an external IP address to KIC (if minikube is used),
   run: ```kubectl patch svc kong-gateway-proxy -n kong -p '{"spec":{"externalIPs":["192.168.49.2"]}}'``` (you may need
   to change the IP address, based on what's available on your side)
4. Apply all the files to the Kubernetes cluster by running (from this folder): ```kubectl apply -f .```
5. To allow image download from minio, add to the bucket ```images``` this [policy](../utils/download_policy.json)
6. To access the various consoles (RabbitMQ, MongoExpress, MinIO), get the public IP address by
   using ```minikube service <service-name> --url```
7. To allow auto-scaling of the 2 main microservices, follow the steps of this [guide](./metric-server/README.md)