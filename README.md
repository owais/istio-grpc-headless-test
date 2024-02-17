## Demonstrate gRPC + headless service issues with istio


### Steps to reproduce

1. Deploy server components `kubectl apply -f server.yaml`
2. Switch to istio-test namespace. `kubens istio-test`
3. Wait for all three server pods to become healthy
    --- step1 screenshot here
4. Deploy client components `kubectl apply -f client.yaml`
5. Observe client is able to send messages to the server pods with `kubectl logs -f -l=component=client`
    --- step 5 screenshot here
6. Scale down server pods with `kubectl scale sts istio-grpc-test-server --replicas 0`
7. Observe errors in client pod logs
  --- step 7 screenshot
8. Look at istio proxy endpoints. Most of the time it'll still list old server pod IPs as healthy. `istioctl proxy-config endpoints <client-pod-name> | grep server`
   --- step 8 screenshot
8. Scale up server pods back to three replicas.  `kubectl scale sts istio-grpc-test-server --replicas 3`
9. Wait for new server pods to come up and become healthy and take note of the new IPs.
  -- step 9 screenshot
10. Describe service to confirm k8s service has updated its endpoints to the new pod IPs
 -- step 10
10. List client proxy endpoints again as in step 8 and notice they are still pointing to old IPs
 -- step 11
11. Look at client pods logs again and confirm that the errors have not resolved even though replacement server pods are up and healthy
 -- step 12

 At this point pretty much the only way to recover the service is to restart the client pod so it gets the new server IPs. Sometimes making some changes to the k8s `service` resource or related istio resources also triggers something and updates client endpoints to the new IPs but this does not work reliably.
 

 ### Deployment have the same issue
 
 I've deployed the servers as a statefulset as that is the closest setup to my real world scenario but it doesn't really matter. I've been able to reproduce it with deployment as well as long as the pods are exposed with a headless service (`clusterIP: None`). 
 

## What works

### Headful services
 
 When using a headful service instead of headless, none of this happensa and the istio proxy is able to discover new endpoints as soon as they become healthy. It also removes the old endpoints as soon as server pods are deleted. This allows the client to recover once server pods become available again. 
 
 It doesn't matter whether the server pods are from a stateful or a deployment. As long as they are exposed via a headless service (`clientIP: None`), I can reliably reproduce the issue. 


### Marking service port as TCP

This only happens when the service port app protocol is set to gRPC by prefixing the port name with `grpc-`. If the port name is prefixed with `tcp-`, everything works as expected even with a headless service.
