# XDG package

![(logo)](./logo.png)

<h4 align="center">The XDG Base Directory Specification implemented in Go.</h4>


<!-- TODO (michal): Add badges -->


The XDG Base Directory Specification (LINK) allows you specify directories where
runtime files, configurations, data and caches are kept. The file discovery
process is automatic and adheres to the XDG standard.

This package supports most Unix-based operating systems. It should work fine on
MacOS. I wrote this package for my personal needs: to deduplicate this kind of
functionality from my other programs, but it is very much a self-contained
implementation.

See [Usage](#usage) section below for examples. Package documentation is
available here: https://pkg.go.dev/github.com/mdm-code/xdg.


## Installation

```sh
go get github.com/mdm-code/xdg
```


## Default locations

The table shows default values for XDG environmental variables for Unix-like systems:

| <a href="#default-locations"><img width="900" height="0"></a> | <a href="#default-locations"><img width="1100" height="0"></a><p>Unix-like</p>  |
| :------------------------------------------------------------: | :----------------------------------------------------------------------------: |
| <kbd><b>XDG_DATA_HOME</b></kbd>                                | <kbd>$HOME/.local/share</kbd>                                                  |
| <kbd><b>XDG_CONFIG_HOME</b></kbd>                              | <kbd>$HOME/.config</kbd>                                                       |
| <kbd><b>XDG_STATE_HOME</b></kbd>                               | <kbd>$HOME/.local/state</kbd>                                                  |
| <kbd><b>XDG_DATA_DIRS</b></kbd>                                | <kbd>/usr/local/share/:/usr/share/</kbd><br/><kbd>/usr/share</kbd>             |
| <kbd><b>XDG_CONFIG_DIRS</b></kbd>                              | <kbd>/etc/xdg</kbd>                                                            |
| <kbd><b>XDG_CACHE_HOME</b></kbd>                               | <kbd>$HOME/.cache</kbd>                                                        |
| <kbd><b>XDG_RUNTIME_DIR</b></kbd>                              | <kbd>$TMPDIR</kbd>                                                             |


## Usage

Here is an example of how to use the public API of the `xdg` package.

```go
package main

import (
	"fmt"

	"github.com/mdm-code/xdg"
)

func main() {
	// XDG base directory paths.
	dirs := []struct {
		msg string
		pth string
	}{
		{"Home data directory: ", xdg.DataHomeDir()},
		{"Config home directory: ", xdg.CacheHomeDir()},
		{"State home directory: ", xdg.StateHomeDir()},
		{"Data directories: ", xdg.DataDirs()},
		{"Config directories: ", xdg.ConfigDirs()},
		{"Cache home directory: ", xdg.CacheHomeDir()},
		{"Runtime home directory: ", xdg.RuntimeDir()},
	}
	for _, d := range dirs {
		fmt.Println(d.msg, d.pth)
	}

	// Search for file in data XDG directories.
	fpath := "/prog/file"
	if f, ok := xdg.Find(xdg.Data, fpath); ok {
		fmt.Println(f)
	} else {
		fmt.Printf("ERROR: couldn't find %s\n", fpath)
	}
}
```


## Development

Consult `Makefile` to see how to format, examine code with `go vet`, run unit test,
get test coverage and check if the package builds all right.


## License

Copyright (c) 2022 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
