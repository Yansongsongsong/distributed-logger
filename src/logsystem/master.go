package logsystem

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

// Master 实质是一个可显示结果的程序
// 1. 维护所有可使用的worker列表
// 		1.1 可以接收worker注册请求
//		1.2 发送心跳包 检测worker状态（可以做成第二个函数）
//		1.3 woker信息 workername + workcer address
// 		1.4 可以断线重连
// 2. 提供一个命令行字段 可以fetch字段
// 3. 有良好的输出界面
// 		3.1 "machineName lineNumber grepField"
//		3.2 提示某个机器是否执行结束
// 		3.3 机器中断&推出提示
//		3.4 可以重定向到文件
//				3.4.1 可以决定是否保留机器名和行号
//				3.4.2 文件名
type Master struct {
	// 可以也把Master做成一个rpc server
	// 1. register()
	// 2. beat()
	// 在worker中 可以轮训这个server有没有build起 利用goroutine 切出去
	// 防止死循环
}

func RunMaster() {
	worker, err := rpc.DialHTTP("tcp", "localhost:9991")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	for {
		var str string
		log.Println("\ntype pattern: ")
		bio := bufio.NewReader(os.Stdin)
		if line, _, err := bio.ReadLine(); err != nil {
			log.Fatal("input: ", err)
		} else {
			str = string(line)
		}
		//str = "grep author"
		log.Println("after typing")
		strs := strings.Fields(str)
		log.Println("cmds: ", strs)
		args := &Cmd{Command: strs[0], Flag: strs[1:]}
		rs := new(ResultSet)

		e := worker.Call("Worker.FetchResults", args, &rs)

		if e != nil {
			log.Println("Worker error: ", e)
		}

		fmt.Printf("call end\n args: %s,\n rs: %s\n", args, rs)
	}
}
