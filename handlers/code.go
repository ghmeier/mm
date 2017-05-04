package handlers

import (
	"flag"
	"fmt"

	"github.com/ghmeier/go-mixmax/codesnippet"
	"github.com/ghmeier/go-mixmax/models"
)

type Code struct {
	*base
	Client *codesnippet.Client
	flags  *flag.FlagSet
	help   bool
	id     string
}

func (a *Code) Init(key string) {
	a.base = new(key, "code", "codesnippets")
	a.flags = flag.NewFlagSet(a.Cmd(), flag.ExitOnError)
	a.flags.StringVar(&a.id, "id", "", "id of specific code snippet")
	a.flags.BoolVar(&a.help, "help", false, "get detailed usage information")

}

func (a *Code) Go(args []string) {
	a.flags.Parse(args)
	if a.help {
		a.Help()
		return
	}

	if a.id != "" {
		code, err := a.Client.Get(a.id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		a.printCode(code)
		return
	}

	codes, err := a.Client.List()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Found %d code snippets:\n", len(codes.Results))
	for _, v := range codes.Results {
		a.printCode(v)
	}
}

func (a *Code) Help() {
	fmt.Println("code: list created formatted code enhancements")
	a.flags.PrintDefaults()
}

func (a *Code) printCode(v *models.CodeSnippet) {
	fmt.Printf("ID: %s Title:\t%s\n", v.ID, v.Title)
	if v.Language != "" {
		fmt.Printf("Language:\t%s\n", v.Language)
	}
}
