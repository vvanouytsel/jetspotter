# Development

## Using the justfile

The [just](https://github.com/casey/just) command runner is used to automate simple tasks.
List all available options:

```bash
$ just list
```

Build the application.

```bash
$ just build
```

Run the application.

```bash
$ just run
```

Test the application.

```bash
$ just test
```

## Devcontainers

You can run this code in a [devcontainer](https://containers.dev/overview). That allows you to run the code locally in a container without needing any dependencies, except `docker`.

## Setting up a local development environment

[Minikube](https://minikube.sigs.k8s.io/docs) can be used for local development.

Install the binary.

```bash
$ curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64
$ sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64
$ minikube version  
minikube version: v1.34.0
commit: 210b148df93a80eb872ecbeb7e35281b3c582c61
```

Add your user to the `libvirt` group, make sure to log out and log back in.

```bash
$ sudo usermod -aG libvirt $(whoami)
$ groups $(whoami)  
myuser : myuser wheel libvirt

# Log out and log back in, verify that you are part of the libvirt group.
$ groups
myuser : myuser wheel libvirt
 ```

> Note: All steps below are automated via the 'just create-dev' command.

Start the local Kubernetes environment.

```bash
$ minikube start
```

Build the container image.

```bash
$ podman build -t ghcr.io/vvanouytsel/jetspotter:dev .
```

Push the container image to the minikube image store.

```bash
$ minikube image load ghcr.io/vvanouytsel/jetspotter:dev
```

Apply the manifests.

```bash
$ kubectl apply -f development
namespace/dev created
secret/jetspotter created
configmap/jetspotter created
service/jetspotter created
deployment.apps/jetspotter created
deployment.apps/postgres created
service/postgres created
```

You can now connect directly to the database if you have `psql` installed.

```bash
$ kubectl exec -ti -n dev $(kubectl get pods  -l app=postgres --no-headers -o custom-columns=":metadata.name" -n dev) -- psql -U postgres
psql (16.3, server 16.6 (Debian 16.6-1.pgdg120+1))
Type "help" for help.

postgres=# 
```

### Cleaning up

Delete your local `minikube` cluster.

```
$ minikube remove
```

## Forward traffic inside cluster to local jetspotter

Using [telepresence](https://www.telepresence.io) it is possble to intercept traffic going to jetspotter in a remote environment and forward it to your local instance. This can be helpful if you want to quickly test changes against an already established environment.

Download the telepresence binary on your local development machine.

```bash
# 1. Download the latest binary (~95 MB):
$ sudo curl -fL https://app.getambassador.io/download/tel2oss/releases/download/v2.21.1/telepresence-linux-amd64 -o /usr/local/bin/telepresence

# 2. Make the binary executable:
$ sudo chmod a+x /usr/local/bin/telepresence
```

Make sure your KUBECONFIG points to the remote Kubernetes cluster where you want to intercept the traffic.  
Install the telepresence helm chart.

```bash
$ telepresence helm install
```

Connect to the namespace where `jetspotter` is running.

```bash
# In this example I am using minikube and jetspotter is running in the 'dev' namespace
$ telepresence connect --namespace dev
Connected to context minikube, namespace dev (https://192.168.39.151:8443)
```

Verify that services are found.

```bash
 ❯ telepresence list
jetspotter: ready to intercept (traffic-agent not yet installed)
postgres  : ready to intercept (traffic-agent not yet installed)
```

List the ports available for the service.

```bash
$ kubectl get service jetspotter --output yaml
```

Intercept traffic going to the service and proxy it locally instead.

```bash
# telepresence intercept <service-name> --port <local-port>[:<remote-named-port>] --env-file <path-to-env-file>
$ telepresence intercept jetspotter --port 7070:metrics --env-file jetspotter.env
Using Deployment jetspotter
   Intercept name         : jetspotter
   State                  : ACTIVE
   Workload kind          : Deployment
   Destination            : 127.0.0.1:7070
   Service Port Identifier: metrics/TCP
   Volume Mount Error     : remote volume mounts are disabled: sshfs is not installed on your local machine
   Intercepting           : all TCP connections
```

Start the application locally on your development machine.  
All traffic going to the jetspotter in your target cluster is now proxied to your locally running instance.

```bash
$ podman run --env-file jetspotter.env ghcr.io/vvanouytsel/jetspotter:dev
2025/01/14 21:50:00 Spotting the following aircraft types within 30 kilometers: [ALL]
2025/01/14 21:50:00 Serving API on port 8085 and path /api
2025/01/14 21:50:00 Serving metrics on port 7070 and path /metrics
2025/01/14 21:50:00 No new matching aircraft have been spotted.
```

### Cleaning up

To clean up, list the active connections and delete it.

```bash
 $ telepresence list
jetspotter: intercepted
   Intercept name         : jetspotter
   State                  : ACTIVE
   Workload kind          : Deployment
   Destination            : 127.0.0.1:7070
   Service Port Identifier: metrics/TCP
   Intercepting           : all TCP connections
postgres  : ready to intercept (traffic-agent not yet installed)

$ telepresence leave jetspotter
```

Uninstall the traffic manager chart.

```bash
$ telepresence helm uninstall
Traffic Manager uninstalled successfully
```
