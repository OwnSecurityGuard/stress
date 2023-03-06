package src

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/rs/zerolog/log"
)

func LoadAllProto() []*desc.FileDescriptor {
	var (
		pkgPath = PkgPath // 依据projectID
		parser  protoparse.Parser
	)

	fileNames := GetDirFileName(pkgPath)

	parser.ImportPaths = append(parser.ImportPaths, pkgPath)
	ds, err := parser.ParseFiles(fileNames...)
	if err != nil {
		log.Error().Err(err)
		return nil
	}
	// GlobalProtoFileDescriptor.Store(projectId, ds)
	return ds
}

// GetDirFileName 获取目录下的文件名
func GetDirFileName(p string) []string {
	var fileName []string

	err := filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), ".proto") {
				fileName = append(fileName, info.Name())
			}
			return nil
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err)
		return nil
	}
	return fileName
}

// GetDyMsgByName 传入 包和协议名, 返回对应的实体mesf
func GetDyMsgByName(fileDescriptors []*desc.FileDescriptor, msgName string) *dynamic.Message {

	for _, fd := range fileDescriptors {
		log.Debug().Str("GetDyMsgByName: ", fd.GetPackage()+"."+msgName)
		// fd.FindMessage(fd.GetPackage() + "." + msgName);
		if m := fd.FindMessage(fd.GetPackage() + "." + msgName); m != nil {
			dy := dynamic.NewMessage(m)
			// log.Debug().Str()
			return dy
		}

	}
	return nil
}

const (
	ItemType_Val   = 1
	ItemType_Array = 2
	ItemType_Obj   = 3
)

type JsonItem struct {
	Key      string
	ItemType int
}

type Param struct {
	Items []*JsonItem
	Val   string
}

func (p *Param) MarshalJson(g *gabs.Container) {
	var (
		item    = gabs.New()
		isFirst bool
		str     string
	)
	for s := len(p.Items) - 1; s >= 0; s -= 1 {

		if len(str) == 0 {
			str = p.Items[s].Key
		} else {
			str = p.Items[s].Key + "." + str
		}
		if p.Items[s].ItemType == ItemType_Array {

			item1 := gabs.New()
			if !isFirst {
				item1.ArrayAppendP(p.Val, str)
				isFirst = true
			} else {
				item1.ArrayAppendP(item, str)

			}
			item = item1
			str = ""
		}
	}

	if len(str) != 0 {
		if !isFirst {
			g.SetP(p.Val, str)
		} else {
			g.SetP(item, str)
		}
	} else {
		// fmt.Println("aaa", item.String())
		g.Merge(item)
		// g.ArrayAppendP(item, p.Items[0].Key)

	}

}

func MapToParamArr(m map[string]interface{}) (ps []*Param, ctxMap map[string][]string) {

	ctxMap = make(map[string][]string)

	for k, v := range m {

		if p := extractParam(k); p != nil {
			p.Val = fmt.Sprintf("%v", v)
			ps = append(ps, p)
		}

		if val, ok := v.(string); ok && strings.HasPrefix(val, "$") { // todo check
			ctxVal := strings.TrimPrefix(val, "$")

			ctxArr := strings.Split(ctxVal, ".") // $LoginResp.ab.val => LoginResp   ab.val
			// log.Printf("ctxArr %v  ctxVal %v  ", ctxArr, ctxVal)
			jmesPath := strings.TrimPrefix(ctxVal, ctxArr[0]+".")

			ctxMap[k] = []string{ctxArr[0], jmesPath}
		}

	}
	return
}

// player[0].username[2].mcnb   ==>  array play  array username  val macnb
func extractParam(s string) *Param { // todo 目前不考虑数组的有序性
	var (
		items []*JsonItem
		// strs []string
		isArray bool
	)
	arr := strings.Split(s, ".")
	if len(arr) == 0 {
		return nil
	}

	for i := range arr {
		item := &JsonItem{Key: arr[i]}
		if isArray || hasSquareBracketsWithNumber(arr[i]) {
			item.ItemType = ItemType_Array
			item.Key = strings.Split(arr[i], "[")[0] // 不考虑属于数组的第几个了
		} else {
			item.ItemType = ItemType_Val
		}

		items = append(items, item)
	}

	return &Param{Items: items}

}

func hasSquareBracketsWithNumber(s string) bool {
	re := regexp.MustCompile(`\[[0-9]+\]`)
	return re.MatchString(s)
}
