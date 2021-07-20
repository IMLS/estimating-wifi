package analysis

type Counter struct {
	Patrons          int
	Devices          int
	Transients       int
	PatronMinutes    int
	DeviceMinutes    int
	TransientMinutes int
}

func NewCounter(minMinutes int, maxMinutes int) *Counter {
	patronMinMins = float64(minMinutes)
	patronMaxMins = float64(maxMinutes)
	return &Counter{0, 0, 0, 0, 0, 0}
}

func (c *Counter) Add(field int, minutes int) {
	switch field {
	case Patron:
		c.Patrons += 1
		c.PatronMinutes += minutes
	case Device:
		c.Devices += 1
		c.DeviceMinutes += minutes
	case Transient:
		c.Transients += 1
		c.TransientMinutes += minutes
	}
}
