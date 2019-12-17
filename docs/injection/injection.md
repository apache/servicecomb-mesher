# Masher injection

Mesher can be used with any application language on any infrastructure. On Kubernetes, we inject Mesher into Pods, and applications can take advantage of all its features. Below we introduce the use of Mesher injector by example. 

## Before
1. Install Servcicecomb [service-center](http://servicecomb.apache.org/docs/service-center/install/#deployment-with-kubernetes)  


2. Download the example of quick_start
```bash
$ git clone https://github.com/apache/servicecomb-mesher.git
$ cd ./servicecomb-mesher/examples/quick_start
```

3. Build docker images
```bash
$ cd ./mesher_injection
$ bash build_images.sh
```

## Automatic Mesher injection  
Meshers can be automatically added to applicable Kubernetes pods using a mutating webhook admission controller provided by Mesher Injector.  
 
1. Create namespace "svccomb-system"
   ```bash
   cp ../../../deployments/kubernetes/injector/*.yaml .
   $ kubectl apply -f svccomb-system.yaml
   namespace/svccomb-system created
   ```
2. Generate Injector's certificatesigningrequest and secret
   ```bash
   $ wget https://raw.githubusercontent.com/morvencao/kube-mutating-webhook-tutorial/master/deployment/webhook-create-signed-cert.sh
   
   $ bash webhook-create-signed-cert.sh --service svccomb-mesher-injector --namespace svccomb-system --secret svccomb-mesher-injector-service-account
   ```
3. Query caBundle and fill it into "mesher-injector.yaml"
   ```bash
   $ CA_BUNDLE=$(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
   
   $ sed -i "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" mesher-injector.yaml
   ```
4. Deploy mesher injecter
   ```bash
   $ kubectl apply -f mesher-injector.yaml
   mutatingwebhookconfiguration.admissionregistration.k8s.io/svccomb-mesher-injector configured
   service/svccomb-mesher-injector created
   deployment.extensions/svccomb-mesher-injector created
   ```
5. Deploy examples to Kubernetes  
   ```bash
   $ kubectl apply -f svccomb-test.yaml
   namespace/svccomb-test created
   
   $ kubectl -n svccomb-test apply -f calculator.yaml
   service/calculator created
   deployment.extensions/calculator-python created
   
   $ kubectl -n svccomb-test apply -f webapp.yaml
   service/webapp created
   deployment.extensions/webapp-node created
   ```
6. Validated results
   ```bash
   $ kubectl get svc 
   NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
   calculator   ClusterIP   10.104.2.143     <none>        5000/TCP         3m43s
   kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP          42d
   webapp       NodePort    10.104.134.148   <none>        5001:30062/TCP   3m35s
   ```
   Open the page "http://127.0.0.1:30062" in your browser, enter your height and weight in the input boxes, and click the submit button, you can see the BMI results about you. 