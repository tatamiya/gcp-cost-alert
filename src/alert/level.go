package alert

type AlertLevel int

const (
	Unexpected AlertLevel = iota
	Low
	Middle
	High
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
