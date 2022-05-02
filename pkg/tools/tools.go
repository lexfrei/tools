package tools

import "regexp"

//nolint:varcheck,deadcode // for the future use
var (
	restrictedCharsInWindows = []rune("<>:\"/\\|?*")
	restrictedCharsInLinux   = []rune("/")

	restrictedBytesInWindows = []byte{
		// Control ASCII chars
		0o0, 0o1, 0o2, 0o3, 0o4, 0o5, 0o6, 0o7,
		10, 11, 12, 13, 14, 15, 16, 17,
		20, 21, 22, 23, 24, 25, 26, 27,
		31, 32, 33, 34, 35, 36, 37,
		127,
	}
	restrictedBytesInLinux = []byte{
		// Null byte
		0o0,
	}

	restrictedDirNamesInLinux   = []string{}
	restrictedDirNamesInWindows = []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	restrictedRegexInWindows = []regexp.Regexp{
		*regexp.MustCompile(`(?m)\s$`),
	}
)

func CleanDirName(dir string) string {
	return ""
}
