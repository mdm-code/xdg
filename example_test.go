package xdg_test

import (
	"fmt"

	"github.com/mdm-code/xdg"
)

func ExampleDataHomeDir() {
	dir := xdg.DataHomeDir()
	fmt.Println("Home data directory: ", dir)
}

func ExampleConfigHomeDir() {
	dir := xdg.ConfigHomeDir()
	fmt.Println("Config home directory: ", dir)
}

func ExampleStateHomeDir() {
	dir := xdg.StateHomeDir()
	fmt.Println("State home directory: ", dir)
}

func ExampleDataDirs() {
	dir := xdg.DataDirs()
	fmt.Println("Data directories: ", dir)
}

func ExampleConfigDirs() {
	dir := xdg.ConfigDirs()
	fmt.Println("Config directories: ", dir)
}

func ExampleCacheHomeDir() {
	dir := xdg.CacheHomeDir()
	fmt.Println("Cache home directory: ", dir)
}

func ExampleRuntimeDir() {
	dir := xdg.RuntimeDir()
	fmt.Println("Runtime home directory: ", dir)
}

func ExampleFind() {
	fpath := "program/file.data"
	if f, ok := xdg.Find(xdg.Data, fpath); ok {
		fmt.Printf("Data file for %s was found at %s\n", fpath, f)
	} else {
		fmt.Printf("ERROR: couldn't find %s\n", fpath)
	}
}
