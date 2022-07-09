package main

import (
	"github.com/canalun/sqloth/cmd"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	cmd.Execute()
}
