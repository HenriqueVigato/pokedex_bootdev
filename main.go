package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// Uma struct no escopo global para que possa acessar o map[string] fora da main
type cmd struct {
	options map[string]cliCommand
}

func main() {
	opt := &cmd{}
	scanner := bufio.NewScanner(os.Stdin)

	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: func () error {
				return commandHelp(opt)
			},
		},
	}
	opt.options = commands

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("")
	fmt.Println("help to see the commands")

	for {
		fmt.Print("> ")

		success := scanner.Scan()
		if !success {
			fmt.Fprintln(os.Stderr, "Wasn't possible read a message")
			break
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "shouldn't see an error scannig a string")
		}

		message := scanner.Text()
		command := cleanInput(message)

		if len(command) == 0 {
			fmt.Println("Comando vazio")
			continue
		}

		if cmd, exist := commands[command[0]]; exist {

			err := cmd.callback()
			if err != nil {
				fmt.Println("Erro:", err)
			}
		} else {
			fmt.Println("Comando desconhecido. Digite 'help' para ver os comandos.")
		}

	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(command *cmd) error {
		fmt.Println("Comandos disponiveis: \n")
	for _, v := range command.options {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}
