# mitto

mitto provides a simple API around application event handling

[![Build Status](https://github.com/xmidt-org/mitto/workflows/CI/badge.svg)](https://github.com/xmidt-org/mitto/actions)
[![codecov.io](http://codecov.io/github/xmidt-org/mitto/coverage.svg?branch=main)](http://codecov.io/github/xmidt-org/mitto?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/mitto)](https://goreportcard.com/report/github.com/xmidt-org/mitto)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/mitto/blob/main/LICENSE)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=xmidt-org_mitto&metric=alert_status)](https://sonarcloud.io/dashboard?id=xmidt-org_mitto)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/xmidt-org/mitto)](https://pkg.go.dev/github.com/xmidt-org/mitto)

## Summary

`mitto` (Latin verb meaning "to send") implements a simple way to manage a set of listeners and
dispatch application events to those listeners. `mitto` is primarily intended for code which
sends events to clients, where clients supply one or more listeners.

## Table of Contents

- [Usage](#usage)
- [Code of Conduct](#code-of-conduct)
- [Install](#install)
- [Contributing](#contributing)

## Usage

The easiest way to use `mitto` is to include it in a type that manages events and listeners:

```go
import github.com/xmidt-org/mitto

type Event {
    // stuff
}

type MyListener interface {
    OnEvent(Event)
}

type MyService struct {
    // could also use mitto.Listeners if concurrency is
    // managed by MyService or some other means.
    listeners mitto.SyncListeners[Event]
}

func (s *MyService) AddListeners(l ...MyListener) {
    mitto.AddListeners(&s.listeners, l...)
}

func (s *MyService) RemoveListeners(l ...MyListener) {
    mitto.RemoveListeners(&s.listeners, l...)
}

func (s *MyService) DoSomething() {
    // time to send an event:
    s.listeners.Send(Event{
        // stuff
    })
}
```

`mitto` allows closures to be used as event sinks.

```go
func (s *MyService) AddListenerFuncs(l ...func(Event)) {
    s.listeners.AddListenerFuncs(l...)
}
```

When using closures, remember that `golang` does not allow comparisons between functions. That means that you can't remove a listener closure later. For cases where a listener closure needs to be removed at some point, `mitto` provides `AsListener` to convert a closure into a comparable `Listener`.

```go
l := mitto.AsListener(func(Event) {})
s := new(MyService)
s.AddListeners(l)

// this will now work
s.RemoveListeners(l)
```

`AsListener` can also be used to adapt a different listener interface.

```go
type DifferentListener interface {
    OnStartEvent(Event)
}

func (s *MyService) AddDifferentListener(l DifferentListener) {
    s.listeners.AddListeners(
        mitto.AsListener(
            l.OnStartEvent,
        ),
    )
}
```

`mitto` also allows channel-based listeners. **Clients are responsible for creating and managing channels to avoid blocking.**

```go
func (s *MyService) AddListenerChans(ch ...chan<- Event) {
    s.listeners.AddListenerChans(ch...)
}
```

## Code of Conduct

This project and everyone participating in it are governed by the [XMiDT Code Of Conduct](https://xmidt.io/docs/community/code_of_conduct/). 
By participating, you agree to this Code.

## Install

go get github.com/xmidt-org/mitto@latest

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
