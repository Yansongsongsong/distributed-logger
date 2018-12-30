// Package prompt 提供一个CLI界面
package prompt

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type workerAddress []string

// new一个存放命令行参数值的slice
func newWorkerAddress(vals []string, p *[]string) *workerAddress {
	*p = vals
	return (*workerAddress)(p)
}

// 实现flag包中的Value接口，将命令行接收到的值用,分隔存到slice里
func (s *workerAddress) Set(val string) error {
	*s = workerAddress(strings.Fields(val))
	return nil
}

// 实现flag包中的Value接口，将命令行接收到的值用,分隔存到slice里
func (s *workerAddress) String() string {
	// default value
	*s = workerAddress([]string{})
	// when it is cast as string, return this value
	return ""
}

var (
	Help     bool
	IsMaster bool
	IsWorker bool

	MasterAddress string
	WorkerAddress []string
	Name          string
	Filepath      string
)

func init() {
	flag.BoolVar(&Help, "h", false, "this help")
	flag.BoolVar(&IsMaster, "m", false, "to set up server for master")
	flag.BoolVar(&IsWorker, "w", false, "to set up server for worker")

	// 注意 `master address`。默认是 -ma string，有了 `master address` 之后，变为 -s master address
	flag.StringVar(&MasterAddress, "ma", "", "set the `master_address` that we can communicate with")
	flag.Var(newWorkerAddress([]string{}, &WorkerAddress), "wa", "set the `worker_address` that we can communicate with, spilt with space")
	n, _ := os.Hostname()
	flag.StringVar(&Name, "n", n, "set `name` for worker")
	flag.StringVar(&Filepath, "f", "", "set the `file` position that the log is from")

	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stdout, `logger, the distributed logger system.
Usage: logger [-?mrh] [-ma masterAddress] [-wa workerAddress] [-n name] [-f file]

Options:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stdout, `
Example:
  1. to set up master
	logger -m -ma localhost:9990 -wa "localhost:9991 localhost:9992 localhost:9993"
  2. to set up worker
	logger -w -ma localhost:9990 -wa localhost:9991 -f ./machine.1.log -n machine.i
  3. to get help
	logger -h
`)
}
