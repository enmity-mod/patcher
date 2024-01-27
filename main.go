package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    false,
	ReportTimestamp: true,
	TimeFormat:      time.TimeOnly,
	Level:           log.DebugLevel,
	Prefix:          "Patcher",
})

var (
	info      map[string]interface{}
	directory string
	assets    string
	ipa       string
)

func main() {
	app := &cli.App{
		Name:  "patcher-ios",
		Usage: "Patches the Discord IPA with icons, utilities and features to ease usability.",
		Action: func(context *cli.Context) error {
			ipa = context.Args().Get(0)

			if ipa == "" {
				logger.Error("Please provide a path to the IPA.")
				os.Exit(1)
			}

			logger.Infof("Requested IPA patch for \"%s\"", ipa)

			extract()
			loadInfo()

			setReactNavigationName()
			setSupportedDevices()
			setFileAccess()
			setAppName()
			setIcons()

			saveInfo()
			archive()

			exit()
			return nil
		},
	}

	assets = os.TempDir()

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
