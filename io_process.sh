#!/bin/bash

# 定义日志文件
TOP_LOG="top_log.txt"
IO_LOG="io_log.txt"

# 每隔 5 秒获取一次 top 信息，持续 60 秒
top -b -n 12 -d 5 > "$TOP_LOG"

# 清空 IO 日志文件
> "$IO_LOG"

# 遍历 top 日志文件中的 PID
while read -r line; do
    # 从 line 中提取 PID
    # 假设 top 输出中，PID 在第二列（根据具体情况可能需要调整）
    PID=$(echo "$line" | awk '{print $1}')

    # 跳过第一行和空行
    if [[ "$line" == *"PID"* || -z "$line" ]]; then
        continue
    fi

    # 检查 PID 是否存在
    if [[ -d "/proc/$PID" ]]; then
        # 从 /proc/PID/io 获取 I/O 使用情况
        IO_INFO=$(cat /proc/$PID/io)
        echo "PID: $PID" >> "$IO_LOG"
        echo "$IO_INFO" >> "$IO_LOG"
        echo "-----------------------------" >> "$IO_LOG"
    fi
done < <(grep -E '^[0-9]' "$TOP_LOG")

echo "I/O information logged to $IO_LOG."
