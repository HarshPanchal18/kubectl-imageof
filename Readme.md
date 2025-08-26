# Kubectl imageof

Quickly retrieve image(s) of pod(s) instead of grepping out from the description

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

### Build Go bin

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

## Publishing

### Generate bins for supporting architecture

```bash
GOOS=linux   GOARCH=amd64 go build -o kubectl-imageof-linux-amd64
GOOS=darwin  GOARCH=amd64 go build -o kubectl-imageof-darwin-amd64
GOOS=darwin  GOARCH=arm64 go build -o kubectl-imageof-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o kubectl-imageof-windows-amd64.exe
```

### Package each binary into .tar.gz (or .zip for Windows)

```bash
tar -czf kubectl-imageof-linux-amd64.tar.gz kubectl-imageof-linux-amd64 LICENSE
tar -czf kubectl-imageof-darwin-amd64.tar.gz kubectl-imageof-darwin-amd64 LICENSE
tar -czf kubectl-imageof-darwin-arm64.tar.gz kubectl-imageof-darwin-arm64 LICENSE
zip kubectl-imageof-windows-amd64.zip kubectl-imageof-windows-amd64.exe LICENSE
```

### Create a GitHub Release (e.g., v0.1.0) and upload these archives as release assets

### Create a manifest

### Generate SHA256 with and apply inside manifest

```bash
sha256sum kubectl-imageof-*.tar.gz
sha256sum kubectl-imageof-*.zip
```

### Test locally with krew

Before publishing, test your manifest with your local krew:

```bash
kubectl krew install --manifest=imageof.yaml --archive=kubectl-imageof-linux-amd64.tar.gz
kubectl imageof -n default -A
```

If it works, you’re ready to publish.

### Submit to krew-index

1. Fork the official repo: <https://github.com/kubernetes-sigs/krew-index>
2. Add your manifest under plugins/imageof.yaml
3. Commit & create a PR to krew-index

Krew maintainers will review your manifest. Once merged, users can install with:

```bash
kubectl krew install imageof
```

And then simply use:

```bash
kubectl imageof -n default -A
```
