package logsystem

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Client 是别名
type Client = rpc.Client

type workerAddr = string

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
	initWorkerSet map[workerAddr]bool
	masterAddress string
	workerMap     map[workerAddr]*Client
}

const (
	// 重连次数
	reDialTimes int = 10
	// 5s 重连
	reDialDuration int = 5
	// 10s beats 心跳
	beatsDuration int = 10
)

var (
	mutex sync.Mutex
)

// todo 有bug 可能存在同时写的错误
func (mr *Master) reDialHTTP(addr string) {
	go func() {
		// 重连次数
		for index := 0; index < reDialTimes; index++ {
			worker, err := rpc.DialHTTP("tcp", addr)
			if err != nil {
				// 重连失败
				log.Printf("When dialing to '%s', happen: %s\n", addr, err)
				// todo 提示 输入grep
				runtime.Gosched()
				time.Sleep(time.Duration(reDialDuration) * time.Second)
				continue
			}
			// 重连成功
			mutex.Lock()
			mr.workerMap[addr] = worker
			log.Println("The connected worker list is: ", mr.workerMap)
			log.Println("The worker list that master holds: ", mr.initWorkerSet)
			mutex.Unlock()
			return
		}
		// 重连过多 移除worker
		mutex.Lock()
		delete(mr.initWorkerSet, addr)
		log.Println("The connected worker list is: ", mr.workerMap)
		log.Println("The worker list that master holds: ", mr.initWorkerSet)
		mutex.Unlock()

	}()
}

func (mr *Master) beats() {
	var workerMapLock sync.Mutex
	go func() {
		for {
			log.Println("Beats begin!")
			for addr, worker := range mr.workerMap {
				args := new(BeatsArg)
				args.MasterAddr = mr.masterAddress
				res := new(BeatsRes)
				e := worker.Call("Worker.Beats", args, &res)
				if e != nil {
					log.Printf("Beats '%s', happen: %s", addr, e)
					// beats 出错
					mr.reDialHTTP(addr)
					workerMapLock.Lock()
					delete(mr.workerMap, addr)
					log.Println("The connected worker list is: ", mr.workerMap)
					log.Println("The worker list that master holds: ", mr.initWorkerSet)
					workerMapLock.Unlock()
					continue
				}
				// beats 成功返回
				log.Println("Beats! ", *res)
			}
			time.Sleep(time.Duration(beatsDuration) * time.Second)
		}
	}()

}

func (mr *Master) outputFormat(args *Cmd, rs *ResultSet) {
	fmt.Printf("'%s' called end\n", rs.WorkerName)
	for _, line := range rs.Lines {
		fmt.Printf("machine: %s, Line: %s, content: %s\n", rs.WorkerName, line.Line, line.S)
	}
}

func (mr *Master) distributedFetch(args *Cmd, rs *ResultSet) {
	for _, worker := range mr.workerMap {
		e := worker.Call("Worker.FetchResults", args, rs)

		if e != nil {
			log.Println("Worker error: ", e)
			continue
		}

		mr.outputFormat(args, rs)

	}

}

// 将输入的字符串变成Cmd
func toPrompt() *Cmd {
	// string type in
	var str string
	log.Println("\nPlease type the grep pattern: ")
	// reading
	bio := bufio.NewReader(os.Stdin)
	if line, _, err := bio.ReadLine(); err != nil {
		log.Fatal("input: ", err)
	} else {
		str = string(line)
	}

	strs := strings.Fields(str)
	log.Println("What you just input is: ", strs)
	if len(strs) == 0 {
		log.Println("Please re-input")
		return toPrompt()
	}
	args := &Cmd{Command: strs[0], Flag: strs[1:]}
	return args
}

// RunMaster 提供了prompt
func RunMaster(mrAddr string, wrAddrs []string) {
	mr := NewMaster(wrAddrs, mrAddr)

	for wradd := range mr.initWorkerSet {
		mr.reDialHTTP(wradd)
	}
	mr.beats()

	for {
		arg := toPrompt()
		rs := new(ResultSet)
		mr.distributedFetch(arg, rs)
	}
}

// NewMaster 工厂方法
func NewMaster(workerAddrs []string, masterAddress string) *Master {
	mr := new(Master)
	mr.initWorkerSet = make(map[string]bool)
	for _, v := range workerAddrs {
		mr.initWorkerSet[v] = true
	}
	mr.masterAddress = masterAddress
	mr.workerMap = make(map[string]*Client)

	return mr
}
