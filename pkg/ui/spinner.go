package ui

import (
	"fmt"
	"time"
)

// Spinner represents a loading spinner
type Spinner struct {
	frames []string
	delay  time.Duration
	active bool
	done   chan bool
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:  100 * time.Millisecond,
		done:   make(chan bool),
	}
}

// NewDotSpinner creates a dot-based spinner
func NewDotSpinner() *Spinner {
	return &Spinner{
		frames: []string{"   ", ".  ", ".. ", "..."},
		delay:  500 * time.Millisecond,
		done:   make(chan bool),
	}
}

// NewSquareSpinner creates a square-based spinner
func NewSquareSpinner() *Spinner {
	return &Spinner{
		frames: []string{"◐", "◓", "◑", "◒"},
		delay:  200 * time.Millisecond,
		done:   make(chan bool),
	}
}

// Start starts the spinner with a message
func (s *Spinner) Start(message string) {
	if s.active {
		return
	}
	s.active = true
	
	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				fmt.Printf("\r%s %s", s.frames[i%len(s.frames)], message)
				i++
				time.Sleep(s.delay)
			}
		}
	}()
}

// Stop stops the spinner and clears the line
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	s.done <- true
	fmt.Print("\r\033[K") // Clear the line
}

// Update updates the spinner message
func (s *Spinner) Update(message string) {
	if s.active {
		fmt.Printf("\r%s", message)
	}
}