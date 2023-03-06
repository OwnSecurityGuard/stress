package src

import (
	"fmt"
	"testing"

	"github.com/Jeffail/gabs/v2"
)

func TestLoadAllProto(t *testing.T) {
	// got := LoadAllProto()
	// for _, g := range got {
	// 	fmt.Println(g.GetName())
	// }
	// j := ` {"name":"232", "password":"apppppp"  }`

	// d := JsonToMes([]byte(j), got, "LoginReq")
	// if d != nil {
	// 	fmt.Println(d.Marshal())
	// }

}

func TestJsonParse(t *testing.T) {

	/*
			  {
			  	"name":"232", "password":"apppppp" ,
		      "testItem": [ {"name":"aa"}, {"name":"vvv"}],
		      "d": [1,2,3,4]
		    }

	*/
	g := gabs.New()

	g.SetP("name", "asdasd")

	b, _ := g.MarshalJSON()
	fmt.Println(string(b))
}

func TestParam_MarshalJson(t *testing.T) {
	type args struct {
		g *gabs.Container
	}
	tests := []struct {
		name string
		p    *Param
		args args
	}{
		{
			"构造带有数组对象的json",
			&Param{
				Items: []*JsonItem{
					{"Player", ItemType_Val},
					{"user", ItemType_Array},
					// { "Player", ItemType_Val,},
					// { "Player", ItemType_Val,},
					// { "Player", ItemType_Val,},
				},
				Val: "a",
			},
			args{gabs.New()},
		},
		{
			"构造没有数组对象的json",
			&Param{
				Items: []*JsonItem{
					{"Player", ItemType_Val},
					{"user", ItemType_Val},
					// { "Player", ItemType_Val,},
					// { "Player", ItemType_Val,},
					// { "Player", ItemType_Val,},
				},
				Val: "a",
			},
			args{gabs.New()},
		},
		{
			"构造第一个为数组对象的json",
			&Param{
				Items: []*JsonItem{
					{"Player", ItemType_Array},
					{"user", ItemType_Val},
					// { "Player", ItemType_Val,},
					// { "Player", ItemType_Val,},
					// { "Player", ItemType_Val,},
				},
				Val: "a",
			},
			args{gabs.New()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.MarshalJson(tt.args.g)
			fmt.Println(tt.args.g.String())
		})
	}
}
