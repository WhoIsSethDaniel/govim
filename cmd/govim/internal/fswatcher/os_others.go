// +build !darwin

package fswatcher

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/tomb.v2"
)

type watcher struct {
	eventCh chan Event
	errCh   chan error
	mw      *fsnotify.Watcher
}

func New(gomodpath string, tomb *tomb.Tomb) (FSWatcher, error) {
	mw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create new watcher: %v", err)
	}

	eventCh := make(chan Event)
	tomb.Go(func() error {
		for {
			e, ok := <-mw.Events
			if !ok {
				break
			}
			switch e.Op {
			case fsnotify.Rename, fsnotify.Remove:
				eventCh <- Event{e.Name, OpRemoved}
			case fsnotify.Create, fsnotify.Chmod, fsnotify.Write:
				eventCh <- Event{e.Name, OpChanged}
			}
		}
		close(eventCh)
		return nil
	})

	return FSWatcher(&watcher{
		eventCh: eventCh,
		errCh:   mw.Errors,
		mw:      mw,
	}), nil
}

func (w *watcher) Add(path string) error {
	return w.mw.Add(path)
}

func (w *watcher) Remove(path string) error {
	return w.mw.Remove(path)
}

func (w *watcher) Close() error {
	return w.mw.Close()
}

func (w *watcher) Events() chan Event {
	return w.eventCh
}

func (w *watcher) Errors() chan error {
	return w.errCh
}
