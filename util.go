package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"howett.net/plist"
)

func fileNameWithoutExtension(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}

func exit() {
	logger.Debug("Cleaning up...")

	if _, e := os.Stat(directory); e == nil {
		if e := os.RemoveAll(directory); e != nil {
			logger.Errorf("Failed to clean up extracted directory: %v", e)
		}
	}

	if _, e := os.Stat(assets); e == nil {
		defer func() {
			if e := os.RemoveAll(assets); e != nil {
				logger.Errorf("Failed to clean up temporary assets directory: %v", e)
			}
		}()
	}

	logger.Info("Cleaned up.")

	os.Exit(0)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}

	return
}

func loadInfo() {
	if info != nil {
		return
	}

	path := filepath.Join(directory, "Payload", "Discord.app", "Info.plist")
	file, err := os.Open(path)

	if err != nil {
		logger.Error("Couldn't find Info.plist. Is the provided zip an IPA file?")
		exit()
	}

	decoder := plist.NewDecoder(file)
	if err := decoder.Decode(&info); err != nil {
		logger.Error("Couldn't find Info.plist. Is the provided zip an IPA file?")
		exit()
	}
}

func saveInfo() {
	path := filepath.Join(directory, "Payload", "Discord.app", "Info.plist")
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0600)

	if err != nil {
		logger.Errorf("Failed to open Info.plist for saving: %v", err)
		exit()
	}

	logger.Debug("Saving Info.plist data...")
	encoder := plist.NewEncoder(file)
	err = encoder.Encode(info)

	if err != nil {
		logger.Errorf("Failed to save Info.plist. %v", err)
		exit()
	}

	logger.Infof("Saved Info.plist data.")
}

func download(url string, path string) {
	out, err := os.Create(path)

	if err != nil {
		logger.Errorf("Failed to pre-write file at %s.", path)
		exit()
	}

	res, err := http.Get(url)

	if err != nil {
		logger.Errorf("Failed to download %s to %s %v", url, path, err)
		exit()
	}

	defer res.Body.Close()
	defer out.Close()

	if res.StatusCode != http.StatusOK {
		logger.Errorf("Received bad status while downloading %s: %s", url, res.Status)
		exit()
	}

	_, err = io.Copy(out, res.Body)

	if err == nil {
		logger.Infof("Successfully downloaded \"%s\" to \"%s\".", url, path)
	} else {
		logger.Errorf("Failed to write \"%s\" to \"%s\": %v.", url, path, err)
		exit()
	}
}
