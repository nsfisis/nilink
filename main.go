package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: nilink <serve|add|remove|list>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "serve":
		cmdServe(os.Args[2:])
	case "add":
		cmdAdd(os.Args[2:])
	case "remove":
		cmdRemove(os.Args[2:])
	case "list":
		cmdList(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func cmdAdd(args []string) {
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	dbPath := flags.String("db", "data/nilink.db", "database path")
	flags.Parse(args)

	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: nilink add [-db path] <url>")
		os.Exit(1)
	}
	url := flags.Arg(0)

	db, err := openDB(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	id, err := insertLink(db, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	short, err := encodeID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s %s\n", short, url)
}

func cmdRemove(args []string) {
	flags := flag.NewFlagSet("remove", flag.ExitOnError)
	dbPath := flags.String("db", "data/nilink.db", "database path")
	flags.Parse(args)

	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: nilink remove [-db path] <short-id>")
		os.Exit(1)
	}
	shortID := flags.Arg(0)

	id, err := decodeID(shortID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	db, err := openDB(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := deleteLink(db, id); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func cmdList(args []string) {
	flags := flag.NewFlagSet("list", flag.ExitOnError)
	dbPath := flags.String("db", "data/nilink.db", "database path")
	flags.Parse(args)

	db, err := openDB(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	links, err := listLinks(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, l := range links {
		short, err := encodeID(l.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s  %s  %s\n", short, l.URL, l.CreatedAt)
	}
}
