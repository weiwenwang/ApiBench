package ser

import (
	"sync"
	"github.com/weiwenwang/ApiBench/common"
)

type Ser struct {
	Pattern string
	List    []common.Request
	UndoList    []common.Request
	DoneList    []common.Request
}

var ins *Ser
var once sync.Once

func GetSer() *Ser {
	once.Do(func() {
		ins = &Ser{}
	})
	return ins
}
