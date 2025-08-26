# Kubectl imageof - Quickly retrieve image(s) of pod(s) instead of grepping out from the description

## Usage

```bash
Usage:
 Quickly retrieve image(s) of pod(s) instead of grepping out from the description.

Syntax:
 kubectl imageof POD_NAME -n NAMESPACE
 kubectl imageof -A -n NAMESPACE

Output:
 CONTAINER: IMAGE

Options:
 -A, --all                List images of all pods in the namespace
 -h, --help               Print plugin usage
 -n, --namespace string   Namespace of the pod(s) (default "default")
 -t, --tree               Show tree view for multiple pods
 -v, --verbose            Show pod name in output

Example:
 $ kubectl imageof redis -n redis
 redis: redis

 $ kubectl imageof redis -n redis -v
 Pod redis:
  redis: redis

 $ kubectl imageof -A -n harbor
 Pod: harbor-core-75dd796d56-8gpld
  core: goharbor/harbor-core:v2.13.2
 Pod: harbor-database-0
  database: goharbor/harbor-db:v2.13.2
 Pod: harbor-jobservice-67f4b8bf4f-bzmxs
  jobservice: goharbor/harbor-jobservice:v2.13.2
 Pod: harbor-nginx-5d755775d9-crx62
  nginx: goharbor/nginx-photon:v2.13.2
 Pod: harbor-portal-856bfddd77-qnbkm
  portal: goharbor/harbor-portal:v2.13.2
 Pod: harbor-redis-0
  redis: goharbor/redis-photon:v2.13.2
 Pod: harbor-registry-565b4b6c6c-wx9vn
  registry: goharbor/registry-photon:v2.13.2
  registryctl: goharbor/harbor-registryctl:v2.13.2
 Pod: harbor-trivy-0
  trivy: goharbor/trivy-adapter-photon:v2.13.2
```

## Build and Test on Local

### Build exe

> Important: the binary name must be prefixed with kubectl- so that kubectl recognizes it as a plugin.

```bash
go build kubectl-imageof main.go
```

### Put it into the `$PATH`

```bash
chmod +x kubectl-imageof
mv kubectl-imageof ~/.local/bin
```

Make sure `~/.local/bin` is in your PATH.

```bash
echo $PATH
```

Or add it via,

```bash
export PATH=$PATH:~/.local/bin
```

### Verify kubectl detects the plugin

```bash
kubectl plugin list
```

You should see something like:

```bash
The following compatible plugins are available:

/home/harsh/.local/bin/kubectl-imageof
```

Now try commands like:

For a specific pod:

```bash
kubectl imageof redis -n redis
```

For all pods in a namespace:

```bash
kubectl imageof -A -n harbor
```

With tree view:

```bash
kubectl imageof -A -n harbor -t
```

With verbose:

```bash
kubectl imageof redis -n redis -v
```
