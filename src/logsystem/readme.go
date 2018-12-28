// Package logsystem 中程序分为两个角色 1. master 2. worker
// 与传统的mapreduce程序相比，这个程序缺少了spilt 与 schedule部分。
// 1. master
// 		master的作用是
//		1. 提供一个client 供人查看
//		2. 维护和所有worker之间的连接
// 2. worker
// 		worker的作用是
//		1. worker负责执行grep，对单机上对某个特定文件进行grep操作
//		2. 执行被master调度过来的reduce操作
package logsystem
