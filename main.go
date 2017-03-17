package main

import (
	"path/filepath"
	"os"
	"fmt"
	"crypto/rand"
	"github.com/secsy/goftp"
	"time"
)

// Configuracao do client ftp
var configura = goftp.Config{
	User:               "client", // User Ftp
	Password:           "pwd_ftp", // Password Ftp
	ConnectionsPerHost: 10,
	Timeout:            10 * time.Second,
	Logger:             os.Stderr,
}
var client, err = goftp.DialConfig(configura, "ip_host") // Informe IP ou dominio do ftp
var maxid = 3 // Aqui voce pode definir quantas go rotines vao executar de uma unica vez
var id = 0

func main() {
	fmt.Println("Por favor aguarde o analise da pasta...")
	checkErr(err)
	defer client.Close()
	searcDir := "./" // Pasta onde ele vai comecar a analisar os arquivos
	filepath.Walk(searcDir, func(path string, f os.FileInfo, err error) error { // A funcao filepath.Walk faz um analise na pasta e subpastas te retornando o caminho completo
		if f.IsDir() == false {
			id++
			go sendFtpFile(path)
			wait()
		}
		return nil
	})
}
// Funcao de analise da extensao e envio do arquivo
func sendFtpFile(file string) {
	ext := filepath.Ext(file)
	if ext == ".xml" {
		arq := make([]byte, 12)
		rand.Read(arq)
		arq2 := fmt.Sprintf("%X", arq)
		b, err := os.Open(file)
		checkErr(err)
		client.Store(arq2+"."+ext, b)
		fmt.Println(file+" - ", arq2+"."+ext)
	}
	id--
}

// Funcao para fazer checagem de erro
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func wait() {
	for {
		if id >= maxid {
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}
}
