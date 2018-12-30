package logsystem

import (
	"runtime"
	"testing"
	"time"
)

var (
	mr = NewMaster([]string{"localhost:9991", "localhost:9992", "localhost:9993"}, "localhost:9990")
)

func TestReDialHTTP(t *testing.T) {

	// mr.reDialHTTP("localhost:9991")
	// mr.reDialHTTP("localhost:9992")
	// mr.reDialHTTP("localhost:9993")
	// runtime.Gosched()
	// t.Log("sleep..")
	// time.Sleep(time.Duration(reDialTimes*reDialDuration+10) * time.Second)
	// t.Log("sleep ends")
	// t.Log("mr.initWorkerSet: ", mr.initWorkerSet)
	// t.Log("mr.workerMap: ", mr.workerMap)
}

func TestBeats(t *testing.T) {
	mr.reDialHTTP("localhost:9991")
	mr.reDialHTTP("localhost:9992")
	mr.reDialHTTP("localhost:9993")
	mr.beats()
	runtime.Gosched()
	t.Log("sleep..")
	time.Sleep(time.Duration(200) * time.Second)
	t.Log("sleep ends")
	t.Log("mr.initWorkerSet: ", mr.initWorkerSet)
	t.Log("mr.workerMap: ", mr.workerMap)
}
