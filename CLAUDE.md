# KPAR — Instruções para Claude

## Build e instalação

Sempre que compilar/instalar o programa, use `make install`. Esse alvo faz o
build e o deploy (`cp` para `~/.local/bin/kpar`) num único passo atômico, e
incrementa sozinho o patch da versão em `VERSION`. **Nunca** rode
`go install`/`go build` isolado para instalar — isso deixa `~/go/bin/kpar` e
`~/.local/bin/kpar` dessincronizados, como já aconteceu antes.

```sh
make install
```

A variável `version` fica em `cmd/kpar/main.go` e o padrão é `"dev"` quando não injetada.
