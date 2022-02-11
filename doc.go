/*
Package xdg implements the XDG Base Directory Specification in Go. The
specification states where specification files should be searched. It does so
by defining base directories relative to which files are located. The details
of the specification are laid out here:
https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

This package supports most Unix-based operating systems. It should work fine
on MacOS.

Usage
	package main

	import (
		"fmt"

		"github.com/mdm-code/xdg"
	)

	func main() {
		// XDG base directory paths.
		dirs := []struct {
			msg string
			f   func() string
		}{
			{"Home data directory: ", xdg.DataHomeDir},
			{"Config home directory: ", xdg.CacheHomeDir},
			{"State home directory: ", xdg.StateHomeDir},
			{"Data directories: ", xdg.DataDirs},
			{"Config directories: ", xdg.ConfigDirs},
			{"Cache home directory: ", xdg.CacheHomeDir},
			{"Runtime home directory: ", xdg.RuntimeDir},
		}
		for _, d := range dirs {
			fmt.Println(d.msg, d.f())
		}

		// Search for file in data XDG directories.
		fpath := "program/file.data"
		if f, ok := xdg.Find(xdg.Data, fpath); ok {
			fmt.Println(f)
		} else {
			fmt.Printf("ERROR: couldn't find %s\n", fpath)
		}
	}
*/
package xdg
