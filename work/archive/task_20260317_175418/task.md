Guia de Compatibilidade: Habilitação de SQLite FTS5

Contexto: O objetivo deste prompt é garantir que o ambiente de desenvolvimento possua um binário/driver de SQLite capaz de executar comandos FTS5 (Full Text Search). Focaremos apenas na compatibilidade do motor antes de aplicar as regras ao Arandu.

⚡ Passo 1: A Escolha do Driver (CGO vs. Pure Go)

O erro "unknown module: fts5" no Go ocorre quase sempre por problemas na compilação do driver go-sqlite3. Para resolver isso de forma definitiva, avalie as duas rotas abaixo:

Rota A (Recomendada SOTA): Driver "Pure Go"

A forma mais robusta de garantir FTS5 em qualquer ambiente sem depender de compiladores C ou bibliotecas do sistema é usar o driver modernc.org/sqlite.

Vantagem: É escrito 100% em Go. Já vem com FTS5 habilitado por padrão.

Instalação: go get modernc.org/sqlite

Uso: Importe como import _ "modernc.org/sqlite" e use o driver name "sqlite".

Rota B: Driver "CGO" (go-sqlite3)

Se for obrigatório usar github.com/mattn/go-sqlite3:

Certificação: O comando CGO_ENABLED=1 deve estar ativo.

Tags: A compilação EXIGE -tags "fts5". Se falhou, verifique se o binário do SQLite instalado no seu SO (sqlite3 --version) suporta FTS5 rodando PRAGMA compile_options;.

📂 Passo 2: Script de Teste de Sanidade (PoC)

Antes de alterar o código do Arandu, crie um arquivo isolado test_fts5.go para validar se o seu motor de banco de dados "enxerga" o módulo FTS5.

package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite" // Ou "[github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)" se optar pela Rota B
)

func main() {
	// Use "sqlite" para modernc ou "sqlite3" para mattn
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// O Teste Real: Tentar criar uma tabela virtual
	_, err = db.Exec("CREATE VIRTUAL TABLE busca USING fts5(conteudo);")
	if err != nil {
		fmt.Printf("❌ FTS5 NÃO COMPATÍVEL: %v\n", err)
		return
	}

	fmt.Println("✅ SUCESSO: O motor SQLite suporta FTS5!")
}


🔍 Passo 3: Verificação de Recursos do Sistema

Se o script acima falhar, execute estes comandos no terminal para diagnosticar o binário local:

Verificar Módulos:

sqlite3 :memory: "PRAGMA compile_options;" | grep FTS5


Se não retornar ENABLE_FTS5, o binário do seu sistema operacional é antigo/limitado.

Verificar CGO (Apenas para Rota B):

go env CGO_ENABLED


🛡️ Checklist de Validação para o Agente

[ ] O script test_fts5.go rodou com sucesso?

[ ] Se o go-sqlite3 falhou, a migração para modernc.org/sqlite foi realizada?

[ ] O comando CREATE VIRTUAL TABLE ... USING fts5 funciona em uma base em memória?

Atenção: Não prossiga com a implementação das migrações do Arandu até que este teste de sanidade retorne "SUCESSO".