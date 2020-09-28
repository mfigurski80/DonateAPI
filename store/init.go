package store

import "os"

var dataPath, logPath string

// Init initializes data store
func Init(dataFolder string) error {
	// create dir
	err := os.MkdirAll(dataFolder, 0755)
	if err != nil {
		return err
	}
	// create data files
	dataPath = dataFolder + "/users.json"
	logPath = dataFolder + "/main.log"
	file, err := os.OpenFile(dataPath, os.O_RDONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	file, err = os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	// return
	return nil
}
