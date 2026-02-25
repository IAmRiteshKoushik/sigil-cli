build:
  @go build -ldflags="-X 'github.com/IAmRiteshKoushik/sigil/cmd.Version=$(git describe --tags)'" -o sigil main.go
