package gl

import (
	"time"
)

type AccessLevel int
type VisibilityLevel int
type NotificationLevel int
type State string

const (
	Guest     = AccessLevel(10)
	Reporter  = AccessLevel(20)
	Developer = AccessLevel(30)
	Master    = AccessLevel(40)
	Owner     = AccessLevel(50)

	Private  = VisibilityLevel(0)
	Internal = VisibilityLevel(10)
	Public   = VisibilityLevel(20)

	NotificationDisabled      = NotificationLevel(0)
	NotificationParticipating = NotificationLevel(1)
	NotificationWatch         = NotificationLevel(2)
	NotificationGlobal        = NotificationLevel(3)

	Active  State = "active"
	Blocked       = "blocked"
)

const (
	dateLayout        = "2006-01-02T15:04:05-07:00"
	jsonDateLayout    = "\"2006-01-02\""
	txtjsonDateLayout = "2006-01-02"
)

// JsonDate represents a date of the form "YYYY-MM-DD" in a json string
type JsonDate struct {
	time.Time
}

func (j *JsonDate) format(layout string) string {
	return j.Time.Format(layout)
}
func parseJsonDate(layout, val string) (time.Time, error) {
	return time.Parse(layout, val)
}
func (j *JsonDate) MarshalText() ([]byte, error) {
	return []byte(j.format(txtjsonDateLayout)), nil
}

func (j *JsonDate) MarshalJSON() ([]byte, error) {
	return []byte(j.format(jsonDateLayout)), nil
}
func (j *JsonDate) UnmarshalJSON(data []byte) (err error) {
	t, err := parseJsonDate(jsonDateLayout, string(data))
	if err != nil {
		return err
	}
	j.Time = t
	return nil
}
func (j *JsonDate) UnmarshalText(data []byte) (err error) {
	t, err := parseJsonDate(txtjsonDateLayout, string(data))
	if err != nil {
		return err
	}
	j.Time = t
	return nil
}
