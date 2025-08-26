package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// getClientset - loads kubeconfig and returns a Kubernetes clientset
func getClientset() (*kubernetes.Clientset, error) {
    kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        return nil, err
    }
    return kubernetes.NewForConfig(config)
}

// printPodImages - prints images for a single pod in the given namespace
func printPodImages(client *kubernetes.Clientset, namespace, podName string, verbose bool) error {
    pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
    if err != nil { return err }

	if verbose { fmt.Printf("Pod %s:\n  ", pod.Name) }

	for _, c := range pod.Spec.Containers {
        fmt.Printf("%s: %s\n", c.Name, c.Image)
    }
    return nil
}

// printAllPodImages - prints images for all pods in the given namespace
func printAllPodImages(client *kubernetes.Clientset, namespace string, treeview bool) error {
    pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil { return err }

	// No pods
    if len(pods.Items) == 0 {
        fmt.Printf("No pods found in namespace %s\n", namespace)
        return nil
    }

	for _, pod := range pods.Items {
		fmt.Printf("Pod: %s\n", pod.Name)

		for _, container := range pod.Spec.Containers {
			if treeview {
				fmt.Printf("└──%s: %s\n", container.Name, container.Image)
			} else {
				fmt.Printf("  %s: %s\n", container.Name, container.Image)
			}
        }
    }
    return nil
}

func printUsage() {
	fmt.Fprintln(os.Stderr,
			`
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
		trivy: goharbor/trivy-adapter-photon:v2.13.2`)
}

func main() {
	var help bool
	var namespace string
	var allPods bool
	var verbose bool
	var tree bool

    pflag.BoolVarP(&help, "help", "h", false, "Print plugin usage")
	pflag.StringVarP(&namespace, "namespace", "n", "default", "Namespace of the pod(s)")
    pflag.BoolVarP(&allPods, "all", "A", false, "List images of all pods in the namespace")
    pflag.BoolVarP(&verbose, "verbose", "v", false, "Show pod name in output")
    pflag.BoolVarP(&tree, "tree", "t", false, "Show tree view for multiple pods")
	pflag.Parse()

	if help {
		printUsage()
		return
	}

	if namespace == "" {
        fmt.Fprintln(os.Stderr, "Error: -n NAMESPACE is required")
        os.Exit(1)
    }

	client, err := getClientset()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	args := pflag.Args()

	if allPods { // Print images for all pods in the namespace
		if err := printAllPodImages(client, namespace, tree); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else if len(args) == 1 { // Print images for a single pod
		pod := args[0]
		if err := printPodImages(client, namespace, pod, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		printUsage()
		return
	}
}