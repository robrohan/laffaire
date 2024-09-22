package ical

import (
	"bytes"
	"fmt"
)

func Prolog(prolog *bytes.Buffer, calendar string, prodId string, timeZone string) {
	prolog.WriteString("BEGIN:VCALENDAR\r\n")
	prolog.WriteString("VERSION:2.0\r\n")
	fmt.Fprintf(prolog, "X-WR-CALNAME:%v\r\n", calendar)
	fmt.Fprintf(prolog, "PRODID:%v\r\n", prodId)
	fmt.Fprintf(prolog, "X-WR-TIMEZONE:%v\r\n", timeZone)
	fmt.Fprintf(prolog, "X-WR-CALDESC:%v\r\n", calendar)
	prolog.WriteString("CALSCALE:GREGORIAN\r\n")
}

func Epilog(epilog *bytes.Buffer) {
	epilog.WriteString("END:VCALENDAR\r\n")
}
