package cron

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/robfig/cron/v3"
	"strconv"
)

type (
	Cron struct {
		cron *cron.Cron
	}
)

func New() contract.Cron {
	return &Cron{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (c *Cron) RunWithFunc(spec string, job func()) {
	c.cron.AddFunc(spec, job)
	c.cron.Start()
	c.cron.Stop()
}

func (c *Cron) RunWithJob(spec string, job contract.CronJob) {
	c.cron.AddJob(spec, job.(cron.Job))
}

func (c *Cron) SecondsWithFunc(seconds int, job func()) {
	string := strconv.Itoa(seconds)
	c.RunWithFunc(`@every `+string+`s`, job)
}

func (c *Cron) SecondsWithJob(seconds int, job contract.CronJob) {
	string := strconv.Itoa(seconds)
	c.RunWithJob(`@every `+string+`s`, job)
}

func (c *Cron) EverySecondWithFunc(job func()) {
	c.RunWithFunc(`@every 1s`, job)
}

func (c *Cron) EverySecondWithJob(job contract.CronJob) {
	c.RunWithJob(`@every 1s`, job)
}

func (c *Cron) MinutesWithFunc(minutes int, job func()) {
	string := strconv.Itoa(minutes)
	c.RunWithFunc(`@every `+string+`m`, job)
}

func (c *Cron) MinutesWithJob(minutes int, job contract.CronJob) {
	string := strconv.Itoa(minutes)
	c.RunWithJob(`@every `+string+`m`, job)
}

func (c *Cron) EveryMinuteWithFunc(job func()) {
	c.RunWithFunc(`@every 1m`, job)
}

func (c *Cron) EveryMinuteWithJob(job contract.CronJob) {
	c.RunWithJob(`@every 1m`, job)
}

func (c *Cron) HoursWithFunc(hours int, job func()) {
	string := strconv.Itoa(hours)
	c.RunWithFunc(`@every `+string+`h`, job)
}

func (c *Cron) HoursWithJob(hours int, job contract.CronJob) {
	string := strconv.Itoa(hours)
	c.RunWithJob(`@every `+string+`h`, job)
}

func (c *Cron) EveryHourWithFunc(job func()) {
	c.RunWithFunc(`@every 1h`, job)
}

func (c *Cron) EveryHourWithJob(job contract.CronJob) {
	c.RunWithJob(`@every 1h`, job)
}

func (c *Cron) DayAtWithFunc(hour int, job func()) {
	//string := strconv.Itoa(hour)
	//c.RunWithFunc(`* * * * * ?`, job)
}

func (c *Cron) DayAtWithJob(hour int, job contract.CronJob) {
	panic("implement me")
}

func (c *Cron) EveryDayWithFunc(job func()) {
	c.RunWithFunc(`@daily`, job)
}

func (c *Cron) EveryDayWithJob(job contract.CronJob) {
	c.RunWithJob(`@daily`, job)
}

func (c *Cron) EveryWeeklyWithFunc(job func()) {
	c.RunWithFunc(`0 0 0 * * 0`, job)
}

func (c *Cron) EveryWeeklyWithJob(job contract.CronJob) {
	c.RunWithJob(`0 0 0 * * 0`, job)
}

func (c *Cron) WeeklyOnWithFunc(week, hour int, job func()) {
	panic("implement me")
}

func (c *Cron) WeeklyOnWithJob(week, hour int, job contract.CronJob) {
	panic("implement me")
}

func (c *Cron) Start() {
	c.cron.Start()
}

func (c *Cron) Stop() {
	c.cron.Stop()
}
