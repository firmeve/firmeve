package contract

import "github.com/robfig/cron/v3"

type (
	CronJob interface {
		cron.Job
	}

	Cron interface {
		Start()

		Stop()

		// Perform tasks at a custom time
		RunWithFunc(spec string, job func())
		RunWithJob(spec string, job CronJob)

		// Execute task in specified seconds
		SecondsWithFunc(seconds int, job func())
		SecondsWithJob(seconds int, job CronJob)

		// Perform a task every second
		EverySecondWithFunc(job func())
		EverySecondWithJob(job CronJob)

		// Execute task in specified minutes
		MinutesWithFunc(minutes int, job func())
		MinutesWithJob(minutes int, job CronJob)

		// Perform a task every minute
		EveryMinuteWithFunc(job func())
		EveryMinuteWithJob(job CronJob)

		// Execute task in specified hours
		HoursWithFunc(hours int, job func())
		HoursWithJob(hours int, job CronJob)

		// Perform a task every hour
		EveryHourWithFunc(job func())
		EveryHourWithJob(job CronJob)

		// Run at a specified time every day
		DayAtWithFunc(hour int, job func())
		DayAtWithJob(hour int, job CronJob)

		// Every day at 0 o'clock
		EveryDayWithFunc(job func())
		EveryDayWithJob(job CronJob)

		// Perform the task at 0:00 every Sunday
		EveryWeeklyWithFunc(job func())
		EveryWeeklyWithJob(job CronJob)

		// Perform a task at N points of N every week
		WeeklyOnWithFunc(week, hour int, job func())
		WeeklyOnWithJob(week, hour int, job CronJob)
	}
)
