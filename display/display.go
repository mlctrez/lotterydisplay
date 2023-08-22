package display

import (
	"sort"
)

// Group represents a group of seven segment displays
type Group struct {
	Digit     [8]byte
	Intensity byte
}

// Clear sets all characters in the display to blank
func (dg *Group) Clear() {
	dg.WriteString("        ")
}

// WriteChar shifts characters 1-7 one to the left and sets character 7 to the provided value
func (dg *Group) WriteChar(char string) {
	for i := 0; i < 7; i++ {
		dg.Digit[i] = dg.Digit[i+1]
	}
	dg.SetDigit(7, char)
}

// WriteString writes each character in the string to the display using WriteChar
func (dg *Group) WriteString(data string) {
	for i := range data {
		dg.WriteChar(data[i : i+1])
	}
}

// Packet copies the internal display structure into a byte array suitable for sending over UDP
func (dg *Group) Packet() []byte {
	b := make([]byte, 9)
	for i, v := range dg.Digit {
		b[i] = v
	}
	b[8] = dg.Intensity
	return b
}

// SetDigit set segments A-G according to dMap
func (dg *Group) SetDigit(num int, char string) {
	// clear A-G, preserving decimal point
	dg.Digit[num] &= sMap["DP"]
	if v, ok := cMap[char]; ok {
		for i := range v {
			dg.Digit[num] |= sMap[v[i:i+1]]
		}
	}
}

/*
  *A*
 F   B
  *G*
 E   C
  *D* DP
*/

// sMap defines segment to binary bit for the display
var sMap = map[string]byte{
	"DP": 0x80, "A": 0x40, "B": 0x20, "C": 0x10, "D": 0x08, "E": 0x04, "F": 0x02, "G": 0x01,
}

func SupportedCharacters() []string {
	var result []string

	for k := range cMap {
		result = append(result, k)
	}
	sort.Strings(result)

	return result
}

// cMap maps characters to segments
var cMap = map[string]string{
	"0": "ABCDEF",
	"1": "BC",
	"2": "ABGED",
	"3": "ABGCD",
	"4": "FGBC",
	"5": "AFGCD",
	"6": "AFEDCG",
	"7": "ABC",
	"8": "AFGCDEB",
	"9": "GFABCD",
	"A": "EFABCG",
	"b": "FEGCD",
	"C": "AFED",
	"c": "GED",
	"d": "BCGED",
	"E": "AFGED",
	"F": "FEAG",
	"G": "AFEDCG",
	"H": "FEBCG",
	"h": "FEGC",
	"I": "FE",
	"i": "C",
	"J": "BCDE",
	"L": "FED",
	"n": "EGC",
	"O": "ABCDEF",
	"o": "GEDC",
	"P": "FEABG",
	"q": "AFGBC",
	"r": "EG",
	"S": "AFGCD",
	"t": "FEDG",
	"U": "FEDCB",
	"u": "EDC",
	"y": "FGBCD",
	"-": "G",
	"_": "D",
	"*": "AFGB",
	"?": "ABGE",
	"]": "ABCD",
	"[": "AFED",
}

// SetDecimalPoint turns the decimal point on or off at digit num
func (dg *Group) SetDecimalPoint(num int, point bool) {
	if point {
		dg.Digit[num] |= sMap["DP"]
	} else {
		dg.Digit[num] &^= sMap["DP"]
	}
}
