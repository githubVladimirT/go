package cmd

import "github.com/gookit/gcli/v3"

var ShowDemo = &gcli.Command{
	Name: "show",
	Func: runShow,
	//
	Desc: "the command will show some data format methods",
}

func runShow(c *gcli.Command, _ []string) error {

	return nil
}
