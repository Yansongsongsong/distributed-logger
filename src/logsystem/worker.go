// Package logsystem 调用依据了包net/rpc
// 使用 @see https://colobu.com/2016/09/18/go-net-rpc-guide/
package logsystem

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
)

type pattern = string
type cFName = string

// Worker 应该
// 1. 接收方法调用，返回fetch结果
// 2. 缓存信息 提高速率
// 3. 识别中断 删除缓存
// 4. 保持心跳 提供生存信息
// 5. 发送注册请求
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
	line string
	s    string
}

// ResultSet 是rpc调用时返回的结果
type ResultSet struct {
	// 机器名
	WorkerName string
	Lines      []result
}

func (wr *Worker) execNonGrepCmd(result *[]byte, theCmd Cmd) (err error) {
	var cmd *exec.Cmd

	// 执行非grep命令
	// 执行连续的shell命令时, 需要注意指定执行路径和参数, 否则运行出错
	cmd = exec.Command(theCmd.Command, theCmd.Flag...)
	if *result, err = cmd.Output(); err != nil {
		return err
	}

	return nil
}

func (wr *Worker) execGrepCmd(result *[]byte, theCmd Cmd) (err error) {
	c1 := exec.Command("cat", wr.holdFile)
	// 显示行号
	flags := append([]string{"-n"}, theCmd.Flag...)
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
	*result = bf.Bytes()
	return nil
}

func clearFile(name string) {
	_, e := os.Stat(name)
	if e != nil {
		// 获取文件状态出错
		if !os.IsNotExist(e) {
			// 打印其他错误
			log.Println(e)
		}
		// 忽略不存在文件的错误
	}
	// 文件存在
	err := os.Remove(name)
	if err != nil {
		// 打印其他错误
		log.Println(err)
	}
}

// ClearAllCache will clear all cache file
func (wr *Worker) clearAllCache() {
	for _, v := range wr.cacheFile {
		clearFile(v)
		log.Println("clear: ", v)
	}
	wr.cacheFile = nil
}

func (wr *Worker) cache(patt string, bytes *[]byte) {
	// 以pattern为文件名缓存文件
	err := ioutil.WriteFile(patt, *bytes, 0666)
	if err != nil {
		log.Println(err)
		// 出错此时清除文件
		clearFile(patt)
	}
	wr.cacheFile[patt] = patt
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

func (wr *Worker) processBytes(bytes []byte) (rs *ResultSet) {
	lines := strings.Split(string(bytes), "\n")
	results := []result{}
	for _, line := range lines {
		i := strings.Index(line, ":")
		if i < 0 {
			continue
		}
		results = append(results, result{line[:i], line[i+1:]})
	}
	rs = &ResultSet{wr.name, results}
	return
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
func (wr *Worker) FetchResults(cmd *Cmd, rs *ResultSet) (err error) {
	var tempResult []byte
	// 检查命令
	if strings.ToLower(cmd.Command) != "grep" {
		// 1. 如果cmd不是grep 正常执行 并返回字符的结果 不缓存
		if err = wr.execNonGrepCmd(&tempResult, *cmd); err != nil {
			log.Println(err)
			return err
		}
		r := result{line: "-1", s: string(tempResult)}
		rs = &ResultSet{WorkerName: wr.name, Lines: []result{r}}
		return nil
	}

	// 2. 如果是grep 则执行
	// 查看是否缓存
	var patt string
	for _, v := range cmd.Flag {
		patt += v
	}

	// 查看是否缓存
	yes, e := wr.checkCache(patt, &tempResult)
	if e != nil {
		log.Println(e)
	}
	if !yes {
		// 2.2 缓存未命中
		// 2.2.1 cat file | grep pattern 返回结果
		if err = wr.execGrepCmd(&tempResult, *cmd); err != nil {
			log.Println(err)
			return err
		}
		// 2.2.2 缓存结果
		wr.cache(patt, &tempResult)
	}
	// 2.1 缓存命中，返回缓存文件内容
	// 2.2.1 未命中但取得结果
	// 处理临时的 []bytes
	rs = wr.processBytes(tempResult)

	return nil
}

func checkFile(filepath string) (err error) {
	fileinfo, e := os.Stat(filepath)
	if e != nil {
		// 文件不存在
		log.Println("File is not found. It will be created.")
		newFile, err := os.Create(filepath)
		if err != nil {
			log.Println("Something wrong when create new file")
			log.Println(err)
			return err
		}
		fileinfo, _ = newFile.Stat()
	}
	if fileinfo.IsDir() {
		log.Println("It can't be one directory!")
		return errors.New("directory")
	}
	log.Printf("The file %s will be held by worker, size: %d", fileinfo.Name(), fileinfo.Size())
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
	name string,
	masterAddress string,
	workerAddress string,
	fetchFilePath string,
) {
	if e := checkFile(fetchFilePath); e != nil {
		os.Exit(1)
	}

	wr := NewWorker(name, workerAddress, fetchFilePath, masterAddress)
	rpc.Register(wr)
	rpc.HandleHTTP()
	index := strings.Index(workerAddress, ":")
	if index < 0 {
		log.Fatal("Wrong net address")
		os.Exit(1)
	}
	host := workerAddress[index:]
	log.Println("Worker' host: ", host)
	l, e := net.Listen("tcp", host)
	if e != nil {
		log.Fatal("listen error:", e)
		os.Exit(1)
	}
	log.Println("Worker is wating for master connecting...")
	http.Serve(l, nil)
}

// 初始化log
func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

// NewWorker 工厂方法
func NewWorker(name string,
	address string,
	holdFile string,
	masterAddress string) *Worker {
	wr := new(Worker)
	wr.name = name
	wr.address = address
	wr.holdFile = holdFile
	wr.masterAddress = masterAddress
	wr.cacheFile = make(map[string]string)

	return wr
}
