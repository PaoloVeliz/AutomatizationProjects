package main
import (
	"fmt"
    "io/ioutil"
    "path/filepath"
	"log"
	"time"
    "os"
    "io"
)

func addTimestamp(filename string) string {
	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHmmss
	fileExtension := filepath.Ext(filename)
	return fmt.Sprintf("%s%s%s", filename[:len(filename)-len(fileExtension)], timestamp, fileExtension)
}

func isDirEmpty(dirPath string) (bool, error) {
    dir, err := os.Open(dirPath)
    if err != nil {
        return false, err
    }
    defer dir.Close()
    _, err = dir.Readdirnames(1)
    if err == nil {
        return false, nil
    }
    if err == io.EOF {
        return true, nil
    }
    return false, err
}

func backup(folder_to_get_files string,folder_to_send_files string) {
    files, err := ioutil.ReadDir(folder_to_get_files)
	if err != nil {
		log.Fatal(err)
	}
    dir_empty, err := isDirEmpty(folder_to_get_files) 
    if dir_empty {
        fmt.Println("The folder is empty.")
    }else {
        for _, file := range files {
            fmt.Println("Moving:", file.Name())
            err := os.Rename(folder_to_get_files + "/" + file.Name(), folder_to_send_files + "/" + file.Name())
            renamed_file := os.Rename(folder_to_send_files + "/" + file.Name(), folder_to_send_files + "/" + addTimestamp(file.Name()))
            if err != nil || renamed_file != nil {
                log.Fatal(err)
            }
            newLogLine(file.Name())
        }
        fmt.Println("BackUp Done")
    }
	
}

func newLogLine(message string){
    logFile, err := os.OpenFile("mylogfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
    defer logFile.Close()
    logger := log.New(logFile, "auto-backup ", log.LstdFlags)
    logger.Printf("Files and folder moved succesfully: %d\n", message)
}

func main() {
    ticker := time.NewTicker(time.Minute)
    backup("where are the files","where you gonna send them")
	for {
        select {
        case <-ticker.C:
            backup("where are the files","where you gonna send them")
		}
	}
}
