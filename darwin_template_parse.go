// ExtractDecryptionResult 从文本中提取解密结果
func ExtractDecryptionResult(input string) ([]string, error) {
    // 使用正则表达式匹配 "解密结果：“任意长度字符串”"
    re := regexp.MustCompile(`解密结果：“([^”]*)”`)

    matches := re.FindAllStringSubmatch(input, -1)
    if matches == nil {
        return nil, fmt.Errorf("没有找到匹配的解密结果")
    }

    var results []string
    for _, match := range matches {
        if len(match) > 1 {
            results = append(results, match[1]) // 添加匹配到的内容
        }
    }
    return results, nil
}
