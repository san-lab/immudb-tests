package color

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const BLACK = "black"
const BLUE = "blue"
const CYAN = "cyan"
const GREEN = "green"
const MAGENTA = "magenta"
const RED = "red"
const WHITE = "white"
const YELLOW = "yellow"
const BOLD = "bold"
const ITALIC = "italic"
const FAINT = "faint"
const UNDERLINE = "underline"

func CPrintln(color, format string, a ...interface{}) {
	switch color {
	case BLACK:
		fmt.Println(promptui.Styler(promptui.FGBlack)(fmt.Sprintf(format, a...)))
	case BLUE:
		fmt.Println(promptui.Styler(promptui.FGBlue)(fmt.Sprintf(format, a...)))
	case CYAN:
		fmt.Println(promptui.Styler(promptui.FGCyan)(fmt.Sprintf(format, a...)))
	case GREEN:
		fmt.Println(promptui.Styler(promptui.FGGreen)(fmt.Sprintf(format, a...)))
	case MAGENTA:
		fmt.Println(promptui.Styler(promptui.FGMagenta)(fmt.Sprintf(format, a...)))
	case RED:
		fmt.Println(promptui.Styler(promptui.FGRed)(fmt.Sprintf(format, a...)))
	case WHITE:
		fmt.Println(promptui.Styler(promptui.FGWhite)(fmt.Sprintf(format, a...)))
	case YELLOW:
		fmt.Println(promptui.Styler(promptui.FGYellow)(fmt.Sprintf(format, a...)))
	case BOLD:
		fmt.Println(promptui.Styler(promptui.FGBold)(fmt.Sprintf(format, a...)))
	case ITALIC:
		fmt.Println(promptui.Styler(promptui.FGItalic)(fmt.Sprintf(format, a...)))
	case FAINT:
		fmt.Println(promptui.Styler(promptui.FGFaint)(fmt.Sprintf(format, a...)))
	case UNDERLINE:
		fmt.Println(promptui.Styler(promptui.FGUnderline)(fmt.Sprintf(format, a...)))
	default:
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

func CPrintf(color, format string, a ...interface{}) {
	switch color {
	case BLACK:
		fmt.Print(promptui.Styler(promptui.FGBlack)(fmt.Sprintf(format, a...)))
	case BLUE:
		fmt.Print(promptui.Styler(promptui.FGBlue)(fmt.Sprintf(format, a...)))
	case CYAN:
		fmt.Print(promptui.Styler(promptui.FGCyan)(fmt.Sprintf(format, a...)))
	case GREEN:
		fmt.Print(promptui.Styler(promptui.FGGreen)(fmt.Sprintf(format, a...)))
	case MAGENTA:
		fmt.Print(promptui.Styler(promptui.FGMagenta)(fmt.Sprintf(format, a...)))
	case RED:
		fmt.Print(promptui.Styler(promptui.FGRed)(fmt.Sprintf(format, a...)))
	case WHITE:
		fmt.Print(promptui.Styler(promptui.FGWhite)(fmt.Sprintf(format, a...)))
	case YELLOW:
		fmt.Print(promptui.Styler(promptui.FGYellow)(fmt.Sprintf(format, a...)))
	case BOLD:
		fmt.Print(promptui.Styler(promptui.FGBold)(fmt.Sprintf(format, a...)))
	case ITALIC:
		fmt.Print(promptui.Styler(promptui.FGItalic)(fmt.Sprintf(format, a...)))
	case FAINT:
		fmt.Print(promptui.Styler(promptui.FGFaint)(fmt.Sprintf(format, a...)))
	case UNDERLINE:
		fmt.Print(promptui.Styler(promptui.FGUnderline)(fmt.Sprintf(format, a...)))
	default:
		fmt.Printf(format, a...)
	}
}

func Shorten(long string, length int) string {
	if len(long) == 0 {
		return "_"
	}
	if len(long) <= length {
		return long
	}
	return long[0:length] + "..."
}
