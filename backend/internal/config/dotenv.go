package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadDotEnv carrega pares KEY=VALUE de um arquivo .env para o ambiente
// do processo, SEM sobrescrever variáveis já definidas. É uma
// conveniência de desenvolvimento local; em produção as variáveis vêm do
// systemd/ambiente e o arquivo .env nem existe (está no .gitignore).
//
// Suporta: linhas em branco, comentários com '#', aspas simples/duplas
// ao redor do valor e o prefixo opcional "export ".
func LoadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // ausência de .env não é erro
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)

		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
	return sc.Err()
}
