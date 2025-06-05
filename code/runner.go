package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	// 内存性能分析相关 ========================
	// 创建内存统计对象，用于记录初始内存状态
	var bm runtime.MemStats
	// 读取当前内存分配情况到bm对象
	runtime.ReadMemStats(&bm)
	// 输出程序启动时的内存使用量（单位KB）
	// Alloc表示已申请且仍在使用的堆内存字节数
	fmt.Printf("KB: %v\n", bm.Alloc/1024)

	// 时间记录相关 ==========================
	// 获取当前时间戳（毫秒级），用于后续性能统计
	now := time.Now().UnixMilli()
	println("当前时间(毫秒) ==> ", now)

	// 子进程执行相关 ========================
	// 创建执行go run命令的进程对象
	// 该命令会运行用户提交的代码文件（假设是加法计算程序）
	cmd := exec.Command("go", "run", "code-user/main.go")

	// 创建缓冲区用于捕获标准输出和错误输出
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr // 重定向标准错误到stderr缓冲区
	cmd.Stdout = &out    // 重定向标准输出到out缓冲区

	// 获取标准输入管道（用于向子程序传递输入参数）
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	// 通过管道向子进程传递测试用例输入
	// 注意：此处未处理可能的写入错误，实际生产环境需要添加错误处理
	_, _ = io.WriteString(stdinPipe, "23 11\n") // 模拟输入两个数字

	// 执行命令并等待完成
	if err := cmd.Run(); err != nil {
		// 如果执行出错，输出错误信息和标准错误内容
		log.Fatalln(err, stderr.String())
	}

	// 输出调试信息 ==========================
	println("Err:", string(stderr.Bytes())) // 打印子进程的错误输出
	fmt.Println(out.String())               // 打印子进程的标准输出

	// 验证结果是否正确（预期输出34，即23+11的结果）
	println(out.String() == "34\n")

	// 内存性能统计结束 ========================
	// 记录程序执行后的内存状态
	var em runtime.MemStats
	runtime.ReadMemStats(&em)
	// 输出执行后的内存使用量（单位KB）
	fmt.Printf("KB: %v\n", em.Alloc/1024)

	// 时间统计结束 ==========================
	end := time.Now().UnixMilli()
	println("当前时间 ==> ", end)
	println("耗时 ==> ", end-now) // 计算并输出总执行耗时
}
