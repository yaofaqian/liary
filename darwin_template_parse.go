 root := "/opt/cli" // 要搜索的根目录
    pattern := "darwin_template" // 文件名前缀

    err := filepath.Walk(root, func(path string, info fileInfo, err error) error {
        if err != nil {
            return err // 返回错误停止遍历
        }
        
        // 检查是否为文件以及文件名是否匹配
        if !info.IsDir() && strings.HasPrefix(info.Name(), pattern) && strings.HasSuffix(info.Name(), ".yaml") {
            fmt.Println(path) // 输出符合条件的文件路径
        }
        return nil
    })

    if err != nil {
        fmt.Printf("遍历文件时发生错误：%v\n", err)
    }
