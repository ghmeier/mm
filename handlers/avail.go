package handlers

import (
	"flag"
	"fmt"

	"github.com/ghmeier/go-mixmax/availability"
	"github.com/ghmeier/go-mixmax/models"
)

type Avail struct {
	*base
	Client *availability.Client
	flags  *flag.FlagSet
	help   bool
	id     string
}

func (a *Avail) Init(key string) {
	a.base = new(key, "avail", "availability")
	a.flags = flag.NewFlagSet(a.Cmd(), flag.ExitOnError)
	a.flags.StringVar(&a.id, "id", "", "id of specific availability")
	a.flags.BoolVar(&a.help, "help", false, "get detailed usage information")

}

func (a *Avail) Go(args []string) {
	a.flags.Parse(args)
	if a.help {
		a.Help()
		return
	}

	if a.id != "" {
		avail, err := a.Client.Get(a.id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		a.printAvailability(avail)
		return
	}

	avails, err := a.Client.List()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Found %d availabilities:\n", len(avails.Results))
	for _, v := range avails.Results {
		a.printAvailability(v)
	}
}

func (a *Avail) Help() {
	fmt.Println("avail: mixmax availability")
	a.flags.PrintDefaults()
}

func (a *Avail) printAvailability(v *models.Availability) {
	fmt.Printf("ID: %s Title:\t%s\n", v.ID, v.Title)
	if v.Location != "" {
		fmt.Printf("Location:\t%s\n", v.Location)
	}
	if v.Description != "" {
		fmt.Printf("Description: %s\n", v.Description)
	}
	for _, t := range v.Timeslots {
		fmt.Printf("\tFrom %s to %s\n", t.Start.Format("Mon, May 2 1:00 AM"), t.End.Format("Mon, May 2 1:00 AM"))
	}
}
