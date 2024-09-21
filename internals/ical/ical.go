package ical

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const (
	subject = iota
	startDate
	startTime
	endDate
	endTime
	allDayEvent
	description
	location
	private
)
const (
	calendar = "Plan"
	prodID   = "//Rob Rohan//Made up go code//EN"
	timeZone = "NZDT"
)

func FormatRecord(record []string, event *bytes.Buffer) error {
	uuid, err := NewId()
	if err != nil {
		panic("bad id gen")
	}

	if len(record)-1 < private {
		fmt.Printf("%v %v", len(record), private)
		return errors.New("bad record length")
	}

	formatDate := strings.Replace(record[startDate], "-", "", -1)
	if formatDate == "" {
		return errors.New("event missing date")
	}

	event.WriteString("BEGIN:VEVENT\r\n")
	fmt.Fprintf(event, "DTSTAMP:%vT000000Z\r\n", formatDate)
	fmt.Fprintf(event, "UID:ROHAN-%v\r\n", uuid)
	fmt.Fprintf(event, "DTSTART;VALUE=DATE:%v\r\n", formatDate)
	fmt.Fprintf(event, "DTEND;VALUE=DATE:%v\r\n", formatDate)
	fmt.Fprintf(event, "SUMMARY:%v\r\n", record[subject])
	fmt.Fprintf(event, "DESCRIPTION:%v\r\n", record[description])
	fmt.Fprintf(event, "CATEGORIES:%v\r\n", calendar)
	event.WriteString("END:VEVENT\r\n")

	return nil
}

func Prolog(prolog *bytes.Buffer) {
	prolog.WriteString("BEGIN:VCALENDAR\r\n")
	prolog.WriteString("VERSION:2.0\r\n")
	fmt.Fprintf(prolog, "X-WR-CALNAME:%v\r\n", calendar)
	fmt.Fprintf(prolog, "PRODID:%v\r\n", prodID)
	fmt.Fprintf(prolog, "X-WR-TIMEZONE:%v\r\n", timeZone)
	fmt.Fprintf(prolog, "X-WR-CALDESC:%v\r\n", calendar)
	prolog.WriteString("CALSCALE:GREGORIAN\r\n")
}

func Epilog(epilog *bytes.Buffer) {
	epilog.WriteString("END:VCALENDAR\r\n")
}

// Thanks internet
// https://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language#15134490
// Close enough for jazz.
func NewId() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u), nil
}
