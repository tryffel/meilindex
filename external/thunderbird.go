package external

import "os/exec"

func OpenById(id string) error {
	cmd := exec.Command("thunderbird", "-thunderlink",
		"thunderlink://messageid="+id)
	return cmd.Run()
}
