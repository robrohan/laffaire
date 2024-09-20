package sheets

import (
	"encoding/json"
	"fmt"
	"unicode"
)

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func loadSheet(dat []byte) ([]Sheet, error) {
	var spreadsheet []Sheet
	err := json.Unmarshal(dat, &spreadsheet)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	for k, v := range spreadsheet[0].Rows {
		// fmt.Printf("--> %v %v\n", k, v)

		if isInt(k) {
			// fmt.Printf("--> %v %v\n", k, v)
			fmt.Printf("-=-> %v \n\n", v)
			cells := v.(map[string]interface{})["cells"]
			fmt.Printf("-~~-> %v \n\n", cells)

			fmt.Printf("%v\n", len(cells))

			// cellRows := cells.(map[string]interface{})
			// fmt.Printf("~~> %v \n\n", cellRows)

		}

		// i, err := strconv.Atoi(k)
		// if err != nil {
		// 	fmt.Printf("<<%v>>", err)
		// 	break
		// }

		// switch v.(type) {
		// case int:
		// 	fmt.Printf("--> %v %v\n", k, v)
		// default:
		// 	fmt.Printf("d'know")
		// }
	}

	return spreadsheet, nil
}
