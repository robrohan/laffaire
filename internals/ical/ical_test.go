package ical

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrologContainsRequiredFields(t *testing.T) {
	var buf bytes.Buffer
	Prolog(&buf, "My Calendar", "-//Laffaire/My Calendar//EN", "America/New_York")
	out := buf.String()

	required := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"CALSCALE:GREGORIAN",
		"X-WR-CALNAME:My Calendar",
		"PRODID:-//Laffaire/My Calendar//EN",
		"X-WR-TIMEZONE:America/New_York",
		"X-WR-CALDESC:My Calendar",
	}
	for _, want := range required {
		if !strings.Contains(out, want) {
			t.Errorf("prolog missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestPrologProdIdHasLeadingDash(t *testing.T) {
	var buf bytes.Buffer
	Prolog(&buf, "Test", "-//Laffaire/Test//EN", "UTC")
	out := buf.String()

	if !strings.Contains(out, "PRODID:-//") {
		t.Errorf("PRODID should start with '-//' for a valid ISO 9070 FPI\ngot:\n%s", out)
	}
}

func TestPrologCalendarNameIsUsed(t *testing.T) {
	var buf bytes.Buffer
	Prolog(&buf, "Summer Tour", "-//Laffaire/Summer Tour//EN", "UTC")
	out := buf.String()

	if !strings.Contains(out, "X-WR-CALNAME:Summer Tour") {
		t.Errorf("X-WR-CALNAME should equal the calendar name\ngot:\n%s", out)
	}
	if !strings.Contains(out, "PRODID:-//Laffaire/Summer Tour//EN") {
		t.Errorf("PRODID should contain the calendar name\ngot:\n%s", out)
	}
}

// Ensure that a title containing "//" is sanitised before being embedded in
// the PRODID so it does not look like a spurious FPI separator.
func TestPrologProdIdSanitisedDoubleSlash(t *testing.T) {
	rawTitle := "AC//DC Tour"
	prodIdName := strings.ReplaceAll(rawTitle, "//", "-")
	prodId := "-//Laffaire/" + prodIdName + "//EN"

	var buf bytes.Buffer
	Prolog(&buf, rawTitle, prodId, "UTC")
	out := buf.String()

	// X-WR-CALNAME keeps the original title
	if !strings.Contains(out, "X-WR-CALNAME:AC//DC Tour") {
		t.Errorf("X-WR-CALNAME should be the raw title\ngot:\n%s", out)
	}
	// PRODID must not contain an extra "//" inside the owner/description segment
	if strings.Contains(out, "PRODID:-//Laffaire/AC//DC") {
		t.Errorf("PRODID should have '//' sanitised out of the calendar name\ngot:\n%s", out)
	}
	if !strings.Contains(out, "PRODID:-//Laffaire/AC-DC Tour//EN") {
		t.Errorf("PRODID should contain sanitised name\ngot:\n%s", out)
	}
}

func TestEpilog(t *testing.T) {
	var buf bytes.Buffer
	Epilog(&buf)
	if !strings.Contains(buf.String(), "END:VCALENDAR") {
		t.Errorf("epilog missing END:VCALENDAR\ngot: %s", buf.String())
	}
}
