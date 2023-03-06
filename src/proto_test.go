package src

import (
	"testing"

	b3 "github.com/liyakai/behavior3go"
	. "github.com/liyakai/behavior3go/config"
	. "github.com/liyakai/behavior3go/core"
	. "github.com/liyakai/behavior3go/loader"
)

func Test(Runt *testing.T) {
	maps := b3.NewRegisterStructMaps()

	maps.Register("proto", new(ReqMsg))
	projectConfig, _ := LoadRawProjectCfg("D:\\studyProject\\robotgo-main\\trees\\game.b3")

	board := NewBlackboard()
	board.SetMem(RobotKey, NewClient())
	root := projectConfig.Data.Trees[0]

	a := CreateBevTreeFromConfig(&root, maps)
	a.Tick(2, board)
}
