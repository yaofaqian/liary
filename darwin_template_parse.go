func GetPodEvents(podName string) (string, error) {
    // 使用 kubectl describe pod 命令
    cmd := exec.Command("kubectl", "describe", "pod", podName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("无法执行命令: %v", err)
    }

    // 查找 Events 部分，正则匹配 "Events" 及其后续内容
    re := regexp.MustCompile(`(?s)Events:\s+(.*)`)
    matches := re.FindStringSubmatch(string(output))
    if matches == nil || len(matches) < 2 {
        return "", fmt.Errorf("未找到 Events 部分")
    }

    return strings.TrimSpace(matches[1]), nil
}
