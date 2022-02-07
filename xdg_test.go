package xdg

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Verify if paths are correctly classified as absolute or relative.
func TestIsAbsolute(t *testing.T) {
	data := []struct {
		name     string
		path     string
		expected bool
	}{
		{"absolute", "/usr/bin", true},
		{"root", "/", true},
		{"relative", "../bin", false},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := isAbsolute(d.path)
			if result != d.expected {
				t.Errorf("Expected %t; got %t", d.expected, result)
			}
		})
	}
}

// Check if Path.require() allows the verify whether Path meets some condition.
func TestPathStructRquire(t *testing.T) {
	data := []struct {
		name     string
		path     path
		checker  func(string) bool
		expected bool
	}{
		{"absolute", path{"/etc"}, isAbsolute, true},
		{"root", path{"/"}, isAbsolute, true},
		{"relative", path{"../"}, isAbsolute, false},
		{"root exists", path{"/"}, isExist, true},
		{"is a list", path{"/bin:/usr/bin"}, isList, true},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			v := filterAdapter(d.checker)
			if result := d.path.require(v); result != d.expected {
				t.Errorf("Expected: %t; got %t", d.expected, result)
			}
		})
	}
}

// Test if the following XDG environmental variables are set.
func TestGettingXdgEnv(t *testing.T) {
	tmpSetEnv := func(key, val string) func() {
		tmp := os.Getenv(key)
		os.Setenv(key, val)
		return func() {
			os.Setenv(key, tmp)
		}
	}
	data := []struct {
		name     string
		key      env
		expected string
		set      func(string, string) func()
	}{
		{string(dataHome), dataHome, "/some/path", tmpSetEnv},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			reset := d.set(string(d.key), d.expected)
			defer reset()
			if result := valueOf(d.key, d.expected); result != d.expected {
				t.Errorf("Expected %s; got %s", d.expected, result)
			}
		})
	}
}

// Verify if joinHome prepends $HOME to the path in a controlled fashion.
func TestJoinHome(t *testing.T) {
	tmp := os.Getenv("HOME")
	os.Setenv("HOME", "/home/user")
	defer os.Setenv("HOME", tmp)

	data := []struct {
		name string
		path string
	}{
		{".local/share", ".local/share"},
		{".local/state", ".local/state"},
		{".config", ".config"},
		{".cache", ".cache"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			home, err := os.UserHomeDir()
			if err != nil {
				t.Fatalf("Failed to execute test %s", err)
			}
			expected := filepath.Join(home, d.path)
			if result := joinHome(d.path); result != expected {
				t.Errorf("Expected: %s; got: %s", expected, result)
			}
		})
	}
}

// Check if the isExist check returns the expected boolean value.
func TestIsExist(t *testing.T) {
	data := []struct {
		path     string
		expected bool
	}{
		{"/", true},
		{"/usr/bin", true},
		{"/etc/", true},
		{"/user", false},
		{"/home/anonymous", false},
	}
	for _, d := range data {
		t.Run(d.path, func(t *testing.T) {
			if result := isExist(d.path); result != d.expected {
				t.Errorf("Expected: %t; got: %t", d.expected, result)
			}
		})
	}
}

// Test if a given string can be seen as a list of paths.
func TestIsList(t *testing.T) {
	data := []struct {
		path     string
		expected bool
	}{
		{"/home", false},
		{"/usr/bin:/bin", true},
		{"", false},
	}
	for _, d := range data {
		t.Run(d.path, func(t *testing.T) {
			if result := isList(d.path); result != d.expected {
				t.Errorf("Expected: %t; got: %t", d.expected, result)
			}
		})
	}
}

// Test splitting up the path on the Unix path list separator.
func TestPathSplits(t *testing.T) {
	data := []struct {
		path     string
		expected []path
	}{
		{"/usr/local/share/:/usr/share/",
			[]path{
				{"/usr/local/share/"},
				{"/usr/share/"},
			},
		},
		{"/bin:/usr/bin:/usr/local/bin",
			[]path{
				{"/bin"},
				{"/usr/bin"},
				{"/usr/local/bin"},
			},
		},
		{"/usr/local/share/,/usr/share/",
			[]path{
				{"/usr/local/share/,/usr/share/"},
			},
		},
		{"",
			[]path{
				{""},
			},
		},
	}
	for _, d := range data {
		t.Run(d.path, func(t *testing.T) {
			p := path{d.path}
			if result := p.split(); !reflect.DeepEqual(result, d.expected) {
				t.Errorf("Expected: %v; got: %v", d.expected, result)
			}
		})
	}
}

