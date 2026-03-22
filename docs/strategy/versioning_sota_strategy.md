Estratégia de Versionamento e Build-time Metadata

O Arandu utiliza metadados injetados durante o tempo de compilação para garantir a rastreabilidade total dos logs e erros.

1. Injeção via LDFLAGS

Não utilizamos ficheiros de configuração para a versão. A versão é definida pelas Git Tags e injetada no pacote internal/platform/version.

O Pacote Go

Criamos um ficheiro internal/platform/version/version.go:

package version

var (
    Version   = "dev"
    Commit    = "none"
    BuildTime = "unknown"
)


O Comando de Build

O agente deve estar ciente de que o build (ou execução) deve preferencialmente usar:

go run -ldflags "-X '[github.com/seu-user/arandu/internal/platform/version.Version=$(git](https://github.com/seu-user/arandu/internal/platform/version.Version=$(git) describe --tags --always)' -X '[github.com/seu-user/arandu/internal/platform/version.Commit=$(git](https://github.com/seu-user/arandu/internal/platform/version.Commit=$(git) rev-parse --short HEAD)'" cmd/arandu/main.go


2. Uso em Logs

O logger centralizado importa este pacote e injeta automaticamente os campos v (version) e commit em cada linha de log JSON. Isto permite filtrar no Grafana: "Mostre-me apenas logs da versão v1.2.0".