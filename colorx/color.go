package colorx

import (
	"fmt"
	"strings"
)

type Color string

const (
	// Reset
	ColorOff = Color("\033[0m") // Reset Color

	// Regular Colors
	Black  = Color("\033[0;30m") // Black
	Red    = Color("\033[0;31m") // Red
	Green  = Color("\033[0;32m") // Green
	Yellow = Color("\033[0;33m") // Yellow
	Blue   = Color("\033[0;34m") // Blue
	Purple = Color("\033[0;35m") // Purple
	Cyan   = Color("\033[0;36m") // Cyan
	White  = Color("\033[0;37m") // White

	// Bold
	BoldBlack  = Color("\033[1;30m") // Black
	BoldRed    = Color("\033[1;31m") // Red
	BoldGreen  = Color("\033[1;32m") // Green
	BoldYellow = Color("\033[1;33m") // Yellow
	BoldBlue   = Color("\033[1;34m") // Blue
	BoldPurple = Color("\033[1;35m") // Purple
	BoldCyan   = Color("\033[1;36m") // Cyan
	BoldWhite  = Color("\033[1;37m") // White

	// Underline
	UnderlineBlack  = Color("\033[4;30m") // Black
	UnderlineRed    = Color("\033[4;31m") // Red
	UnderlineGreen  = Color("\033[4;32m") // Green
	UnderlineYellow = Color("\033[4;33m") // Yellow
	UnderlineBlue   = Color("\033[4;34m") // Blue
	UnderlinePurple = Color("\033[4;35m") // Purple
	UnderlineCyan   = Color("\033[4;36m") // Cyan
	UnderlineWhite  = Color("\033[4;37m") // White

	// Background
	BackgroundBlack  = Color("\033[40m") // Black
	BackgroundRed    = Color("\033[41m") // Red
	BackgroundGreen  = Color("\033[42m") // Green
	BackgroundYellow = Color("\033[43m") // Yellow
	BackgroundBlue   = Color("\033[44m") // Blue
	BackgroundPurple = Color("\033[45m") // Purple
	BackgroundCyan   = Color("\033[46m") // Cyan
	BackgroundWhite  = Color("\033[47m") // White

	// High Intensty
	LightBlack  = Color("\033[0;90m") // Black
	LightRed    = Color("\033[0;91m") // Red
	LightGreen  = Color("\033[0;92m") // Green
	LightYellow = Color("\033[0;93m") // Yellow
	LightBlue   = Color("\033[0;94m") // Blue
	LightPurple = Color("\033[0;95m") // Purple
	LightCyan   = Color("\033[0;96m") // Cyan
	LightWhite  = Color("\033[0;97m") // White

	// Bold High Intensty
	BoldLightBlack  = Color("\033[1;90m") // Black
	BoldLightRed    = Color("\033[1;91m") // Red
	BoldLightGreen  = Color("\033[1;92m") // Green
	BoldLightYellow = Color("\033[1;93m") // Yellow
	BoldLightBlue   = Color("\033[1;94m") // Blue
	BoldLightPurple = Color("\033[1;95m") // Purple
	BoldLightCyan   = Color("\033[1;96m") // Cyan
	BoldLightWhite  = Color("\033[1;97m") // White

	// High Intensty backgrounds
	LightBackgroundBlack  = Color("\033[0;100m") // Black
	LightBackgroundRed    = Color("\033[0;101m") // Red
	LightBackgroundGreen  = Color("\033[0;102m") // Green
	LightBackgroundYellow = Color("\033[0;103m") // Yellow
	LightBackgroundBlue   = Color("\033[0;104m") // Blue
	LightBackgroundPurple = Color("\033[10;95m") // Purple
	LightBackgroundCyan   = Color("\033[0;106m") // Cyan
	LightBackgroundWhite  = Color("\033[0;107m") // White
)

var (
	AllColors = []Color{
		ColorOff,
		Black, Red, Green, Yellow, Blue, Purple, Cyan, White,
		BoldBlack, BoldRed, BoldGreen, BoldYellow, BoldBlue, BoldPurple, BoldCyan, BoldWhite,
		UnderlineBlack, UnderlineRed, UnderlineGreen, UnderlineYellow, UnderlineBlue, UnderlinePurple, UnderlineCyan, UnderlineWhite,
		BackgroundBlack, BackgroundRed, BackgroundGreen, BackgroundYellow, BackgroundBlue, BackgroundPurple, BackgroundCyan, BackgroundWhite,
		LightBlack, LightRed, LightGreen, LightYellow, LightBlue, LightPurple, LightCyan, LightWhite,
		BoldLightBlack, BoldLightRed, BoldLightGreen, BoldLightYellow, BoldLightBlue, BoldLightPurple, BoldLightCyan, BoldLightWhite,
		LightBackgroundBlack, LightBackgroundRed, LightBackgroundGreen, LightBackgroundYellow, LightBackgroundBlue, LightBackgroundPurple, LightBackgroundCyan, LightBackgroundWhite,
	}
)

func Text(txt string, color Color) string {
	return fmt.Sprint(color, txt, ColorOff)
}

func Textf(color Color, txt string, args ...interface{}) string {
	return Text(fmt.Sprintf(txt, args...), color)
}

func Remove(msg string) string {
	for _, color := range AllColors {
		msg = strings.ReplaceAll(msg, string(color), "")
	}
	return msg
}

func Removes(msg []string) []string {
	for i, txt := range msg {
		msg[i] = Remove(txt)
	}
	return msg
}
