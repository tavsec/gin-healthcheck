package checks

type Check interface {
	Pass() bool
	Name() string
}
