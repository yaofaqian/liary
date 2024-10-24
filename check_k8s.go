package check_k8s

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "time"
)

// CheckAPIServer checks the health of the k8s API server
func CheckAPIServer() {
    url := "https://127.0.0.1:6443/healthz"
    // Skip TLS verification for local testing
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Timeout: 5 * time.Second, Transport: tr}

    resp, err := client.Get(url)
    if err != nil {
        fmt.Printf("Failed to connect to API server: %v\n", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        fmt.Println("API Server is healthy.")
    } else {
        fmt.Printf("API Server returned status: %v\n", resp.Status)
    }
}
// CheckNodes gets the node status using kubectl command
func CheckNodes() {
    cmd := exec.Command("kubectl", "get", "nodes", "-o", "wide")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error getting nodes: %v\n", err)
        return
    }
    
    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.Contains(line, "Ready") || strings.Contains(line, "NotReady") {
            fmt.Println(line) // Output the node name and status
        }
    }
}
// CheckPods checks the status of kube-system pods and investigates abnormal ones
func CheckPods() {
    cmd := exec.Command("kubectl", "get", "pods", "-n", "kube-system")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error getting pods: %v\n", err)
        return
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.Contains(line, "Running") || strings.Contains(line, "Completed") {
            fmt.Println(line)
        } else if strings.Contains(line, "Error") || strings.Contains(line, "CrashLoopBackOff") {
            fields := strings.Fields(line)
            if len(fields) > 0 {
                podName := fields[0]
                fmt.Printf("Investigating pod %s\n", podName)
                CheckPodEvents(podName)
            }
        }
    }
}
// CheckPodEvents describes a problematic pod to get more details
func CheckPodEvents(podName string) {
    cmd := exec.Command("kubectl", "describe", "pod", podName, "-n", "kube-system")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error describing pod %s: %v\n", podName, err)
        return
    }
    fmt.Println(string(output))
}
// CheckRouting checks the system routing table for k8s routes
func CheckRouting() {
    cmd := exec.Command("route", "-n")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error checking route: %v\n", err)
        return
    }

    fmt.Println(string(output))
}
