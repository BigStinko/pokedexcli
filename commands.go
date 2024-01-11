package main

import (
	"math/rand"
	"errors"
	"fmt"
	"os"
)

func commandHelp(cfg *config, args ...string) error {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args ...string) error {
	locationsresp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {	return err }

	cfg.nextLocationsURL = locationsresp.Next
	cfg.prevLocationsURL = locationsresp.Previous

	for _, location := range locationsresp.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil { return err }

	cfg.nextLocationsURL = locationResp.Next
	cfg.prevLocationsURL = locationResp.Previous

	for _, location := range locationResp.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("incorrect number of arguments")
	}

	name := args[0]
	location, err := cfg.pokeapiClient.GetLocation(name)
	if err != nil {	return err }

	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon: ")
	for _, encounter := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("incorrect number of arguments")
	}

	name := args[0]
	pokemon, err := cfg.pokeapiClient.GetPokemon(name)
	if err != nil { return err }

	res := rand.Intn(pokemon.BaseExperience)

	fmt.Printf("Trying to catch %s...\n", pokemon.Name)
	if res > 40 {
		fmt.Print("Failed\n")
		return nil
	}

	fmt.Printf("Caught %s\n", pokemon.Name)

	cfg.caughtPokemon[pokemon.Name] = pokemon
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("incorrect number of arguments")
	}

	name := args[0]
	pokemon, ok := cfg.caughtPokemon[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	
	fmt.Println("Stats:")

	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")

	for _, t := range pokemon.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for name := range cfg.caughtPokemon {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}
