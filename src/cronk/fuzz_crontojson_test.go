package cronk

import (
	"testing"
)

func FuzzCronToJson(f *testing.F) {
	f.Add("# comment\n0 0 * * * echo test")
	f.Add("0 * * * * echo hello\n# outro")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		_ = CronToJson(input)
	})
}
