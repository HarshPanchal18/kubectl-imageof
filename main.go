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

func main() {
	var namespace string
	flag.StringVar(&namespace, "n", "default", "Namespace")
	flag.Parse()

	// args := flag.Args()
	// if(len(args)) < 1 {
	// 	fmt.Fprintf(os.Stderr, "Usage: kubectl-imageof POD_NAME [flags]\n")
	// 	os.Exit(1)
	// }
	// podName := args[0]

	// Load kubeconfig
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// List all pods in the namespace
    pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    if len(pods.Items) == 0 {
        fmt.Printf("No pods found in namespace %s\n", namespace)
        return
    }

    for _, pod := range pods.Items {
        fmt.Printf("Image(s) for pod %s:\n", pod.Name)
        for _, container := range pod.Spec.Containers {
            fmt.Printf("  %s: %s\n", container.Name, container.Image)
        }
    }
}