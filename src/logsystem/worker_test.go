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

func TestProcessBytes(t *testing.T) {
	var data []byte
	cmd := Cmd{"grep", []string{"author"}}
	e := wr.execGrepCmd(&data, cmd)

	if e != nil {
		t.Fatal("Wrong: ", e)
	}

	rs := wr.processBytes(data)

	t.Log("resultset: \n", rs)
}

func TestFetchResults(t *testing.T) {
	var bytes []byte

	if yes, _ := wr.checkCache("author", &bytes); yes {
		t.Fatal("wrong")
	}

	cmd := Cmd{"grep", []string{"author"}}
	rs, e := wr.FetchResults(&cmd)

	if e != nil {
		t.Fatal("Wrong: ", e)
	}

	if yes, _ := wr.checkCache("author", &bytes); !yes {
		t.Fatal("wrong")
	}

	t.Log("resultset: \n", rs)
}

func TestFetchResultsWithCache(t *testing.T) {
	cmd := Cmd{"grep", []string{"author"}}
	rs, e := wr.FetchResults(&cmd)

	if e != nil {
		t.Fatal("Wrong: ", e)
	}

	t.Log("resultset: \n", rs)
}

func TestFetchResultsWithOtherCmd(t *testing.T) {
	cmd := Cmd{"echo", []string{"string1", "string2", "string3"}}
	rs, e := wr.FetchResults(&cmd)

	if e != nil {
		t.Fatal("Wrong: ", e)
	}

	t.Log("resultset: \n", rs)
}

func TestClearAllCache(t *testing.T) {
	t.Log("All files: ", wr.cacheFile)
	wr.ClearAllCache()
	if wr.cacheFile != nil {
		t.Fatal("wrong: \ncacheFile poniter: ", wr.cacheFile)
	}
}

func TestCheckFile(t *testing.T) {
	checkFile("master.go")
	checkFile("../.gitignore")
	e := checkFile("../main")
	checkFile("askjbkasjcb12,1le13")
	if e != nil {
		if e.Error() != "directory" {
			t.Fatal(e)
		}
	}
}

func TestClearFile(t *testing.T) {
	clearFile("askjbkasjcb12,1le13")
	clearFile("../.gitignore")

}
