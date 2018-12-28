// Package logsystem 调用依据了包net/rpc
// 使用 @see https://colobu.com/2016/09/18/go-net-rpc-guide/
package logsystem

// Worker assess
type Worker struct {
	name          string
	address       string
	holdFile      string
	masterAddress string
}

// Cmd 是rpc调用时需要传入的参数
type Cmd struct {
	Command string
	Flag    string
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
