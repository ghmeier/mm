package handlers

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/ghmeier/go-service"
)

type Code struct {
	*base
	flags *flag.FlagSet
	help  bool
	id    string
}

type CodeSnippets struct {
	Results []*CodeSnippet `json:"results,omitempty"`
	Next    bool           `json:"hasNext"`
}
type CodeSnippet struct {
	ID         string `json:"_id"`
	UserId     string `json:"userId,omitempty"`
	HTML       string `json:"html,omitempty"`
	Title      string `json:"title,omitempty"`
	Background string `json:"background,omitempty"`
	Code       string `json:"code,omitempty"`
	Theme      string `json:"theme,omitempty"`
	Language   string `json:"language,omitempty"`
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
		a.getByID(a.id)
		return
	}

	a.get()
}

func (a *Code) Help() {
	fmt.Println("code: list created formatted code enhancements")
	a.flags.PrintDefaults()
}

func (a *Code) get() {
	results := make([]*CodeSnippet, 0)

	var res CodeSnippets
	res.Next = true
	for res.Next {
		err := a.s.Send(&service.Request{
			Method:  http.MethodGet,
			URL:     fmt.Sprintf("%s", a.Url()),
			Headers: map[string]string{"X-API-Token": a.Key()},
		}, &res)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		results = append(results, res.Results...)
	}

	fmt.Printf("Found %d code snippets:\n", len(results))
	for _, v := range res.Results {
		a.printCode(v)
	}
}

func (a *Code) getByID(id string) {
	var res CodeSnippet
	err := a.s.Send(&service.Request{
		Method:  http.MethodGet,
		URL:     fmt.Sprintf("%s/%s", a.Url(), id),
		Headers: map[string]string{"X-API-Token": a.Key()},
	}, &res)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	a.printCode(&res)
}

func (a *Code) printCode(v *CodeSnippet) {
	fmt.Printf("ID: %s Title:\t%s\n", v.ID, v.Title)
	if v.Language != "" {
		fmt.Printf("Language:\t%s\n", v.Language)
	}
}
