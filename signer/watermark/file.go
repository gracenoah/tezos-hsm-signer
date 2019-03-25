package watermark

import (
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

// FileWatermark stores the last-signed level in local files
type FileWatermark struct {
	file    string
	session SessionWatermark
	mux     sync.Mutex
}

// GetFileWatermark returns a new file watermark manager
func GetFileWatermark(file string) *FileWatermark {
	// If file is not set, create a new file in our home directory
	if len(file) == 0 {
		file = path.Join(os.Getenv("HOME"), ".hsm-signer-watermarks")
	}
	// Load from disk
	watermarkEntries, err := loadFromDisk(file)
	if err != nil {
		log.Fatal("Unable to load watermark entries from: " + file)
	}

	wm := FileWatermark{
		file: file,
		session: SessionWatermark{
			watermarkEntries: watermarkEntries,
			mux:              sync.Mutex{},
		},
		mux: sync.Mutex{},
	}
	// Verify we can write to disk before returning
	if wm.saveToDisk() != nil {
		log.Fatal("Could not write to watermark file")
	}
	return &wm
}

func loadFromDisk(file string) ([]*watermarkEntry, error) {
	watermarkEntries := []*watermarkEntry{}

	// If file doesn't exist, return empty
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Println("Watermark file did not exist.  Initializing: ", file)
	} else {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println("Warning: Unable to read watermark file: ", file)
			return nil, err
		}
		err = yaml.Unmarshal(yamlFile, &watermarkEntries)
		if err != nil {
			log.Println("Warning: Unable to parse watermark file: ", file)
			return nil, err
		}
	}
	return watermarkEntries, nil
}

// save the watermark entries to disk
func (wm *FileWatermark) saveToDisk() error {
	bytes, err := yaml.Marshal(wm.session.watermarkEntries)
	if err != nil {
		log.Println("Unable to marshall watermark entries")
		return err
	}

	err = ioutil.WriteFile(wm.file, bytes, 0644)
	if err != nil {
		log.Println("Unable to write lockfile: " + wm.file)
		return err
	}
	return nil
}

// IsSafeToSign returns true if the provided (key, chainID, opType) tuple has
// not yet been signed at this or greater levels
func (wm *FileWatermark) IsSafeToSign(keyHash string, chainID string, opType uint8, level *big.Int) bool {
	wm.mux.Lock()
	defer wm.mux.Unlock()

	// Verify logic is safe
	isSessionSafe := wm.session.IsSafeToSign(keyHash, chainID, opType, level)

	// Update File
	err := wm.saveToDisk()
	if err != nil {
		return false
	}
	return isSessionSafe
}
