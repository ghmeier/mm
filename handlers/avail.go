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
}

type Availabilities struct {
	Results []*Availability `json:"results,omitempty"`
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
	a.flags.BoolVar(&a.help, "help", false, "get detailed usage information")

}

func (a *Avail) Go(args []string) {
	a.flags.Parse(args)
	if a.help {
		a.Help()
		return
	}

	a.get()
}

func (a *Avail) Help() {
	fmt.Println("avail: mixmax availability")
	a.flags.PrintDefaults()
}

func (a *Avail) get() {
	var res Availabilities
	err := a.s.Send(&service.Request{
		Method:  http.MethodGet,
		URL:     fmt.Sprintf("%s", a.Url()),
		Headers: map[string]string{"X-API-Token": a.Key()},
	}, &res)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range res.Results {
		fmt.Printf("%+v\n", v)
	}
}
