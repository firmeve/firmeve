package kernel

type (
	SchedulerRecoverError struct {
		Message string
		Recover interface{}
	}
)

func (s SchedulerRecoverError) Error() string  {
	return s.Message
}
