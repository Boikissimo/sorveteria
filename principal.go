package main

import (
	_ "bufio"
	_ "bytes"
	"context"
	_ "crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
)

var configs configuracoes

type configuracoes struct {
	Ip      string `json:"Ip"`
	Porta   string `json:"Porta"`
	Usuario string `json:"Usuário"`
	//Senha	string `json:"Senha"`
}

//var senha string
//var salvo bool

func init() {
	teste, er := os.Stat("configuração.config")
	if os.IsNotExist(er) {
		configuracao()
	} else {
		if teste.Size() < 4 {
			configuracao()
		} else {
			carregar()
		}
	}
}

func main() {
	senha := ""
	fmt.Println("Por favor, insira a senha:")
	fmt.Scanln(&senha)
	data, err := pgx.Connect(context.Background(), "postgres://"+configs.Usuario+":"+senha+"@"+configs.Ip+":"+configs.Porta+"/sorveteria")
	if err != nil {
		fmt.Fprintf(os.Stderr, "É miga, não deu certo:\n %v\n", err)
		os.Exit(1)
	}
	tabela, err := data.Query(context.Background(), "SELECT tablename from pg_tables WHERE schemaname = 'producao'")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Falhou em pegar a linha, mas não se preoculpe, ligue somente durante o programa\n %v\n", err)
		os.Exit(1)
	}

	for tabela.Next() {
		teste := ""
		tabela.Scan(&teste)
		fmt.Println(teste)
	}

	data.Close(context.Background())
	fmt.Scanln()
}

func configuracao() {
	configs.Usuario = usuario()
	configs.Ip = ip()
	configs.Porta = porta()
	bytes, err := json.MarshalIndent(configs, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("configuração.config", bytes, 0644)
}

func carregar() {
	conteudo, err := ioutil.ReadFile("configuração.config")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(conteudo, &configs)
	if err != nil {
		fmt.Println(err)
	}
}

func usuario() string {
	resposta := ""
	fmt.Println("Você deseja usar o usuário padrão (postgres)? (s ou n)")
	fmt.Scanln(&resposta)
	if strings.ToLower(resposta) == "s" {
		return "postgres"
	}
	for {
		temp := ""
		fmt.Println("Por favor, insira o nome de usúario que você deseja usar:")
		fmt.Scanln(&resposta)
		fmt.Println("O nome de usuário está correto?\nNome de usuário: " + resposta + "\nDigite n para mudar, aperte enter para continuar")
		fmt.Scanln(&temp)
		temp = strings.ToLower(temp)
		if temp == "n" {
			continue
		}
		break
	}
	return resposta
}
func ip() string {
	resposta := ""
	fmt.Println("O banco de dados está em outra maquina? (s ou n)")
	fmt.Scanln(&resposta)
	resposta = strings.ToLower(resposta)
	if resposta == "s" {
		fmt.Println("Digite o IP ou endereço do banco de dados, excluindo a porta: ")
		for {
			fmt.Scanln(&resposta)
			if strings.Count(resposta, ".") != 3 && strings.Count(resposta, ":") != 7 {
				fmt.Println("Você digitou um ip invalido, por favor digite o ip novamente")
				continue
			}
			break
		}
		if strings.Count(resposta, ":") == 7 {
			resposta = "[" + resposta + "]"
		}
		return resposta
	} else {
		return "localhost"
	}
}
func porta() string {
	resposta := ""
	fmt.Println("A porta é a padrão \"5432\" ou é outra? (s ou n)")
	fmt.Scanln(&resposta)
	resposta = strings.ToLower(resposta)
	if resposta == "n" {
		fmt.Println("Digite a porta:")
		fmt.Scanln(&resposta)
	} else {
		resposta = "5432"
	}
	return resposta
}

/*func senha() string {
	//Isto deixa a senha salva em um arquivo json, deixando criptografado
}*/
