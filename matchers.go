package href

import "strings"

// path begins with "/" or is empty:
// 	path-abempty  = *( "/" segment )
func matchPathAbEmpty(path string) bool {
	return path == "" || path[0] == '/'
}

// path begins with "/" but not "//":
// 	path-absolute = "/" [ segment-nz *( "/" segment ) ]
func matchPathAbsolute(path string) bool {
	return path != "" && path[0] == '/' && (len(path) == 1 || path[1] != '/')
}

// path begins with a segment:
// 	path-rootless = segment-nz *( "/" segment )
func matchPathRootless(path string) bool {
	return path != "" && path[0] != '/'
}

// path with zero characters:
// 	path-empty = 0<pchar>
func matchPathEmpty(path string) bool {
	return path == ""
}

// begins with a non-colon segment
//	path-noscheme = segment-nz-nc *( "/" segment )
func matchPathNoScheme(path string) bool {
	segments := strings.SplitN(path, "/", 2)

	return segments[0] != "" && !strings.Contains(segments[0], ":")
}
