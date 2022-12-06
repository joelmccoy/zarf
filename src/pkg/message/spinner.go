// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package message provides a rich set of functions for displaying messages to the user.
package message

import (
	"fmt"

	"github.com/pterm/pterm"
)

var activeSpinner *Spinner

// Spinner is a wrapper around pterm.SpinnerPrinter.
type Spinner struct {
	spinner   *pterm.SpinnerPrinter
	startText string
}

// NewProgressSpinner creates a new progress spinner.
func NewProgressSpinner(format string, a ...any) *Spinner {
	if activeSpinner != nil {
		Debug("Active spinner already exists")
		return activeSpinner
	}

	var spinner *pterm.SpinnerPrinter
	text := fmt.Sprintf(format, a...)
	if NoProgress {
		Info(text)
	} else {
		spinner, _ = pterm.DefaultSpinner.
			WithRemoveWhenDone(false).
			// Src: https://github.com/gernest/wow/blob/master/spin/spinners.go#L335
			WithSequence(`  ⠋ `, `  ⠙ `, `  ⠹ `, `  ⠸ `, `  ⠼ `, `  ⠴ `, `  ⠦ `, `  ⠧ `, `  ⠇ `, `  ⠏ `).
			Start(text)
	}

	activeSpinner = &Spinner{
		spinner:   spinner,
		startText: text,
	}

	return activeSpinner
}

// Write the given text to the spinner.
func (p *Spinner) Write(text []byte) (int, error) {
	size := len(text)
	if NoProgress {
		return size, nil
	}
	Debug(string(text))
	return len(text), nil
}

// Updatef updates the spinner text.
func (p *Spinner) Updatef(format string, a ...any) {
	if NoProgress {
		return
	}

	text := fmt.Sprintf(format, a...)
	p.spinner.UpdateText(text)
}

// Debugf prints a debug message.
func (p *Spinner) Debugf(format string, a ...any) {
	if logLevel >= DebugLevel {
		text := fmt.Sprintf("Debug: "+format, a)
		if NoProgress {
			Debug(text)
		} else {
			p.spinner.UpdateText(text)
		}
	}
}

// Stop the spinner.
func (p *Spinner) Stop() {
	if p.spinner != nil && p.spinner.IsActive {
		_ = p.spinner.Stop()
	}
	activeSpinner = nil
}

// Success prints a success message and stops the spinner.
func (p *Spinner) Success() {
	p.Successf(p.startText)
}

// Successf prints a success message with the spinner and stops it.
func (p *Spinner) Successf(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	if p.spinner != nil {
		p.spinner.Success(text)
		activeSpinner = nil
	} else {
		Info(text)
	}
}

// Warnf prints a warning message with the spinner.
func (p *Spinner) Warnf(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	if p.spinner != nil {
		p.spinner.Warning(text)
	} else {
		Warn(text)
	}
}

// Errorf prints an error message with the spinner.
func (p *Spinner) Errorf(err error, format string, a ...any) {
	p.Warnf(format, a...)
	Debug(err)
}

// Fatal calls message.Fatalf with the given error.
func (p *Spinner) Fatal(err error) {
	p.Fatalf(err, p.startText)
}

// Fatalf calls message.Fatalf with the given error and format.
func (p *Spinner) Fatalf(err error, format string, a ...any) {
	if p.spinner != nil {
		p.spinner.RemoveWhenDone = true
		_ = p.spinner.Stop()
		activeSpinner = nil
	}
	Fatalf(err, format, a...)
}