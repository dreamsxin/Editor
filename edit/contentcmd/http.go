package contentcmd

import (
	"net/url"
	"os/exec"

	"github.com/jmigpin/editor/edit/cmdutil"
)

// Opens http/https lines in x-www-browser.
func http(erow cmdutil.ERower, s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		return false
	}
	go func() {
		cmd := exec.Command("x-www-browser", u.String())
		err := cmd.Run()
		if err != nil {
			ed := erow.Editorer()
			ed.Error(err)
			ed.UI().RequestTreePaint()
		}
	}()
	return true
}
