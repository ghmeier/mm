package handlers

import (
	"flag"
	"fmt"

	"github.com/ghmeier/go-mixmax/appointmentlinks"
)

type Appt struct {
	*base
	Client *appointmentlinks.Client
	flags  *flag.FlagSet
	name   string
	update string
	help   bool
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
		links, err := a.Client.Me()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Mixmax Calendar URL: https://cal.mixmax.com/%s\n", links.Name)
		return
	}

	if a.name == "" && a.update != "" {
		err := a.Client.Set(a.update)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Successfully Updated Mixmax Calendar URL.")
		return
	}

	_, err := a.Client.Get(a.name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Retrieved Appointment Links.")
}

func (a *Appt) Help() {
	fmt.Println("appt: get mixmax calendar url")
	a.flags.PrintDefaults()
}
