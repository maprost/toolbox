package tbsound

import "github.com/hegedustibor/htgo-tts"

// TextToSpeech needs mplayer installed
// apt install mplayer
//
// install htgo-tts
// go get "github.com/hegedustibor/htgo-tts"
func TextToSpeech(txt string) {
	speech := htgotts.Speech{Folder: "audio", Language: "de"}
	speech.Speak(txt)
}
