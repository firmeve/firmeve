package kernel

type (
	Error struct {
		err     error
		code    int
		message string
	}
)
