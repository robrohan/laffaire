package sheets

// type CellStyle struct {
// 	Width int32 `json:"width"`
// 	Style int16 `json:"style"`
// }

type Cell struct {
	Text string `json:"text"`
}

// type Cols struct {
// 	// Length int16                `json:"len"`
// 	Cols map[string]CellStyle
// }

type Cells map[int16]Cell

type Sheet struct {
	Name   string                 `json:"name"`
	Freeze string                 `json:"freeze"`
	Cols   map[string]interface{} `json:"cols"`
	Rows   map[string]interface{} `json:"rows"`
}

// type Spreadsheet struct {
// 	Sheet []Sheet
// }
