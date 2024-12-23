#!/bin/bash

# 定义当前时间
CURRENT_TIME=$(date "+%Y-%m-%d-%H-%M-%S")
TOP_LOG="top_log_$CURRENT_TIME.txt"
IO_LOG="io_log_$CURRENT_TIME.txt"

# 获取一次 top 信息，持续 1 秒
top -b -n 1 > "$TOP_LOG"

# 清空 IO 日志文件
> "$IO_LOG"

# 获取当前时间并格式化
echo "Log Time: $CURRENT_TIME" >> "$IO_LOG"

# 将 top_log 中的每一行逐行读取
while IFS= read -r line; do
    # 从 line 中提取 PID 和进程名
    PID=$(echo "$line" | awk '{print $1}')
    PROCESS_NAME=$(echo "$line" | awk '{print $12}')

    # 跳过第一行和空行
    if [[ "$line" == *"PID"* || -z "$line" ]]; then
        continue
    fi

    # 检查 PID 是否存在
    if [[ -d "/proc/$PID" ]]; then
        # 从 /proc/PID/io 获取 I/O 使用情况
        IO_INFO=$(cat /proc/$PID/io)

        # 输出格式化信息
        echo "ProcessName: $PROCESS_NAME" >> "$IO_LOG"
        echo "PID: $PID # 进程id" >> "$IO_LOG"
        
        # 解析 I/O 信息并添加中文注释
        while IFS= read -r io_line; do
            case $io_line in
                rchar*) echo "$io_line # 读取的字符数" >> "$IO_LOG" ;;
                wchar*) echo "$io_line # 写入的字符数" >> "$IO_LOG" ;;
                syscr*) echo "$io_line # 读取的系统调用次数" >> "$IO_LOG" ;;
                syscw*) echo "$io_line # 写入的系统调用次数" >> "$IO_LOG" ;;
                read_bytes*) echo "$io_line # 读取的字节数" >> "$IO_LOG" ;;
                write_bytes*) echo "$io_line # 写入的字节数" >> "$IO_LOG" ;;
                cancel_bytes*) echo "$io_line # 取消的字节数" >> "$IO_LOG" ;;
                *) echo "$io_line" >> "$IO_LOG" ;;  # 其他字段直接输出
            esac
        done <<< "$IO_INFO"

        echo "-----------------------------" >> "$IO_LOG"
    fi
done < "$TOP_LOG"

echo "I/O information logged to $IO_LOG."
