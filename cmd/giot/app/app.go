package app

import (
	"giot/internal/conf"
	"giot/internal/log"
	"giot/internal/server"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func NewGiotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "giot 设备接入平台",
		Short: "giot",
		Long: dedent.Dedent(`
				┌──────────────────────────────────────────────────────────┐
			    │ FlYING                                                   │
			    │ Cloud Native Distributed Configuration Center             │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/ZhangChengJi/flyingv2/issues           │
			    └──────────────────────────────────────────────────────────┘
		`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf.InitConf()
			log.InitLogger()
			s := server.NewServer()
			errSig := make(chan error, 5)
			s.Start(errSig)
			// Signal received to the process externally.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case sig := <-quit:
				log.Infof("The Manager API server receive %s and start shutting down", sig.String())
				s.Stop()
				log.Infof("See you next time!")
			case err := <-errSig:
				log.Errorf("The Manager API server start failed: %s", err.Error())
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&conf.ConfigFile, "config", "c", "", "config file")
	cmd.PersistentFlags().StringVarP(&conf.WorkDir, "work-dir", "p", ".", "current work directory")
	return cmd
}
