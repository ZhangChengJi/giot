package app

import (
	"giot/conf"
	"giot/internal/scheduler/server"
	"giot/pkg/log"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func NewSchedulerCommand() *cobra.Command {
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
			zap.ReplaceGlobals(log.New())
			s := server.NewServer()
			errSig := make(chan error, 5)
			s.Start(errSig)
			// Signal received to the process externally.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case sig := <-quit:
				log.Sugar.Infof("The Manager API server receive %s and start shutting down", sig.String())
				s.Stop()
				log.Sugar.Info("See you next time!")
			case err := <-errSig:
				log.Sugar.Errorf("The Manager API server start failed: %s", zap.Error(err))
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&conf.ConfigFile, "config", "c", "", "config file")
	cmd.PersistentFlags().StringVarP(&conf.WorkDir, "work-dir", "p", ".", "current work directory")
	return cmd
}
