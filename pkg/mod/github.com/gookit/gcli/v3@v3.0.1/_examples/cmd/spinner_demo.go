package cmd

import (
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/progress"
)

type spinnerDemo struct {
	speed    int
	themeNum int
}

var spOpts = spinnerDemo{}

var SpinnerDemo = &gcli.Command{
	Name:    "spinner",
	Desc:    "there are some CLI spinner bar run demos",
	Aliases: []string{"spr", "spr-demo"},
	// Func:    spOpts.Run,
	Config: func(c *gcli.Command) {
		c.IntOpt(&spOpts.speed, "speed", "s", 100, "setting the spinner running speed")
		c.IntOpt(&spOpts.themeNum, "theme-num", "t", 0, "setting the theme numbering. allow: 0 - 16")

		c.AddArg("name",
			"spinner type name. allow: loading,roundTrip",
			false,
		)
	},
	Examples: `Loading spinner:
  {$fullCmd} loading
roundTrip spinner:
  {$fullCmd} roundTrip`,
	Func: func(c *gcli.Command, _ []string) error {
		name := c.Arg("name").String()

		switch name {
		case "", "spinner", "load", "loading":
			spOpts.runLoadingSpinner()
		case "rt", "roundTrip":
			spOpts.runRoundTripSpinner()
		default:
			return c.Errorf("the spinner type name only allow: loading,roundTrip. input is: %s", name)
		}
		return nil
	},
}

func (sd *spinnerDemo) runRoundTripSpinner() {
	s := progress.RoundTripLoading(
		progress.GetCharTheme(sd.themeNum),
		time.Duration(sd.speed)*time.Millisecond,
	)

	// s.Start("%s work handling ... ...")
	s.Start("[%s] work handling ... ...")

	// Run for some time to simulate work
	time.Sleep(4 * time.Second)
	s.Stop("work handle complete")
}

func (sd *spinnerDemo) runLoadingSpinner() {
	s := progress.LoadingSpinner(
		progress.GetCharsTheme(sd.themeNum),
		time.Duration(sd.speed)*time.Millisecond,
	)

	s.Start("%s work handling ... ...")
	// Run for some time to simulate work
	time.Sleep(4 * time.Second)
	s.Stop("work handle complete")
}
