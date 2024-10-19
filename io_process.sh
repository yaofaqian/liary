#!/bin/bash

# 定义日志文件
TOP_LOG="top_log.txt"
CURRENT_TIME=$(date "+%Y-%m-%d-%H-%M-%S")
IO_LOG="io_log_$CURRENT_TIME.txt"

# 每隔 5 秒获取一次 top 信息，持续 60 秒
top -b -n 12 -d 5 > "$TOP_LOG"

# 清空 IO 日志文件
> "$IO_LOG"

# 获取当前时间并格式化
echo "Log Time: $CURRENT_TIME" >> "$IO_LOG"

# 声明一个数组来存储进程的 I/O 信息
declare -A io_data

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

        # 获取 read_bytes 和 write_bytes
        READ_BYTES=$(echo "$IO_INFO" | grep "read_bytes" | awk '{print $2}')
        WRITE_BYTES=$(echo "$IO_INFO" | grep "write_bytes" | awk '{print $2}')

        # 如果没有读取到值，默认设为 0
        READ_BYTES=${READ_BYTES:-0}
        WRITE_BYTES=${WRITE_BYTES:-0}
        TOTAL_IO=$((READ_BYTES + WRITE_BYTES))

        # 存储 I/O 信息到数组
        io_data["$PID"]="$PROCESS_NAME $READ_BYTES $WRITE_BYTES $TOTAL_IO"

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

# 按 read_bytes + write_bytes 对 I/O 信息进行排序
{
    echo "ProcessName # PID # 读取的字节数 # 写入的字节数 # 总 I/O 字节数"
    for key in "${!io_data[@]}"; do
        echo "${io_data[$key]} # $key"
    done
} | sort -k4 -n -r >> "$IO_LOG"

echo "I/O information logged to $IO_LOG."
