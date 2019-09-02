package fswatcher

type Event struct {
	Path string
	Op   Op
}

type Op string

const (
	OpChanged Op = "changed"
	OpRemoved Op = "removed"
)

type FSWatcher interface {
	Add(path string) error
	Remove(path string) error
	Close() error
	Events() chan Event
	Errors() chan error
}
