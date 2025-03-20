package check

type Check interface {
	Name() string
	PassedMessage() string
	FailedMessage() string
	Run() error
	Passed() bool
	IsRunnable() bool
	UUID() string
	Status() string
	RequiresRoot() bool
}
