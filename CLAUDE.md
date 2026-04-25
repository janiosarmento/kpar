# KPAR — Instruções para Claude

## Build e instalação

Sempre que compilar/instalar o programa, incremente a versão em `-ldflags`:

```sh
go install -ldflags "-X main.version=X.Y.Z" ./cmd/kpar
/bin/cp -f ~/go/bin/kpar ~/.local/bin/kpar
```

A variável `version` fica em `cmd/kpar/main.go` e o padrão é `"dev"` quando não injetada.
