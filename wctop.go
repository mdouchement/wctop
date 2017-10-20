package main

import (
	"fmt"

	"github.com/mdouchement/wctop/server"
	"github.com/spf13/cobra"
)

//go:generate esc -ignore .DS_Store -o server/assets.go -pkg server -prefix server server/assets
//go:generate go run ./generate_buildtime.go

var (
	binding string
	port    string
)

func main() {
	c := &cobra.Command{
		Use:   "wctop",
		Short: "Web UI for ctop, containers top",
		Long:  "Web UI for ctop, containers top",
		RunE:  action,
	}

	c.Flags().StringVarP(&binding, "binding", "b", "0.0.0.0", "Server's binding")
	c.Flags().StringVarP(&port, "port", "p", "5000", "Server's port")

	if err := c.Execute(); err != nil {
		fmt.Println(err)
	}
}

func action(c *cobra.Command, args []string) error {
	return server.Run(fmt.Sprintf("%s:%s", binding, port))
}
