package emotestore

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"

	"github.com/chpwssn/emote/emote"
)

const storeDir string = "localdata"

// Emotestore handles reading, writing, and finding Emotes in the datastore
type Emotestore struct {
	Rootpath string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (store Emotestore) localDataPath() string {
	return filepath.Join(store.Rootpath, storeDir)
}

func (store Emotestore) getEmotePath(name string) string {
	return filepath.Join(store.localDataPath(), name)
}

func (store Emotestore) getAllRecords() []emote.Emote {
	localDataPath := store.localDataPath()
	c, err := ioutil.ReadDir(localDataPath)
	check(err)
	var records []emote.Emote
	r, _ := regexp.Compile("(.*).json$")

	for _, entry := range c {
		match := r.MatchString(entry.Name())
		if match {
			name := r.FindStringSubmatch(entry.Name())
			record, _ := store.GetEmoteRecord(name[1])
			records = append(records, record)
		}
	}

	return records
}

func (store Emotestore) findRecord(name string) (emote.Emote, error) {
	var defaultResult emote.Emote
	for _, entry := range store.getAllRecords() {
		if entry.Name == name {
			return entry, nil
		}
	}
	return defaultResult, errors.New("No entry found with that name")
}

//Init initialize the datastore if necessary
func (store Emotestore) Init() {
	err := os.MkdirAll(store.localDataPath(), 0755)
	check(err)
}

// GetEmoteRecord fetch an emote record by name
func (store Emotestore) GetEmoteRecord(name string) (emote.Emote, error) {
	var emoteObj emote.Emote
	dat, err := ioutil.ReadFile(store.getEmotePath(name + ".json"))
	if err == nil {
		json.Unmarshal(dat, &emoteObj)
		return emoteObj, nil
	}
	return emoteObj, errors.New("Emote not found with that name")
}

// AllEmotes list the emotes in the local datastore
func (store Emotestore) AllEmotes() []emote.Emote {
	return store.getAllRecords()
}

// GetEmoteFileContents read the file contents of an emote's image
func (store Emotestore) GetEmoteFileContents(name string) ([]byte, error) {
	var result []byte
	emote, err := store.GetEmoteRecord(name)
	if err == nil {
		result, _ := ioutil.ReadFile(filepath.Join(store.localDataPath(), emote.Filename))
		return result, nil
	}
	return result, errors.New("Error reading file contents")
}

// StoreNewEmote write a new emote to the datastore
func (store Emotestore) StoreNewEmote(name string, credit string, file multipart.File, header multipart.FileHeader) (emote.Emote, error) {
	var result emote.Emote
	_, err := store.GetEmoteRecord(name)
	if err == nil {
		return result, errors.New("An emote with that name already exists")
	}
	outfile, err := os.Create(filepath.Join(store.localDataPath(), name))
	check(err)
	_, err = io.Copy(outfile, file)
	check(err)
	result = emote.Emote{
		Name:             name,
		Filename:         name,
		OriginalFilename: header.Filename,
		Credit:           credit,
	}
	f, err := os.Create(filepath.Join(store.localDataPath(), name+".json"))
	check(err)
	defer f.Close()
	json.NewEncoder(f).Encode(result)
	return result, nil
}
