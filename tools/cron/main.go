// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// 定义要执行的定时任务函数
func myTask() {
	fmt.Printf("定时任务执行：%s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	// 1. 创建 cron 实例（默认支持秒级表达式，若需标准 5 字段，需用 WithSeconds(false) 禁用）
	c := cron.New(cron.WithSeconds())

	// 2. 添加定时任务（cron 表达式：每 5 秒执行一次）
	// AddFunc 返回任务 ID，用于后续管理（如删除任务）
	taskID, err := c.AddFunc("*/1 * * * * *", myTask)
	if err != nil {
		log.Fatalf("添加任务失败：%v", err)
	}
	fmt.Printf("已添加任务，ID：%d\n", taskID)

	// 3. 启动 cron 调度器（非阻塞，会启动一个后台协程）
	c.Start()
	defer c.Stop() // 程序退出时停止调度器

	// 防止程序立即退出（模拟业务逻辑阻塞）
	select {}
}
