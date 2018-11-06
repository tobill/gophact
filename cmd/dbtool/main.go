package main

import (
	"path/filepath"
	"os"
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
)

func printInfo(opts badger.Options) error {
	dbClient, err :=  badger.Open(opts)
	if err != nil 	{
		fmt.Printf("%v", err)
		return	err	
	}
	defer dbClient.Close()
	lsm, vlog := dbClient.Size()
	fmt.Printf("Size of LSM: %d, size of vlog: %d", lsm, vlog)
	return err
}

func removeLockfileAndTruncate(opts badger.Options) error {
	lockFile := "LOCK"
	absLockPath := filepath.Join(opts.Dir, lockFile)
	err := os.Remove(absLockPath)
	opts.Truncate = true
	dbClient, err :=  badger.Open(opts)
	defer dbClient.Close()
	if err != nil 	{
		fmt.Printf("%v", err)
		return	err	
	}
	return err
}

func main() {
	action :=  flag.String("action", "info", "action to do")
	dbPath := flag.String("dbPath", "", "db path")
	
	flag.Parse()

	opts := badger.DefaultOptions
	opts.Dir = *dbPath
	opts.ValueDir = *dbPath

	switch *action {
	case "removeLock":
		err := removeLockfileAndTruncate(opts)
		if (err != nil) {
			fmt.Printf("%v", err)
		}
	case "info":
		err := printInfo(opts)
		if (err != nil) {
			fmt.Printf("%v", err)
		}
	default:
		fmt.Printf("Nothing to do")
	} 

}