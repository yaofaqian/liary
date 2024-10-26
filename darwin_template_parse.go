package darwin_template_parse
func DarwinTemplateParse(filePath string){
// 定义查找的目录和文件模式
    dir := "/opt/cli/plan_client"
    pattern := "darwin_template_*.yaml"

    // 检查目录是否存在
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        fmt.Printf("错误：目录 %s 不存在。\n", dir)
        return
    }

    // 使用 filepath.Glob 查找匹配文件
    matches, err := filepath.Glob(filepath.Join(dir, pattern))
    if err != nil {
        fmt.Printf("错误：查找文件时出现问题，原因是 %v\n", err)
        return
    }

    // 检查是否有匹配的文件
    if len(matches) == 0 {
        fmt.Println("没有找到符合条件的文件。")
        return
    }

    // 提取文件名而非完整路径
    fileName := matches[0]
    fmt.Printf("找到第一个符合条件的文件：%s\n", filepath.Base(fileName))
  }
