package app

import (
	"giot/internal/server"
	"github.com/spf13/cobra"
)

func NewGiotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manager-api [flags]",
		Short: "APISIX Manager API",
		RunE: func(cmd *cobra.Command, args []string) error {
			//log.InitLogger()
			server.NewServer()
			return nil
		},
	}
	return cmd
}
