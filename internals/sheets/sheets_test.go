package sheets

import (
	"fmt"
	"os"
	"testing"
)

func TestParseSheet(t *testing.T) {
	dat, err := os.ReadFile("../../test_data/s1.json")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// fmt.Printf("%v\n", string(dat))

	spreadsheet, err := loadSheet(dat)
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
		return
	}

	fmt.Printf("%v\n", spreadsheet)
}
