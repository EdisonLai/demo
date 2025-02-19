package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	criapi "k8s.io/cri-api/pkg/apis/runtime/v1" // 使用 v1 版本
)

const (
	containerdSocket = "/run/containerd/containerd.sock"
	crioSocket       = "/var/run/crio/crio.sock"
)

func main() {
	conn, err := grpc.Dial(
		getCRISocket(),
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.Dial("unix", addr)
		}),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect: %v", err))
	}
	defer conn.Close()

	client := criapi.NewRuntimeServiceClient(conn)

	pods, err := client.ListPodSandbox(context.Background(), &criapi.ListPodSandboxRequest{})
	if err != nil {
		panic(fmt.Sprintf("ListPodSandbox failed: %v", err))
	}

	for _, pod := range pods.Items {
		status, err := client.PodSandboxStatus(context.Background(), &criapi.PodSandboxStatusRequest{
			PodSandboxId: pod.Id,
		})
		if err != nil {
			fmt.Printf("Failed to get pod status: %v\n", err)
			continue
		}

		// 直接从 status.Info 中解析 PID
		info := status.GetInfo()
		if info == nil {
			continue
		}

		// 解析 info 中的 "pid" 字段（JSON 格式）
		pid := parsePIDFromInfo(info["info"])
		fmt.Printf("Pod: %s, Pause容器PID: %d\n", pod.Metadata.Name, pid)
	}
}

// 从 info JSON 字符串中解析 PID
func parsePIDFromInfo(infoStr string) int {
	// 这里需要实际解析 JSON 字符串，例如：
	// {"pid": 12345, ...}
	// 为简化示例，假设直接提取 PID
	// 实际代码中应使用 json.Unmarshal
	var pid int
	fmt.Sscanf(infoStr, `{"pid":%d`, &pid) // 简化处理，实际需要完整解析
	return pid
}

// 自动检测 CRI Socket
func getCRISocket() string {
	for _, path := range []string{containerdSocket, crioSocket} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	panic("未找到可用的 CRI Socket")
}
