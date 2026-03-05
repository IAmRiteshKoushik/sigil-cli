package main

import (
	"github.com/IAmRiteshKoushik/sigil/cmd"
	"github.com/IAmRiteshKoushik/sigil/pkg"
)

func init() {
	pkg.LoadConfig()
	pkg.SetupLogger()
}

func main() {
	cmd.Execute()
}
