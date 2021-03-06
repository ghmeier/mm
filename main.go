package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghmeier/go-mixmax"
	"github.com/ghmeier/mm/handlers"
)

const keyFile = ".mixmax_key"

func main() {
	m := &mm{}
	m.Init()

	if m.key == "" {
		flag.Parse()
		m.handleKey()
		return
	}

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		return
	}

	m.Handle()
}

type mm struct {
	key   string
	c     *mixmax.Client
	nkey  string
	help  bool
	appt  handlers.Handler
	avail handlers.Handler
	code  handlers.Handler
}

func (m *mm) Init() {
	m.key = apiKey()
	m.c = mixmax.New(m.key)

	flag.StringVar(&m.nkey, "key", "", "your mixmax api key")
	flag.BoolVar(&m.help, "help", false, "get detailed usage information")
	flag.Usage = m.Usage

	m.appt = &handlers.Appt{Client: m.c.AppointmentLinks}
	m.appt.Init(m.key)

	m.avail = &handlers.Avail{Client: m.c.Availability}
	m.avail.Init(m.key)

	m.code = &handlers.Code{Client: m.c.CodeSnippet}
	m.code.Init(m.key)
}

func (m *mm) Handle() {
	switch os.Args[1] {
	case m.appt.Cmd():
		m.appt.Go(os.Args[2:])
	case m.avail.Cmd():
		m.avail.Go(os.Args[2:])
	case m.code.Cmd():
		m.code.Go(os.Args[2:])
	default:
		flag.Parse()
		if m.nkey != "" {
			m.handleKey()
			return
		}

		if m.help {
			m.Usage()
		}
	}
}

func (m *mm) Usage() {
	flag.PrintDefaults()
	m.appt.Help()
	m.avail.Help()
	m.code.Help()
}

func (m *mm) handleKey() {
	if m.nkey == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	writeKey(m.nkey)
	fmt.Println("API Key is successfully set")
}

func apiKey() string {
	raw, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return ""
	}

	return string(raw)
}

func writeKey(k string) {
	ioutil.WriteFile(keyFile, []byte(k), os.ModePerm)
}
