package handlers

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/ghmeier/go-service"
)

type Appt struct {
	*base
	flags  *flag.FlagSet
	name   string
	update string
	help   bool
}

type AppointmentLinks struct {
	UserId string `json:"userId,omitempty"`
	Name   string `json:"name"`
}

func (a *Appt) Init(key string) {
	a.base = new(key, "appt", "appointmentlinks")
	a.flags = flag.NewFlagSet(a.Cmd(), flag.ExitOnError)
	a.flags.StringVar(&a.name, "name", "", "search for appointment with this name")
	a.flags.StringVar(&a.update, "update", "", "update mixmax calendar name with to this name")
	a.flags.BoolVar(&a.help, "help", false, "get detailed usage information")

}

func (a *Appt) Go(args []string) {
	a.flags.Parse(args)
	if a.help {
		a.Help()
		return
	}

	if a.name == "" && a.update == "" {
		a.me()
		return
	}

	if a.name == "" && a.update != "" {
		a.set()
		return
	}

	// send other request
	a.get()
}

func (a *Appt) Help() {
	fmt.Println("appt: get mixmax calendar url")
	a.flags.PrintDefaults()
}

func (a *Appt) me() {
	var res AppointmentLinks
	err := a.s.Send(&service.Request{
		Method:  http.MethodGet,
		URL:     fmt.Sprintf("%s/me", a.Url()),
		Headers: map[string]string{"X-API-Token": a.Key()},
	}, &res)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Appointments at https://cal.mixmax.com/%s\n", res.Name)
}

func (a *Appt) set() {
	req := &AppointmentLinks{Name: a.update}

	err := a.s.Send(&service.Request{
		Method: http.MethodPatch,
		URL:    fmt.Sprintf("%s/me", a.Url()),
		Data:   req,
		Headers: map[string]string{
			"X-API-Token": a.Key(),
		},
	}, nil)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Updated mixmax calendar name to https://cal.mixmax.com/%s\n", a.update)
}

func (a *Appt) get() {
	var res AppointmentLinks
	err := a.s.Send(&service.Request{
		Method:  http.MethodGet,
		URL:     fmt.Sprintf("%s?name=%s", a.Url(), a.name),
		Headers: map[string]string{"X-API-Token": a.Key()},
	}, &res)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Got appointments:\n")
}
