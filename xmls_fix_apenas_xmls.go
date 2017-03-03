package main

import (
	"path/filepath"
	"os"
	"fmt"
	"crypto/rand"
	"github.com/secsy/goftp"
	"time"
)

func main() {
	fmt.Println("Por favor aguarde o analise da pasta...")
	configura := goftp.Config{
		User: "user_ftp",
		Password: "pwd_ftp",
		ConnectionsPerHost: 10,
		Timeout: 10 * time.Second,
		Logger: os.Stderr,
	}

	client, err :=goftp.DialConfig(configura,"ip_host")
	checkErr(err)

	defer client.Close()

	searcDir := "./"
	fileList := []string{}
	filepath.Walk(searcDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() == false {
			fileList = append(fileList, path)
		}
		return nil
	})

	for _, file := range fileList {
		ext := string(file[len(file)-3:])
		if ext =="xml" || ext=="zip"{
			arq := make([]byte,12)
			rand.Read(arq)
			arq2 := fmt.Sprintf("%X",arq)
			b, err := os.Open(file)
			checkErr(err)
			client.Store(arq2+"."+ext,b)
			fmt.Println(file + " - ",arq2+"."+ext)
		}

	}
}

func checkErr (err error){
	if err != nil {
		panic(err.Error())
	}
}
