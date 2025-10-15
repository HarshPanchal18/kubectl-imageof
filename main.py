#!/usr/bin/env python3

import argparse
import os
import sys

from kubernetes import client, config
from kubernetes.client.rest import ApiException

def get_client() -> client.CoreV1Api:
    kubeconfig = os.path.expanduser("~/.kube/config")
    try:
        config.load_kube_config(config_file=kubeconfig)
    except Exception as e:
        print(f"Error loading kubeconfig: {e}", file=sys.stderr)
        sys.exit(1)
    return client.CoreV1Api()

def print_pod_images(api, namespace, pod_name, verbose) -> None:
    try:
        pod = api.read_namespaced_pod(name=pod_name, namespace=namespace)
    except ApiException as e:
        print(f"Error retrieving pod '{pod_name}' in namespace '{namespace}': {e}", file=sys.stderr)
        sys.exit(1)

    if verbose:
        print(f"Pod {pod.metadata.name}:\n  ", end="")

    for container in pod.spec.containers:
        print(f"{container.name}: {container.image}")

def print_all_pod_images(api, namespace, treeview) -> None:
    try:
        pods = api.list_namespaced_pod(namespace=namespace)
    except ApiException as e:
        print(f"Error listing pods in namespace '{namespace}': {e}", file=sys.stderr)
        sys.exit(1)

    items = pods.items
    if not items:
        print(f"No pods found in namespace {namespace}")
        return

    for pod in items:
        print(f"Pod: {pod.metadata.name}")
        for container in pod.spec.containers:
            if treeview:
                print(f"└──{container.name}: {container.image}")
            else:
                print(f"  {container.name}: {container.image}")

def usage() -> str:
    return """
Usage:
    Quickly retrieve image(s) of pod(s) instead of grepping out from the description.

Syntax:
    kubectl_imageof POD_NAME -n NAMESPACE
    kubectl_imageof -A -n NAMESPACE

Output:
    CONTAINER: IMAGE

Options:
    -A, --all                List images of all pods in the namespace
    -h, --help               Print plugin usage
    -n, --namespace string   Namespace of the pod(s) (default "default")
    -t, --tree               Show tree view for multiple pods
    -v, --verbose            Show pod name in output

Example:
    $ kubectl_imageof redis -n redis
    redis: redis

    $ kubectl_imageof redis -n redis -v
    Pod redis:
      redis: redis

    $ kubectl_imageof -A -n harbor
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
        """

def main():
    parser = argparse.ArgumentParser(add_help=False)
    parser.add_argument("pod_name", nargs="?", help="Name of the pod")
    parser.add_argument("-n", "--namespace", default="default", help="Namespace of the pod(s)")
    parser.add_argument("-A", "--all", action="store_true", help="List images of all pods in the namespace")
    parser.add_argument("-v", "--verbose", action="store_true", help="Show pod name in output")
    parser.add_argument("-t", "--tree", action="store_true", help="Show tree view for multiple pods")
    parser.add_argument("-h", "--help", action="store_true", help="Print plugin usage")

    args = parser.parse_args()

    if args.help:
        print(usage())
        return

    if not args.namespace:
        print("Error: -n NAMESPACE is required", file=sys.stderr)
        sys.exit(1)

    api = get_client()

    if args.all:
        print_all_pod_images(api, args.namespace, args.tree)
    elif args.pod_name:
        print_pod_images(api, args.namespace, args.pod_name, args.verbose)
    else:
        print(usage())

if __name__ == "__main__":
    main()