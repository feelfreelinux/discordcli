package core

import (
	"strconv"

	"github.com/aybabtme/rgbterm"
)

/*
GetColoredNick returns rgb escaped nick with given color in dec format
*/
func GetColoredNick(nick string, color int) string {
	hex := strconv.FormatInt(int64(color), 16)
	if len(hex) == 6 {
		r, _ := strconv.ParseInt(hex[0:2], 16, 0)
		g, _ := strconv.ParseInt(hex[2:4], 16, 0)
		b, _ := strconv.ParseInt(hex[4:6], 16, 0)
		return rgbterm.FgString(nick, uint8(r), uint8(g), uint8(b))
	}
	return nick

}
