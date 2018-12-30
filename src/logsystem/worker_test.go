package logsystem

import (
	"os"
	"testing"
)

var (
	patt = "patt"
	by   = []byte("bytes")
	wr   = NewWorker("wr", "localhost:9999", "../test/Test.txt", "localhost:9990")
)

func TestCheckCacheNo(t *testing.T) {
	var data []byte
	yes, e := wr.checkCache(patt, &data)

	if e != nil {
		t.Fatal(e)
	}

	if yes {
		t.Fatal("cacheData: ", string(data))
	} else {
		t.Log("pointer: ", data)
	}

}

func TestCache(t *testing.T) {
	wr.cache(patt, &by)

	fileinfo, e := os.Stat(patt)
	if e != nil {
		t.Fatal(e)
	}
	// 文件存在
	t.Log("It is ok: \n fileinfo name: ", fileinfo.Name(), "size: ", fileinfo.Size())
}

func TestCheckCacheYes(t *testing.T) {
	var data []byte
	yes, e := wr.checkCache(patt, &data)

	if e != nil {
		t.Fatal(e)
	}

	if yes {
		t.Log("cacheData: ", string(data))
	} else {
		t.Fatal("pointer: ", data)
	}
}

func TestExecNonGrepCmd(t *testing.T) {
	var data []byte
	cmd := Cmd{"echo", []string{"string1", "string2", "string3"}}

	e := wr.execNonGrepCmd(&data, cmd)
	if e != nil {
		t.Fatal("Wrong: ", e)
	}

	t.Log("result: \n", string(data))
}

func TestExecGrepCmd(t *testing.T) {
	var data []byte
	cmd := Cmd{"grep", []string{"author"}}
	e := wr.execGrepCmd(&data, cmd)

	if e != nil {
		t.Fatal("Wrong: ", e)
	}

	t.Log("result: \n", string(data))
}
func TestClear(t *testing.T) {
	clearFile(patt)

	fileinfo, e := os.Stat(patt)
	if e != nil {
		t.Log(e)
	} else {
		// 文件存在
		t.Fatal("It is ok: \n fileinfo ", fileinfo)
	}
}
