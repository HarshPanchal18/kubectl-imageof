package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
func printPodImages(client *kubernetes.Clientset, namespace, podName string) error {
    pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
    if err != nil {
        return err
    }
    fmt.Printf("Image(s) for pod %s in namespace %s:\n", pod.Name, namespace)
    for _, c := range pod.Spec.Containers {
        fmt.Printf("  %s: %s\n", c.Name, c.Image)
    }
    return nil
}

// printAllPodImages - prints images for all pods in the given namespace
func printAllPodImages(client *kubernetes.Clientset, namespace string) error {
    pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return err
    }
    if len(pods.Items) == 0 {
        fmt.Printf("No pods found in namespace %s\n", namespace)
        return nil
    }
    for _, pod := range pods.Items {
        fmt.Printf("Images for pod %s:\n", pod.Name)
        for _, c := range pod.Spec.Containers {
            fmt.Printf("  %s: %s\n", c.Name, c.Image)
        }
    }
    return nil
}

func main() {
	var namespace string
	var allPods bool

	flag.StringVar(&namespace, "n", "default", "Namespace")
	flag.StringVar(&namespace, "namespace", "default", "Namespace (required)")
    flag.BoolVar(&allPods, "A", false, "If set, list images of all pods in the namespace")
    flag.BoolVar(&allPods, "all", false, "If set, list images of all pods in the namespace")
	flag.Parse()

	client, err := getClientset()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	args := flag.Args()

	if allPods {
		// Print images for all pods in the namespace
		if err := printAllPodImages(client, namespace); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else if len(args) == 1 {
		// Print images for a single pod
		if err := printPodImages(client, namespace, args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr,
			"Usage:\n",
			"  kubectl imageof POD_NAME -n NAMESPACE\n",
			"  kubectl imageof -A -n NAMESPACE")
		os.Exit(1)
	}
}