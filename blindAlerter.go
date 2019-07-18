package poker

import (
	"fmt"
	"io"
	"time"
)

//BlindAlerter interface sets a time and amount
type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amount int, alertsDestination io.Writer)
}

//BlindAlerterFunc converts the interface into a func
type BlindAlerterFunc func(duration time.Duration, amount int, alertsDestination io.Writer)

// ScheduledAlertAt sets the duration and amount via the Func
func (a BlindAlerterFunc) ScheduledAlertAt(duration time.Duration, amount int, alertsDestination io.Writer) {
	a(duration, amount, alertsDestination)
}

// Alerter applies the duration and amount to the stdout
func Alerter(duration time.Duration, amount int, alertsDestination io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(alertsDestination, "Blind is now %d\n", amount)
	})
}
