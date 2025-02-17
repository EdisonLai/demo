package main

import (
	"context"
	"cri-demo/global"
	"fmt"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

func main() {
	args := os.Args
	_, err := sandboxFind(args[1])
	fmt.Println(err)
}

func containerdFind() (int, error) {
	client, err := containerd.New(global.ContainerdRuntimeAddress)
	if err != nil {
		return 0, fmt.Errorf("error creating containerd client:%v", err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), global.K8sNamespace)
	containers, err := client.Containers(ctx)
	if err != nil {
		return 0, fmt.Errorf("error listing containers in containerd: %v", err)
	}

	for _, ctr := range containers {
		labels, err := ctr.Labels(ctx)
		if err != nil {
			return 0, fmt.Errorf("error getting labels for container %s: %v", ctr.ID(), err)
		}
		fmt.Printf("container %s labels: %v\n", ctr.ID(), labels)
	}
	return 0, nil
}

func sandboxFind(podName string) (int, error) {
	client, err := containerd.New(global.ContainerdRuntimeAddress)
	if err != nil {
		return 0, fmt.Errorf("error creating containerd client:%v", err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), global.K8sNamespace)
	sandbox, err := client.LoadSandbox(ctx, podName)
	if err != nil {
		return 0, fmt.Errorf("error listing containers in containerd: %v", err)
	}

	fmt.Printf("sandbox %+v\n", sandbox)
	return 0, nil
}
