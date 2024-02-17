## Demonstrate gRPC + headless service issues with istio

issue: https://github.com/istio/istio/issues/49391


### Steps to reproduce

1. Deploy server components `kubectl apply -f server.yaml`
2. Switch to istio-test namespace. `kubens istio-test`
3. Wait for all three server pods to become healthy
<img width="690" alt="step-1" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/e8590713-6f23-4d4d-919d-302594415a77">

4. Deploy client components `kubectl apply -f client.yaml`
5. Observe client is able to send messages to the server pods with `kubectl logs -f -l=component=client`
<img width="998" alt="step-5" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/0d7a0fde-313b-4af2-a2ba-62ec6175b25e">

6. Scale down server pods with `kubectl scale sts istio-grpc-test-server --replicas 0`
7. Observe errors in client pod logs
  <img width="1000" alt="step-7" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/efcaebd8-72d7-4cf0-b259-7cff832c5e72">
  
8. Look at istio proxy endpoints. Most of the time it'll still list old server pod IPs as healthy. `istioctl proxy-config endpoints <client-pod-name> | grep server`
<img width="1160" alt="step-8" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/8349a1fe-5299-458c-9641-81b228446b80">

9. Scale up server pods back to three replicas and Wait for new server pods to come up and become healthy and take note of the new IPs.  `kubectl scale sts istio-grpc-test-server --replicas 3`
<img width="1113" alt="step-9" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/151511e9-4e84-4f17-af22-a71d0b5c2f2e">

10. Describe service to confirm k8s service has updated its endpoints to the new pod IPs
<img width="783" alt="step-10" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/de4da675-1cf3-4cc1-9a33-541961be17a6">

11. List client proxy endpoints again as in step 8 and notice they are still pointing to old IPs
<img width="1278" alt="step-11" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/65c770d6-cd4c-4c72-b127-c77f86ab9dbc">

12. Look at client pods logs again and confirm that the errors have not resolved even though replacement server pods are up and healthy
<img width="994" alt="step-12" src="https://github.com/owais/istio-grpc-headless-test/assets/46186/574cc15b-a1a7-4400-baa5-d785b828a4f8">


 At this point pretty much the only way to recover the service is to restart the client pod so it gets the new server IPs. Sometimes making some changes to the k8s `service` resource or related istio resources also triggers something and updates client endpoints to the new IPs but this does not work reliably.
 

 ### Deployment have the same issue
 
 I've deployed the servers as a statefulset as that is the closest setup to my real world scenario but it doesn't really matter. I've been able to reproduce it with deployment as well as long as the pods are exposed with a headless service (`clusterIP: None`). 

### accessing individual pods vs service behaves the same
It doesn't matter whether the client tries to connect to the headless service (`istio-grpc-test-server.istio-test.svc.cluster.local`) or a specific pod (`istio-grpc-test-server-0.istio-grpc-test-server.istio-test.svc.cluster.local`). Both cases behave exactly the same.

## What works

### Headful services
 
 When using a headful service instead of headless, none of this happensa and the istio proxy is able to discover new endpoints as soon as they become healthy. It also removes the old endpoints as soon as server pods are deleted. This allows the client to recover once server pods become available again. 
 
 It doesn't matter whether the server pods are from a stateful or a deployment. As long as they are exposed via a headless service (`clientIP: None`), I can reliably reproduce the issue. 


### Marking service port as TCP

This only happens when the service port app protocol is set to gRPC by prefixing the port name with `grpc-`. If the port name is prefixed with `tcp-`, everything works as expected even with a headless service.
