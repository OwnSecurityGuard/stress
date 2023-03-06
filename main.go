package main

import (
	_ "net/http/pprof"
	"stress/src"
	. "stress/src"
	"sync"
	"time"

	// "github.com/jmespath/go-jmespath"
	b3 "github.com/liyakai/behavior3go"
	. "github.com/liyakai/behavior3go/config"
	. "github.com/liyakai/behavior3go/core"
	. "github.com/liyakai/behavior3go/loader"
	"github.com/rs/zerolog"
)

func main() {

	// go func() {
	// 	log.Println(http.ListenAndServe(":6060", nil))
	// }()
	src.InitFlag()

	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	ch := make(chan *src.MsgTpsPoint, 100)

	src.InitReportCollect(ch)

	var wg sync.WaitGroup

	wg.Add(src.Num)
	tick := time.NewTicker(1 * time.Second)
	for start := 0; start < src.Num; {
		<-tick.C
		var n = src.Spawn
		if start+src.Spawn > src.Num {
			n = start + src.Spawn - src.Num

		}
		for i := 0; i < n; i++ {
			go func() {
				maps := b3.NewRegisterStructMaps()

				maps.Register("ReqMsg", new(ReqMsg))
				maps.Register("RespMsg", new(RespMsg))
				projectConfig, _ := LoadRawProjectCfg("D:\\workspace\\xxxxx\\gt.b3")
				board := NewBlackboard()
				robot := NewRobot(ServerAddr, ch)
				board.SetMem(RobotKey, robot)
				root := projectConfig.Data.Trees[0]
				tree := CreateBevTreeFromConfig(&root, maps)
				go robot.Listen()
				tree.Tick(1, board)
				wg.Done()
			}()
		}
		start += src.Spawn
	}

	wg.Wait()
	close(ch)
	time.Sleep(time.Second * 3)
	// close(ch)

	// tree.Tick(3, board)
}
