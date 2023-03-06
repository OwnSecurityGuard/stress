package src

import "flag"

var (
	TcpAddr  string
	Num      int
	Spawn    int
	ProtoDir string
	TreePath string
)

func InitFlag() {
	flag.StringVar(&TcpAddr, "tcpAddr", "127.0.0.1:7898", "server addr")
	flag.IntVar(&Num, "num", 10, "number of robot")
	flag.IntVar(&Spawn, "spawn", 10, "robot spawn per second")
	flag.StringVar(&ProtoDir, "protoDir", "./protoA", "proto file path")
	flag.StringVar(&TreePath, "treePath", "./bz.bs", "behavior tree file path")
	flag.Parse()
}
