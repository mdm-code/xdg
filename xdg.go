package xdg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Data XDG base directory type.
	Data dir = iota
	// Config XDG base directory type.
	Config
	// State XDG base directory type.
	State
	// Cache XDG base directory type.
	Cache
	// Runtime XDG base directory type.
	Runtime

	dataHome   env = "XGD_DATA_HOME"
	configHome env = "XDG_CONFIG_HOME"
	stateHome  env = "XDG_STATE_HOME"
	dataDirs   env = "XDG_DATA_DIRS"
	configDirs env = "XDG_CONFIG_DIRS"
	cacheHome  env = "XDG_CACHE_HOME"
	runtimeDir env = "XDG_RUNTIME_DIR"
)

// XDG base directory type.
type dir uint8

type env string

type filterAdapter func(string) bool

type joinerAdapter func(path, string) path

// filter defines an interface to check whether some string meets
// specified condition(s).
type filter interface {
	filter(string) bool
}

// joiner defines an interface that handles joining two paths together.
type joiner interface {
	join(path, string) path
}

type path struct {
	value string
}

func (f filterAdapter) filter(s string) bool {
	return f(s)
}

func (j joinerAdapter) join(p path, s string) path {
	return j(p, s)
}

// require checks wether the object meets a certain condition.
func (p path) require(v filter) bool {
	return v.filter(p.value)
}

// join joins the value of the path struct with the function attribute.
// It returns a copy of the path struct.
func (p path) join(s string) path {
	return path{filepath.Join(p.value, s)}
}

// split splits the value of the path struct into paths separated by the path
// list separator.
func (p path) split() []path {
	result := []path{}
	sep := string(os.PathListSeparator)
	for _, s := range strings.Split(p.value, sep) {
		result = append(result, path{s})
	}
	return result
}

// A single base directory relative to which user-specific files should
// be written. Default: $HOME/.local/share.
func DataHomeDir() string {
	return valueOf(dataHome, joinHome(".local/share"))
}

// A single base directory relative to which user-specific configuration
// files should be written. Default: $HOME/.config.
func ConfigHomeDir() string {
	return valueOf(configHome, joinHome(".config"))
}

// A single base directory relative to which user-specific state data
// should be written. Default: $HOME/.local/state.
func StateHomeDir() string {
	return valueOf(stateHome, joinHome(".local/state"))
}

// A set of preference-ordered base directories relative to which data
// files should be searched. Default: /usr/local/share/:/usr/share/.
func DataDirs() string {
	return valueOf(dataDirs, "/usr/local/share/:/usr/share/")
}

// A set of preference-ordered base directories relative to which
// configuration files should be searched. Default: /etc/xdg.
func ConfigDirs() string {
	return valueOf(configDirs, "/etc/xdg")
}

// A single base directory relative to which user-specific, non-essential
// (cached) data should be written. Default: $HOME/.cache.
func CacheHomeDir() string {
	return valueOf(cacheHome, joinHome(".cache"))
}

// A single base directory relative to which user-specific, runtime files
// and other file objects should be placed. It defaults to $TMPDIR on Unix
// if non-empty else /tmp.
func RuntimeDir() string {
	return valueOf(runtimeDir, os.TempDir())
}

// isAbsolute verifies if a path is absolute, not a relative one. XDG Base
// Directory Standard states that all XDG environmental variables must be
// absolute. A relative path is specified to be considered invalid and ignored.
func isAbsolute(s string) bool {
	return filepath.IsAbs(s)
}

// isExist checks whether a path exists on the file system.
func isExist(s string) bool {
	_, err := os.Stat(s)
	return !errors.Is(err, os.ErrNotExist)
}

// isList verifies if a path is a comma-separated list of paths.
func isList(s string) bool {
	sep := string(os.PathListSeparator)
	if strings.Contains(s, sep) {
		return true
	}
	return false
}

// valueOf returns an XDG environmental variable for the user. If it is unset,
// empty or it is not an absolute path, it will return the default XDG base
// directory or preference-ordered base directories.
func valueOf(key env, fallback string) string {
	val := os.Getenv(string(key))
	p := path{val}
	// NOTE: Empty strings are not absolute.
	if !p.require(filterAdapter(isAbsolute)) {
		return fallback
	}
	return val
}

// joinHome joins the value of the $HOME environmental variable with the path
// function parameter.
func joinHome(s string) string {
	home, _ := os.UserHomeDir()
	result := filepath.Join(home, s)
	return result
}

// applyFilter filters out paths based on the provided filter criteria.
func applyFilter(paths []path, filters ...filter) []path {
	if len(filters) == 0 {
		return paths
	}

	result := []path{}

	for _, p := range paths {
		keep := true

		for _, f := range filters {
			if !f.filter(p.value) {
				keep = false
				break
			}
		}
		if keep {
			result = append(result, p)
		}
	}
	return result
}

// mapJoin joins all paths with the provided string.
func mapJoin(paths []path, s string, j joiner) []path {
	result := []path{}
	for _, p := range paths {
		result = append(result, j.join(p, s))
	}
	return result
}

// Finds searches for a given path in specified XDG directories. It returns an empty
// string and false if the path is not found.
func Find(d dir, path string) (result string, ok bool) {
	switch d {
	case Data:
		result, ok = find(path, DataHomeDir(), DataDirs())
	case Config:
		result, ok = find(path, ConfigHomeDir(), ConfigDirs())
	case State:
		result, ok = find(path, StateHomeDir())
	case Cache:
		result, ok = find(path, CacheHomeDir())
	case Runtime:
		result, ok = find(path, RuntimeDir())
	}
	return result, ok
}

// find looks for a given path in XDG directories. It returns an empty string and
// false if none was found.
func find(pth string, dirs ...string) (string, bool) {
	paths := make([]path, 0, 4)

	for _, d := range dirs {
		p := path{d}
		if p.require(filterAdapter(isList)) {
			j := joinerAdapter(path.join)
			paths = append(paths, mapJoin(p.split(), pth, j)...)
		} else {
			paths = append(paths, p.join(pth))
		}
	}

	paths = applyFilter(
		paths,
		filterAdapter(isAbsolute),
		filterAdapter(isExist),
	)

	if len(paths) == 0 {
		return "", false
	}

	return paths[0].value, true
}
