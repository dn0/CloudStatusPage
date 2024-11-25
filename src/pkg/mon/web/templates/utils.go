package templates

import (
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

const (
	thinSpace             = "&thinsp;"
	timeDay               = 24 * time.Hour
	timeFormatWithSeconds = "yyyy-LL-dd HH:mm:ss ZZZZ"
	timeFormatWithoutTZ   = "yyyy-LL-dd HH:mm"
	// Hardcoded in base.templ.
	// timeFormatDefault = "yyyy-LL-dd HH:mm ZZZZ".
)

func ShortUUID(id string) templ.Component {
	return templ.Raw(`<span class="text-nowrap">` + id[:13] + `</span>`)
}

func Timestamp(t time.Time) templ.Component {
	return timestamp(t, "")
}

func TimestampWithSeconds(t time.Time) templ.Component {
	return timestamp(t, timeFormatWithSeconds)
}

func TimestampWithoutTZ(t time.Time) templ.Component {
	return timestamp(t, timeFormatWithoutTZ)
}

func timestamp(tim time.Time, fmt string) templ.Component {
	sb := new(strings.Builder)
	sb.WriteString(`<span class="timestamp-raw"`)

	if fmt == "" {
		sb.WriteString(">")
	} else {
		sb.WriteString(` data-timestamp_fmt="`)
		sb.WriteString(fmt)
		sb.WriteString(`">`)
	}

	sb.WriteString(tim.In(time.UTC).Format(time.RFC3339))
	sb.WriteString(`</span>`)

	return templ.Raw(sb.String())
}

func DurationMilliseconds(d, round time.Duration, decimalPlaces int) templ.Component {
	dms := strconv.FormatFloat(float64(d.Round(round))/1e6, 'f', decimalPlaces, 64)
	return templ.Raw(`<span class="text-nowrap">` + dms + thinSpace + "ms</span>")
}

//nolint:varnamelen,durationcheck,funlen,cyclop,mnd // Multiplications of durations are OK here.
func Duration(d, round time.Duration) templ.Component {
	if d < 2*round {
		switch round {
		case time.Hour:
			round = time.Minute
		case time.Minute:
			round = time.Second
		case time.Second:
			round = time.Millisecond
		case time.Millisecond:
			round = 10 * time.Microsecond
		case time.Microsecond:
			round = 10 * time.Nanosecond
		case time.Nanosecond:
			round = 0
		}
	}
	d = d.Round(round)

	sb := new(strings.Builder)
	sb.WriteString(`<span class="text-nowrap">`)
	empty := true

	if d >= timeDay {
		days := d / timeDay
		d -= days * timeDay
		sb.WriteString(strconv.Itoa(int(days)))
		sb.WriteString("d")
		empty = false
	}

	if d >= time.Hour {
		hours := d / time.Hour
		d -= hours * time.Hour
		if !empty {
			sb.WriteString(thinSpace)
		}
		sb.WriteString(strconv.Itoa(int(hours)))
		sb.WriteString("h")
		empty = false
	}

	if d >= time.Minute {
		mins := d / time.Minute
		d -= mins * time.Minute
		if !empty {
			sb.WriteString(thinSpace)
		}
		sb.WriteString(strconv.Itoa(int(mins)))
		sb.WriteString("m")
		empty = false
	}

	if empty {
		sb.WriteString(d.String())
	} else if round < time.Minute && d > 0 {
		sb.WriteString(thinSpace)
		sb.WriteString(d.String())
	}
	sb.WriteString(`</span>`)

	return templ.Raw(sb.String())
}
