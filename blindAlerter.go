package poker

import (
	"time"
	"fmt"
	"io"
)

//BlindAlerter interface sets a time and amount
type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amount int, to io.Writer)
}

//BlindAlerterFunc converts the interface into a func
type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

// ScheduledAlertAt sets the duration and amount via the Func
func (a BlindAlerterFunc) ScheduledAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}

// Alerter applies the duration and amount to the stdout
func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}