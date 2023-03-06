package src

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"

	"net"

	// . "github.com/liyakai/behavior3go/config"
	. "github.com/liyakai/behavior3go/core"
	// . "github.com/liyakai/behavior3go/loader"
	"github.com/rs/zerolog/log"
)

var (
	m            sync.Mutex
	MsgIdMapResp = make(map[MsgIdType]*RespMsg)
	// RespIdMapWaitCh = make(map[string]chan struct{})
)

type Robot struct {
	Id              string
	ctx             Context
	conn            net.Conn
	codec           Codec
	WaittingRespId  string
	RespIdMapWaitCh map[string]chan struct{}
	MsgIdMapResp    map[MsgIdType]*RespMsg
	TimeoutCh       chan struct{}
	ReportCh        chan *MsgTpsPoint
}

func NewRobot(addr string, reportCh chan *MsgTpsPoint) *Robot {
	conn, _ := net.Dial("tcp", addr)

	return &Robot{
		Id:              randStr(10),
		conn:            conn,
		ctx:             NewContext(),
		codec:           TestCodec{},
		RespIdMapWaitCh: make(map[string]chan struct{}),
		MsgIdMapResp:    make(map[MsgIdType]*RespMsg),
		// WaitCh:    make(chan struct{}),
		TimeoutCh: make(chan struct{}),
		ReportCh:  reportCh,
	}
}

func NewClient() net.Conn {
	conn, err := net.Dial("tcp", ServerAddr)
	if err != nil {
		fmt.Println("err : ", err)
		return nil
	}
	return conn
}

func GetRobot(tick *Tick) *Robot {
	if inter := tick.Blackboard.GetMem(RobotKey); inter != nil {
		r := inter.(*Robot)
		return r
	}
	return nil

}

func (r *Robot) Send(msg *ReqMsg) { // 发送协议
	data, _ := msg.Serialize()

	data, _ = r.codec.Encode(MsgIdType(msg.MsgId), data) // todo check
	log.Debug().Str("Client => Server:", msg.MsgName)
	r.SetStartTime(msg.MsgName)
	r.conn.Write(data)

}

func (r *Robot) Listen() { // 接收服务端推送协议消息
	// r.MsgIdMapResp
	for {
		msgId, data, err := r.codec.Decode(r.conn)
		// log.Debug().Msg("Server => Clien")
		if err != nil {
			log.Error().Err(err)
			return
		}
		if msgId == 0 {
			return
		}
		if msg, ok := r.MsgIdMapResp[MsgIdType(msgId)]; ok {
			log.Debug().Msg("Server => Client")
			msg.DyMsg.Unmarshal(data) // todo check
			jsonData, _ := msg.DyMsg.MarshalJSON()
			// log.Debug().Str("Server => Client: 11", string(jsonData))
			// nowTime := time.Now()
			r.ReportCh <- &MsgTpsPoint{
				TimestampSecond: time.Now().Second(),
				MsgName:         msg.MsgName,
				MsgSize:         uint32(len(data)),
				RespTime:        int64(time.Since(r.GetStartTime(msg.ReqMsgName))),
			}

			r.ctx.Set(msg.MsgName, jsonData)
			// time.Sleep(1 * time.Second)
			if ch, ok := r.RespIdMapWaitCh[r.WaittingRespId]; ok {
				close(ch)
			}

		} else {
			log.Debug().Str("Not Map MesId ", fmt.Sprintf("%v", msgId))
		}
	}

}

func (r *Robot) GetWaitCh(resp *RespMsg) chan struct{} {
	inter, _ := r.ctx.Get(resp.Id + resp.MsgName + "waitCh")
	return inter.(chan struct{})
}

func (r *Robot) GetStartTime(msgReqName string) time.Time { //
	val, ok := r.ctx.Pop(msgReqName + "StartTime")
	if !ok {
		fmt.Println(" errrrrrr", r.Id)
		return time.Now()
	}
	return val.(time.Time)
}

func (r *Robot) SetStartTime(msgReqName string) { //
	r.ctx.Set(msgReqName+"StartTime", time.Now())

}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
