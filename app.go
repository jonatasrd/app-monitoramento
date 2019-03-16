package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 10
const delay = 10

func main() {
	exibeMenu()

	switch leComando() {
	case 1:
		iniciarMonitoramento()
	case 2:
		imprimeLogs()
	case 3:
		fmt.Println("Até mais...")
		os.Exit(0)
	default:
		fmt.Println("Comando desconhecido!")
		os.Exit(-1)
	}

}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("3 - Sair do programa")
}

func leComando() int {
	var comando int
	fmt.Scan(&comando)
	return comando
}

func iniciarMonitoramento() {
	//sites := []string{"https://jonatas.dev", "https://github.com/jonatasrd", "https://www.alura.com"}
	fmt.Println("******************* INICIANDO O MONITORAMENTO *******************")
	sites := leSitesDoArquivo()
	for i := 0; i < monitoramentos; i++ {
		for _, site := range sites {
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
	}
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso.")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, "está com problemas na conexão. Status code:", resp.StatusCode)
		registraLog(site, false)
	}
}

func leSitesDoArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")
	//arquivo, err := ioutil.ReadFile("sites.txt")

	if err != nil {
		fmt.Println("ocorreu um erro:", err)
		os.Exit(-1)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}
	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	arquivo.WriteString(getHorarioAgora() + site + " - online: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func getHorarioAgora() string {
	return time.Now().Format("02/01/2006 15:04:05 ")
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	fmt.Println(string(arquivo))
}
