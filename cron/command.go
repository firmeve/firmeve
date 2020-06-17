package cron

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

type Command struct {
}

func (c Command) CobraCmd() *cobra.Command {
	cmd := new(cobra.Command)
	cmd.Use = "cron"
	cmd.Short = "Scheduled tasks"
	return cmd
}

func (c Command) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	var (
		logger = root.Resolve(`logger`).(contract.Loggable)
		cron   = root.Resolve(`cron`).(contract.Cron)
	)

	logger.Debug(`Start crontab`)
	cron.Start()

	logger.Debug(`Crontab start successful`)

	logger.Debug(`Signal listen SIGTERM(kill),SIGINT(kill -2),SIGKILL(kill -9)`)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit // 阻塞直至有信号传入

	//
	cron.Stop()
	logger.Debug(`Stoped crontab`)
}
