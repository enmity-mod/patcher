package patcher

import (
	"errors"
	"log"
	"os"

	"howett.net/plist"
)

const (
	DEFAULT_IPA_PATH    = "files/Discord.ipa"
	DEFAULT_ICONS_PATH  = "files/icons.zip"
	DEFAULT_ENMITY_PATH = "files/Enmity.dylib"
)

const (
	IPA_URL    = "https://files.enmity.app/Discord.ipa"
	ICONS_URL  = "https://files.enmity.app/icons.zip"
	ENMITY_URL = "https://files.enmity.app/Enmity.dylib"
)

func PatchDiscord(discordPath *string, iconsPath *string, dylibPath *string) {
	log.Println("starting patcher")

	checkFile(discordPath, DEFAULT_IPA_PATH, IPA_URL)
	checkFile(iconsPath, DEFAULT_ICONS_PATH, ICONS_URL)
	checkFile(dylibPath, DEFAULT_ENMITY_PATH, ENMITY_URL)

	extractDiscord(discordPath)

	log.Println("renaming Discord to Enmity")
	if err := patchName(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Discord renamed")

	log.Println("adding Enmity url scheme")
	if err := patchSchemes(); err != nil {
		log.Fatalln(err)
	}
	log.Println("url scheme added")

	log.Println("remove devices whitelist")
	if err := patchDevices(); err != nil {
		log.Fatalln(err)
	}
	log.Println("device whitelist removed")

	log.Println("patch Discord icons")
	extractIcons(iconsPath)
	if err := patchIcon(); err != nil {
		log.Fatalln(err)
	}
	log.Println("icons patched")

	packDiscord()
	log.Println("cleaning up")
	clearPayload()

	log.Println("done!")
}

// Check if file exists, download if not found
func checkFile(path *string, defaultPath string, url string) {
	_, err := os.Stat(*path)
	if errors.Is(err, os.ErrNotExist) {
		if *path == defaultPath {
			log.Println("downloading", url, "to", *path)
			err := downloadFile(url, path)
			if err != nil {
				log.Println("error downloading", url)
				log.Fatalln(err)
			}
		} else {
			log.Fatalln("file not found", *path)
		}
	}
}

// Delete the payload folder
func clearPayload() {
	err := os.RemoveAll("Payload")
	if err != nil {
		log.Panicln(err)
	}
}

// Load Discord's plist file
func loadPlist() (map[string]interface{}, error) {
	infoFile, err := os.Open("Payload/Discord.app/Info.plist")
	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	decoder := plist.NewDecoder(infoFile)
	if err := decoder.Decode(&info); err != nil {
		return nil, err
	}

	return info, nil
}

// Save Discord's plist file
func savePlist(info *map[string]interface{}) error {
	infoFile, err := os.OpenFile("Payload/Discord.app/Info.plist", os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	encoder := plist.NewEncoder(infoFile)
	err = encoder.Encode(*info)
	return err
}

// Patch Discord's name
func patchName() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	info["CFBundleName"] = "Enmity"
	info["CFBundleDisplayName"] = "Enmity"

	err = savePlist(&info)
	return err
}

// Patch Discord's URL scheme to add Enmity's URL handler
func patchSchemes() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	info["CFBundleURLTypes"] = append(
		info["CFBundleURLTypes"].([]interface{}),
		map[string]interface{}{
			"CFBundleURLName": "Enmity",
			"CFBundleURLSchemes": []string{
				"enmity",
			},
		},
	)

	err = savePlist(&info)
	return err
}

// Remove Discord's device limits
func patchDevices() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	delete(info, "UISupportedDevices")

	err = savePlist(&info)
	return err
}

// Patch the Discord icon to use Enmity's icon
func patchIcon() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	info["CFBundleIcons"].(map[string]interface{})["CFBundlePrimaryIcon"].(map[string]interface{})["CFBundleIconName"] = "EnmityIcon"
	info["CFBundleIcons"].(map[string]interface{})["CFBundlePrimaryIcon"].(map[string]interface{})["CFBundleIconFiles"] = []string{"EnmityIcon60x60"}

	info["CFBundleIcons~ipad"].(map[string]interface{})["CFBundlePrimaryIcon"].(map[string]interface{})["CFBundleIconName"] = "EnmityIcon"
	info["CFBundleIcons~ipad"].(map[string]interface{})["CFBundlePrimaryIcon"].(map[string]interface{})["CFBundleIconFiles"] = []string{"EnmityIcon60x60", "EnmityIcon76x76"}

	err = savePlist(&info)
	return err
}
