package patcher

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

func getLastModified(url *string) (*time.Time, error) {
	res, err := http.Head(*url)
	if err != nil {
		return nil, err
	}

	lastModified := res.Header.Get("Last-Modified")
	if lastModified == "" {
		return nil, errors.New("couldn't get Last Modified header")
	}

	lastModifiedTime, err := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", lastModified)
	if err != nil {
		return nil, err
	}

	return &lastModifiedTime, nil
}

func compareLastModified(file *string, lastModified *time.Time) (bool, error) {
	stat, err := os.Stat(*file)

	if os.IsNotExist(err) {
		return true, nil
	}

	localLastModified := stat.ModTime()

	if localLastModified.Equal(*lastModified) {
		return false, nil
	}

	if localLastModified.Before(*lastModified) {
		return true, nil
	}

	return false, nil
}

func downloadFile(url string, file *string) error {
	lastModified, err := getLastModified(&url)
	if err != nil {
		return err
	}

	outdated, err := compareLastModified(file, lastModified)
	if err != nil {
		return err
	}

	if outdated {
		res, err := http.Get(url)

		if err != nil {
			return err
		}

		defer res.Body.Close()

		if _, err := os.Stat("files/"); os.IsNotExist(err) {
			err := os.Mkdir("files/", 0755)
			if err != nil {
				return err
			}
		}

		out, err := os.OpenFile(*file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, res.Body)
		if err != nil {
			return err
		}

		err = os.Chtimes(*file, *lastModified, *lastModified)
		if err != nil {
			return err
		}
	}

	return nil
}
