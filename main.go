package main

import (
    "fmt"
    "TQServerClusterTools/check_k8s"
)

func main() {
    fmt.Println("Starting cluster check...")

    // Check k8s API server
    check_k8s.CheckAPIServer()

    // Check k8s nodes
    check_k8s.CheckNodes()

    // Check kube-system pods
    check_k8s.CheckPods()

    // Check system routing
    check_k8s.CheckRouting()

    fmt.Println("Cluster check complete.")
}
