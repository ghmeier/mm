package handlers

import ()

type Handler interface {
	Init(string)
	Go([]string)
	Cmd() string
	Help()
}

type base struct {
	cmd string
}

func new(key, cmd, path string) *base {
	return &base{
		cmd: cmd,
	}
}

func (b *base) Cmd() string {
	return b.cmd
}
