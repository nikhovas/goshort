package goshort

import (
	"flag"
	"goshort/utils"
)

func main() {
	useColor := flag.String("config", "", "display colorized output")
	flag.Parse()

	utils.SetupViper(*useColor)
	AppObject = App{}
	AppObject.Initialize()
	AppObject.Run()
}
