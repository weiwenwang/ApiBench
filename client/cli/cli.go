package cli

import (
	"sync"
)

type Cli struct {
	Pattern string
	Port    int
}


var ins *Cli
var once sync.Once

func GetCli() *Cli {
	once.Do(func() {
		ins = &Cli{}
	})
	return ins
}
