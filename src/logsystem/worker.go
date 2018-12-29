// Package logsystem 调用依据了包net/rpc
// 使用 @see https://colobu.com/2016/09/18/go-net-rpc-guide/
package logsystem

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type pattern = string
type cFName = string

// Worker has
type Worker struct {
	name          string
	address       string
	holdFile      string
	masterAddress string
	cacheFile     map[pattern]cFName
}

// Cmd 是rpc调用时需要传入的参数
type Cmd struct {
	Command string
	Flag    []string
}

// 行号 + 结果
type result struct {
	line int
	s    string
}

// ResultSet 是rpc调用时返回的结果
type ResultSet struct {
	// 机器名
	WorkerName string
	Lines      []result
}

func (wr *Worker) execNonGrepCmd(result []byte, todo Cmd) (err error) {
	var cmd *exec.Cmd

	// 执行非grep命令
	// 执行连续的shell命令时, 需要注意指定执行路径和参数, 否则运行出错
	cmd = exec.Command(todo.Command, todo.Flag...)
	if result, err = cmd.Output(); err != nil {
		return err
	}

	return nil
}

func (wr *Worker) execGrepCmd(result []byte, todo Cmd) (err error) {
	c1 := exec.Command("cat", wr.holdFile)
	// 显示行号
	flags := append([]string{"-n"}, todo.Flag...)
	c2 := exec.Command("grep", flags...)
	// 使用shell管道
	c2.Stdin, err = c1.StdoutPipe()
	if err != nil {
		return
	}
	var bf bytes.Buffer
	c2.Stdout = &bf
	err = c2.Start()
	if err != nil {
		return
	}
	err = c1.Run()
	if err != nil {
		return
	}
	err = c2.Wait()
	if err != nil {
		return
	}
	result = bf.Bytes()
	return nil
}

func clearFile(name string) {
	_, e := os.Stat(name)
	if e != nil {
		// 获取文件状态出错
		if !os.IsNotExist(e) {
			// 打印其他错误
			log.Fatalln(e)
		}
		// 忽略不存在文件的错误
	}
	// 文件存在
	err := os.Remove(name)
	if err != nil {
		// 打印其他错误
		log.Fatalln(err)
	}
}

func (wr *Worker) cache(patt string, bytes *[]byte) {
	// 以pattern为文件名缓存文件
	err := ioutil.WriteFile(patt, *bytes, 0666)
	if err != nil {
		log.Fatalln(err)
		// 出错此时清除文件
		clearFile(patt)
	}
}

func (wr *Worker) checkCache(patt string, bytes *[]byte) (bool, error) {
	p := pattern(patt)
	if _, ok := wr.cacheFile[p]; !ok {
		// 缓存未命中
		return false, nil
	}

	// 缓存命中
	file, err := os.Open(wr.cacheFile[p])
	defer file.Close()

	if err != nil {
		return false, err
	}

	var e error
	*bytes, e = ioutil.ReadAll(file)
	if e != nil {
		return false, e
	}

	return true, nil
}

// FetchResults 是Rpc调用方法
// 1. 如果cmd不是grep 则执行
// 		1.1 正常执行 并返回字符的结果 不缓存
// 2. 如果是grep 则执行
// 		查看是否缓存，
//			2.1 是，返回缓存文件内容
//		  2.2 否
// 		   	2.2.1 cat file | grep pattern 返回结果
// 				2.2.2 缓存结果
func (wr *Worker) FetchResults(cmd *Cmd, rs *ResultSet) error {
	var tempResult []byte
	var err error
	// 检查命令
	if strings.ToLower(cmd.Command) != "grep" {
		// 1. 如果cmd不是grep 正常执行 并返回字符的结果 不缓存
		if err = wr.execNonGrepCmd(tempResult, *cmd); err != nil {
			log.Fatal(err)
			return err
		}
		r := result{line: -1, s: string(tempResult)}
		rs = &ResultSet{WorkerName: wr.name, Lines: []result{r}}
		return nil
	}

	// 2. 如果是grep 则执行
	// 查看是否缓存
	if err = wr.execGrepCmd(tempResult, *cmd); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// RunWorker 在初始化worker时被调用
// 需要执行的责任
// 1. 建立一个httpserver 提供rpc
// 2. 保持rpc server的稳健不退出
// 3. 能够提供rpc功能：
//		3.1 接收方法调用，返回fetch结果
//		3.2 缓存信息 提高速率
//		3.3 识别中断 删除缓存
//		3.4 保持心跳 提供生存信息
//    3.5 发送注册请求
func RunWorker(
	masterAddress string,
	workerAddress string,
	fetchFilePath string, // todo: 待resolve
) {

}

// 初始化log
func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}
