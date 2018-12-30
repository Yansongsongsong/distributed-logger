// Package main 提供一个CLI界面
package main

import (
	"flag"
	"logsystem"
	"prompt"
)

var (
	help     = &prompt.Help
	isMaster = &prompt.IsMaster
	isWorker = &prompt.IsWorker

	masterAddress = &prompt.MasterAddress
	workerAddress = &prompt.WorkerAddress
	name          = &prompt.Name
	filepath      = &prompt.Filepath
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
	}

	if !*isMaster && !*isWorker {
		flag.Usage()
	}

	if *isMaster {
		logsystem.RunMaster(*masterAddress, (*workerAddress))

	}

	if *isWorker {
		logsystem.RunWorker(*name, *masterAddress, (*workerAddress)[0], *filepath)
	}

}
