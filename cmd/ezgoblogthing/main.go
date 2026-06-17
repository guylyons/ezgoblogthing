package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"ezgoblogthing/internal/site"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	command := "build"
	if len(args) > 0 {
		command = args[0]
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}
	out := filepath.Join(root, "dist")

	switch command {
	case "build":
		return build(root, out)
	case "serve":
		if err := build(root, out); err != nil {
			return err
		}
		addr := ":8080"
		fmt.Printf("Serving %s at http://localhost%s\n", out, addr)
		return http.ListenAndServe(addr, http.FileServer(http.Dir(out)))
	default:
		return errors.New("usage: go run ./cmd/ezgoblogthing [build|serve]")
	}
}

func build(root string, out string) error {
	if err := site.Build(root, out); err != nil {
		return err
	}
	fmt.Printf("Generated %s\n", out)
	return nil
}
