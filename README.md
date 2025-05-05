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

### Basic

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
    // could also use mitto.Set if concurrency is
    // managed by MyService or some other means.
    listeners mitto.SyncSet[Event]
}

func (s *MyService) Add(l ...MyListener) {
    mitto.Add(&s.listeners, l...)
}

func (s *MyService) Remove(l ...MyListener) {
    mitto.Remove(&s.listeners, l...)
}

func (s *MyService) DoSomething() {
    // time to send an event:
    s.listeners.Send(Event{
        // stuff
    })
}
```

### AsListener

`mitto` provides `AsListener` to convert other common types into listeners.

```go
f := func(event Event) { /* ... */ }
ch := make(chan Event, 10)

s.listeners.Add(
    mitto.AsListener[Event](f),
    mitto.AsListener[Event](ch),
)
```

#### Important note on channels

A client is responsible for ensuring that a channel is properly managed to reduce or avoid blocking. In addition, a client must remove the channel listener *before* closing the channel, otherwise `Send` may panic.

```go
ch := make(chan Event, 10)
l := AsListener[Event](ch)
s.listeners.Add(l)

// remove the listener BEFORE closing the channel
s.listeners.Remove(l)
close(ch)
```

### Adapting custom Listener interfaces

`AsListener` can also be used to adapt a different listener interface or type.

```go
type DifferentListener interface {
    OnStartEvent(Event)
}

func (s *MyService) AddDifferentListener(l DifferentListener) {
    s.listeners.Add(
        mitto.AsListener[Event](
            l.OnStartEvent,
        ),
    )
}
```

## Code of Conduct

This project and everyone participating in it are governed by the [XMiDT Code Of Conduct](https://xmidt.io/docs/community/code_of_conduct/). 
By participating, you agree to this Code.

## Install

go get github.com/xmidt-org/mitto@latest

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
