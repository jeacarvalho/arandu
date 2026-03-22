package version

// Version é a versão da aplicação (injetada via ldflags)
var Version = "dev"

// Commit é o hash do commit git (injetado via ldflags)
var Commit = "unknown"

// BuildTime é o timestamp do build (injetado via ldflags)
var BuildTime = "unknown"

// AppName é o nome da aplicação
const AppName = "arandu"