// Test filtering paths with a filter function.
func TestFilterPath(t *testing.T) {
	paths := []path{
		{"/"},
		{"/etc"},
		{"/usr/bin"},
		{"Documents"},
	}
	expected := []path{
		{"/"},
		{"/etc"},
		{"/usr/bin"},
	}
	result := applyFilter(paths, filterAdapter(isAbsolute))
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v; got: %v", expected, result)
	}
}

// Test if the same paths are returned without any filter being applied to them.
func TestApplyNoFilter(t *testing.T) {
	paths := []path{
		{"/"},
		{"/etc"},
		{"/usr/bin"},
		{".."},
		{"../Documents"},
	}
	result := applyFilter(paths)
	if !reflect.DeepEqual(result, paths) {
		t.Errorf("Expected: %v; got: %v", paths, result)
	}
}

// Test the joined value of the returned path struct.
func TestJoinPath(t *testing.T) {
	p := path{"/usr/local/share"}
	result := p.join("prog/file")
	expected := path{"/usr/local/share/prog/file"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v; got: %v", expected, result)
	}
}

// Verify detaching struct method and mapping it over a list of path structs.
func TestMapJoin(t *testing.T) {
	data := []struct {
		name     string
		paths    []path
		s        string
		expected []path
	}{
		{
			"/usr/share/prog/file",
			[]path{{"/usr/share"}, {"/usr/local/share"}},
			"prog/file",
			[]path{{"/usr/share/prog/file"}, {"/usr/local/share/prog/file"}},
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := mapJoin(d.paths, d.s, joinerAdapter(path.join))
			if !reflect.DeepEqual(result, d.expected) {
				t.Errorf("Expected: %v; got: %v", d.expected, result)
			}
		})
	}
}

// Check whether Find discovers the first valid path from the preference ordered
// set.
func TestFind(t *testing.T) {
	// NOTE: I am testing for empty default dirs. I do not want to create any dirs.
	data := []struct {
		name string
		dir  dir
	}{
		{"data", Data},
		{"config", Config},
		{"state", State},
		{"cache", Cache},
		{"runtime", Runtime},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			if _, ok := Find(d.dir, ""); !ok {
				t.Errorf("Find failed for %s", d.name)
			}
		})
	}
}

// Verify if unexported find() fails when len(paths) == 0.
func TestFindFails(t *testing.T) {
	dirs := []struct {
		name     string
		dir      string
		expected bool
	}{
		{"non-existent", "/non/existent/dir", false},
		{"non-absolute", "../non_absolute", false},
	}
	for _, d := range dirs {
		t.Run(d.name, func(t *testing.T) {
			if _, ok := find("missing/file", d.dir); ok != d.expected {
				t.Errorf("Expected: %t; got: %t", d.expected, ok)
			}
		})
	}
}

// Test if variables for XDG base directories are dynamically changed when the
// value is changed.
func TestDynamic(t *testing.T) {
	tmpSetEnv := func(key, val string) func() {
		tmp := os.Getenv(key)
		os.Setenv(key, val)
		return func() {
			os.Setenv(key, tmp)
		}
	}
	data := []struct {
		name     string
		key      env
		expected string
		callable func() string
		set      func(string, string) func()
	}{
		{string(dataHome), dataHome, "/some/data/home", DataHomeDir, tmpSetEnv},
		{string(configHome), configHome, "/some/config/path", ConfigHomeDir, tmpSetEnv},
		{string(stateHome), stateHome, "/some/state/home", StateHomeDir, tmpSetEnv},
		{string(dataDirs), dataDirs, "/some/data/:/data/dirs/", DataDirs, tmpSetEnv},
		{string(configDirs), configDirs, "/some/config/:/config/dirs/", ConfigDirs, tmpSetEnv},
		{string(cacheHome), cacheHome, "/random/cache/home", CacheHomeDir, tmpSetEnv},
		{string(runtimeDir), runtimeDir, "/local/temp/dir", RuntimeDir, tmpSetEnv},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			reset := d.set(string(d.key), d.expected)
			defer reset()
			if result := d.callable(); result != d.expected {
				t.Errorf("Expected %s; got %s", d.expected, result)
			}
		})
	}
}
