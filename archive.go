package main

import (
	"compress/flate"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

func extract() {
	logger.Debugf("Attempting to extract \"%s\"", ipa)
	format := archiver.Zip{}
	directory = fileNameWithoutExtension(filepath.Base(ipa))

	if _, err := os.Stat(ipa); err != nil {
		logger.Errorf("Couldn't find \"%s\". Does it exist?", ipa)
		exit()
	}

	if _, err := os.Stat(directory); err == nil {
		logger.Debug("Detected previously extracted directory, cleaning it up...")

		err := os.RemoveAll(directory)
		if err != nil {
			logger.Errorf("Failed to clean up previously extracted directory: %s", err)
			exit()
		}

		logger.Info("Previously extracted directory cleaned up. ")
	}

	err := format.Unarchive(ipa, directory)
	if err != nil {
		logger.Errorf("Failed to extract %s: **%v**", ipa, err)
		os.Exit(1)
	}

	logger.Infof("Successfully extracted to \"%s\"", directory)
}

func archive() {
	logger.Debugf("Attempting to archive \"%s\"", directory)

	format := archiver.Zip{CompressionLevel: flate.BestCompression}
	zip := directory + ".zip"

	if _, err := os.Stat(zip); err == nil {
		logger.Debug("Detected previous archive, cleaning it up...")

		err := os.Remove(zip)
		if err != nil {
			logger.Errorf("Failed to clean up previous archive: %s", err)
			exit()
		}

		logger.Info("Previous archive cleaned up.")
	}

	logger.Infof("Archiving \"%s\" to \"%s\"", directory, zip)
	err := format.Archive([]string{filepath.Join(directory, "Payload")}, zip)
	if err != nil {
		logger.Errorf("Failed to archive \"%s\": %v", zip, err)
		exit()
	}

	if _, err := os.Stat("Enmity.ipa"); err == nil {
		logger.Debug("Detected previous Enmity IPA, cleaning it up...")

		err := os.Remove("Enmity.ipa")
		if err != nil {
			logger.Errorf("Failed to clean up previous Enmity IPA: %s", err)
			exit()
		}

		logger.Info("Previous Enmity IPA cleaned up.")
	}

	err = os.Rename(zip, "Enmity.ipa")
	if err != nil {
		logger.Errorf("Failed to rename \"%s\": %v", zip, err)
		exit()
	}

	logger.Infof("Successfully archived \"%s\" to \"Enmity.ipa\"", zip)
}
