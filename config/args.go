package config

import "flag"

type arg struct {
	Env string
}

var Arg *arg

func InitArgs() {

	Arg = &arg{}

	flag.StringVar(&Arg.Env, "env", "debug", "环境变量")
	flag.Parse()
}

