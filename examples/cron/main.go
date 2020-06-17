package main

import (
	"fmt"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/cron"
	"github.com/firmeve/firmeve/kernel/contract"
)

func main() {
	firmeve.RunWithSupportFunc(func(application contract.Application) {
		cron := application.Resolve(`cron`).(contract.Cron)
		cron.EverySecondWithFunc(EverySecond)
		cron.SecondsWithFunc(5, Second5)
		cron.EveryMinuteWithFunc(EveryMinute)
	}, firmeve.WithConfigPath("./config.yaml"), firmeve.WithProviders([]contract.Provider{
		new(cron.Provider),
	}), firmeve.WithCommands([]contract.Command{
		new(cron.Command),
	}))
}

func EverySecond() {
	fmt.Println(`every second`)
}

func Second5() {
	fmt.Println(`second 5`)
}

func EveryMinute() {
	fmt.Println(`every second`)
}
