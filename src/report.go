package src

import (
	"fmt"
	"sort"
)

var GlobalReport MsgTps

func InitReportCollect(ch chan *MsgTpsPoint) {

	go func() {

		for point := range ch {
			// log.Printf("point %+v", point)
			GlobalReport = append(GlobalReport, point)
			// log.Printf("GlobalReport")

		}
		// fmt.Println("SSSSSSSSSSSSSSSSSSSSSSSS ", i)
		GlobalReport.GetTpsData()
	}()

}

type MsgTpsPoint struct {
	TimestampSecond int
	MsgName         string
	MsgSize         uint32
	RespTime        int64
}

type MsgTps []*MsgTpsPoint

func (m MsgTps) Len() int { return len(m) }

func (m MsgTps) GetTpsData() {
	var (
		msgCount          = make(map[string]int)
		msgCountPerSecond = make(map[string]map[int]int)
		msgRespTime       = make(map[string][]int64)
		msgTps            = make(map[string]float64)
	)

	for i := range m {
		msgCount[m[i].MsgName] += 1
		msgRespTime[m[i].MsgName] = append(msgRespTime[m[i].MsgName], m[i].RespTime)
		if countPerSecond, ok := msgCountPerSecond[m[i].MsgName]; ok {
			countPerSecond[m[i].TimestampSecond] += 1
		} else {
			msgCountPerSecond[m[i].MsgName] = map[int]int{m[i].TimestampSecond: 1}
		}
	}

	for k, v := range msgCountPerSecond {
		msgTps[k] = CalculateTps(v)
	}

	fmt.Println("report result:")
	fmt.Println("协议名  请求数量  tps ")
	for k, v := range msgTps {
		fmt.Printf("%s  %d  %v \r\n", k, msgCount[k], v)
	}
	// todo print time
}

func CalculateTps(m map[int]int) float64 {
	var (
		tps   = make([]int, len(m))
		start = 0
	)
	for _, v := range m {
		tps[start] = v
		start += 1
	}
	sort.Ints(tps)
	r3 := len(tps) / 4 * 3
	r1 := len(tps) / 4

	if r1 > 0 {
		irq := float64(tps[r3] - tps[r1])
		bigVal := float64(tps[r3]) + 1.5*irq
		lesVal := float64(tps[r1]) - 1.5*irq

		var newTps []int
		for _, v := range tps {
			if v <= int(bigVal) && v >= int(lesVal) {
				newTps = append(newTps, v)
			}
		}

		return Avg(newTps)
	}
	return Avg(tps)

}

func Avg(data []int) float64 {
	sum := 0
	for _, d := range data {
		sum += d
	}
	// fmt.Println(data)
	// fmt.Println("aaaaaaa ", float64(sum), "   ", float64(len(data)))
	return float64(sum) / float64(len(data))
}
