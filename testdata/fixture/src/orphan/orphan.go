package orphan

// helper lives in its own directory with no anchor anywhere in it — an orphan
// code-area the engine must flag. (Its sibling src/ files are anchored, so only
// this directory is unanchored.)
func helper() int {
	return 42
}
