package handlers

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/ghmeier/go-service"
)

type Avail struct {
	*base
	flags *flag.FlagSet
	help  bool
	id    string
}

type Availabilities struct {
	Results []*Availability `json:"results,omitempty"`
	Next    bool            `json:"hasNext"`
}
type Availability struct {
	ID             string     `json:"_id"`
	UserId         string     `json:"userId,omitempty"`
	Location       string     `json:"location,omitempty"`
	Description    string     `json:"description,omitempty"`
	Title          string     `json:"title,omitempty"`
	Timezone       string     `json:"timezone,omitempty"`
	CalendarID     string     `json:"calendarId,omitempty"`
	CalendarName   string     `json:"calendarName,omitempty"`
	Modify         bool       `json:"guestsCanModify,omitempty"`
	Timeslots      []Timeslot `json:"timeslots,omitempty"`
	DoubleBookings bool       `json:"allowDoubleBookings,omitempty"`
	Guests         []Guest    `json:"guests,omitempty"`
}

type Timeslot struct {
	Start  time.Time `json:"start,omitempty"`
	End    time.Time `json:"end,omitempty"`
	Events []Event   `json:"events,omitempty"`
}

type Event struct {
	Guest Guest  `json:"guest,omitempty"`
	ID    string `json:"id,omitempty"`
}

type Guest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
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
		a.getByID(a.id)
		return
	}

	a.get()
}

func (a *Avail) Help() {
	fmt.Println("avail: mixmax availability")
	a.flags.PrintDefaults()
}

func (a *Avail) get() {
	results := make([]*Availability, 0)

	var res Availabilities
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

	fmt.Printf("Found %d availabilities:\n", len(results))
	for _, v := range res.Results {
		a.printAvailability(v)
	}
}

func (a *Avail) getByID(id string) {
	var res Availability
	err := a.s.Send(&service.Request{
		Method:  http.MethodGet,
		URL:     fmt.Sprintf("%s/%s", a.Url(), id),
		Headers: map[string]string{"X-API-Token": a.Key()},
	}, &res)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	a.printAvailability(&res)
}

func (a *Avail) printAvailability(v *Availability) {
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
