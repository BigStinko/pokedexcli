package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BigStinko/pokedexcli/pokeapi"
)

const PROMPT = "pokedex >> "

type config struct {
	pokeapiClient pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name string
	description string
	callback func(*config, ...string) error
}

func startREPL(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(PROMPT)
		
		if ok := scanner.Scan(); !ok {
			return
		}

		input := strings.Fields(scanner.Text())

		if len(input) == 0 {
			continue
		}

		commandName := input[0]
		args := []string{}

		if len(input) > 1 {
			args = input[1:]
		}

		if command, ok := getCommands()[commandName]; ok {
			err := command.callback(cfg, args...)
			if err != nil { fmt.Println(err) }
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}

	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the program",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "get the next page of locations",
			callback: commandMap,
		},
		"mapb": {
			name: "map",
			description: "get the previous page of locations",
			callback: commandMapb,
		},
		"explore": {
			name: "explore <location_name>",
			description: "Explore a location",
			callback: commandExplore,
		},
		"catch": {
			name: "catch <pokemon_name>",
			description: "Attempt to catch a pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect <pokemon_name>",
			description: "inspects a pokemon in your pokedex",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "lists the pokemon in your pokedex",
			callback: commandPokedex,
		},
	}
}
