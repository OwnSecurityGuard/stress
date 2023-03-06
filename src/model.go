package src

import (
	"encoding/json"

	"github.com/Jeffail/gabs/v2"
	"github.com/jhump/protoreflect/dynamic"
	b3 "github.com/liyakai/behavior3go"
	. "github.com/liyakai/behavior3go/config"
	. "github.com/liyakai/behavior3go/core"
	"github.com/rs/zerolog/log"

	// . "github.com/liyakai/behavior3go/loader"
	"github.com/jmespath/go-jmespath"
)

var (
	TestAllProto         = LoadAllProto() // todo 测试使用
	GlobalReqMapRespName = make(map[string]string)
)

// ReqMsg 请求协议
type ReqMsg struct {
	MsgId    uint32
	MsgName  string
	Title    string
	P        []*Param
	ExtraMap map[string][]string // 存放的值 存放需要从上文获取的参数
	Param    *gabs.Container
	DyMsg    *dynamic.Message
	Composite
}

func (p *ReqMsg) Initialize(setting *BTNodeCfg) {
	// p.Action.Initialize(setting)

	p.MsgId, p.MsgName, p.Title = GetDecoratorInitVal(setting)

	p.DyMsg = GetDyMsgByName(TestAllProto, p.MsgName)
	p.P, p.ExtraMap = MapToParamArr(setting.Properties)
	p.Param = gabs.New()
	for _, v := range p.P {
		v.MarshalJson(p.Param)
	}

	for _, c := range setting.Children {
		m.Lock()
		GlobalReqMapRespName[c] = p.MsgName
		m.Unlock() // TODO
	}

	log.Debug().Msg("ReqMsg init complete ExtraMap  ")

}

func (p *ReqMsg) OnTick(tick *Tick) b3.Status {
	p.SetExtraVal(tick)
	robot := GetRobot(tick)
	log.Debug().Msg("Send Msg")
	robot.Send(p)

	for i := 0; i < p.GetChildCount(); i++ {
		var status = p.GetChild(i).Execute(tick)
		if status != b3.SUCCESS {
			return status
		}
	}
	return b3.SUCCESS
}

func (p *ReqMsg) Serialize() ([]byte, error) {

	p.DyMsg.UnmarshalJSON(p.Param.Bytes())
	return p.DyMsg.Marshal()
}

func (p *ReqMsg) SetExtraVal(tick *Tick) {
	robot := GetRobot(tick)
	for k, v := range p.ExtraMap {
		jsonVal, _ := robot.ctx.Get(v[0]) // todo check
		var data interface{}
		json.Unmarshal(jsonVal.([]byte), &data)
		ctxVal, _ := jmespath.Search(v[1], data)
		p.Param.SetP(ctxVal, k)

	}
}

type RespMsg struct {
	Id         string
	MsgId      uint32
	MsgName    string
	Title      string
	IsInit     bool
	IsExport   bool
	DyMsg      *dynamic.Message
	ReqMsgName string
	WaitCh     chan struct{}
	Composite
}

func (p *RespMsg) Initialize(setting *BTNodeCfg) {
	p.Id = setting.Id
	p.MsgId, p.MsgName, p.Title = GetDecoratorInitVal(setting)
	p.DyMsg = GetDyMsgByName(TestAllProto, p.MsgName)
	// m.Lock()
	// MsgIdMapResp[MsgIdType(p.MsgId)] = p
	// m.Unlock()
	// if resp, ok := MsgIdMapResp[MsgIdType(p.MsgId)]; ok {
	// 	resp[p.Id] = p
	// } else {
	// 	MsgIdMapResp[MsgIdType(p.MsgId)] = map[string]*RespMsg{p.Id: p}
	// }
	p.ReqMsgName = GlobalReqMapRespName[p.Id]
	// fmt.Println(GlobalReqMapRespName)

	// RespIdMapWaitCh[p.Id] = p.WaitCh
}

func (p *RespMsg) OnTick(tick *Tick) b3.Status {

	p.WaitCh = make(chan struct{})
	robot := GetRobot(tick)
	robot.RespIdMapWaitCh[p.Id] = p.WaitCh
	robot.MsgIdMapResp[MsgIdType(p.MsgId)] = p
	robot.WaittingRespId = p.Id
	<-p.WaitCh
	delete(robot.RespIdMapWaitCh, p.Id)
	// a, _ := robot.ctx.Get(p.MsgName)
	// log.Debug().Msg("Resp:" + string(a.([]byte)))
	for i := 0; i < p.GetChildCount(); i++ {
		var status = p.GetChild(i).Execute(tick)
		if status != b3.SUCCESS {
			return status
		}
	}
	return b3.SUCCESS
}

func GetDecoratorInitVal(setting *BTNodeCfg) (msgId uint32, msgName, title string) {
	msgId = uint32(setting.Properties["MsgId"].(float64))
	delete(setting.Properties, "MsgId")
	msgName = setting.Properties["MsgName"].(string)
	delete(setting.Properties, "MsgName")
	title = setting.Title
	return
}
