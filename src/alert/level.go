package alert

// AlertLevel defines the level of the alert
// based on the exeeded threshold.
// The headline of a notification message changes with this level.
type AlertLevel int

const (
	Unexpected AlertLevel = iota
	Low                   // 0.2 <= threshold < 0.5
	Middle                // 0.5 <= threshold < 1.0
	High                  // threshold >= 1.0
)

func newAlertLevel(threshold float64) AlertLevel {
	switch {
	case threshold >= 0.2 && threshold < 0.5:
		return Low
	case threshold >= 0.5 && threshold < 1.0:
		return Middle
	case threshold >= 1.0:
		return High
	default:
		return Unexpected
	}
}

func (l *AlertLevel) headline() string {
	switch *l {
	case Low:
		return ":warning: 注意 :warning:"
	case Middle:
		return ":rotating_light: 警報 :rotating_light:"
	case High:
		return ":fire::fire::fire: 警告 :fire::fire::fire:"
	default:
		return ":asyncparrot:"
	}
}
