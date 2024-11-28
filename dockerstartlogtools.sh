#!/bin/bash

# 定义输出目录
output_dir="docker_start_log"
zip_file="${output_dir}.zip"

# 清理旧文件和文件夹
if [ -d "$output_dir" ]; then
    rm -rf "$output_dir"
fi
if [ -f "$zip_file" ]; then
    rm -f "$zip_file"
fi

# 创建输出目录
mkdir -p "$output_dir"

# 获取所有容器ID
container_ids=$(docker ps -q)

# 检查是否有容器运行
if [ -z "$container_ids" ]; then
    echo "没有运行的容器，脚本退出。"
    exit 1
fi

# 循环处理每个容器ID
for container_id in $container_ids; do
    # 输出日志到对应的txt文件
    log_file="${output_dir}/${container_id}.txt"
    echo "正在处理容器ID: $container_id"
    journalctl | grep "$container_id" > "$log_file"
done

# 压缩输出目录为zip文件
zip -r "$zip_file" "$output_dir" > /dev/null 2>&1

# 检查压缩是否成功
if [ $? -eq 0 ]; then
    echo "日志已成功压缩为 $zip_file"
else
    echo "日志压缩失败，请检查。"
fi
