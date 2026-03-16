package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/observation"
	"arandu/internal/infrastructure/repository/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("🔍 Validação Completa da Implementação REQ-01-02-02")
	fmt.Println("====================================================")
	fmt.Println()

	// 1. Configurar banco de dados em memória para teste
	fmt.Println("1. Configurando ambiente de teste...")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteDB := &sqlite.DB{db}

	// 2. Inicializar repositório
	fmt.Println("2. Inicializando repositório...")
	repo := sqlite.NewObservationRepository(sqliteDB)
	if err := repo.InitSchema(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✅ Schema inicializado com campo updated_at")

	// 3. Criar serviço
	fmt.Println("3. Criando serviço...")
	service := services.NewObservationService(repo)

	// 4. Testar criação de observação
	fmt.Println("4. Testando criação de observação...")
	obs, err := service.CreateObservation("session-test-123", "Observação clínica inicial para validação")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Observação criada: ID=%s\n", obs.ID)
	fmt.Printf("   📝 Conteúdo: %s\n", obs.Content)
	fmt.Printf("   📅 Criada em: %s\n", obs.CreatedAt.Format("2006-01-02 15:04:05"))

	// 5. Testar recuperação da observação
	fmt.Println("5. Testando recuperação da observação...")
	retrieved, err := service.GetObservation(obs.ID)
	if err != nil {
		log.Fatal(err)
	}
	if retrieved == nil {
		log.Fatal("Observação não encontrada após criação")
	}
	fmt.Println("   ✅ Observação recuperada com sucesso")

	// 6. Testar atualização da observação
	fmt.Println("6. Testando atualização da observação...")
	originalCreatedAt := retrieved.CreatedAt
	time.Sleep(10 * time.Millisecond) // Garantir diferença de tempo

	err = service.UpdateObservation(obs.ID, "Observação clínica ATUALIZADA após edição")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✅ Observação atualizada com sucesso")

	// 7. Verificar atualização
	fmt.Println("7. Verificando dados atualizados...")
	updated, err := service.GetObservation(obs.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar conteúdo atualizado
	if updated.Content != "Observação clínica ATUALIZADA após edição" {
		log.Fatalf("Conteúdo não atualizado: %s", updated.Content)
	}
	fmt.Println("   ✅ Conteúdo atualizado corretamente")

	// Verificar created_at preservado
	if !updated.CreatedAt.Equal(originalCreatedAt) {
		log.Fatal("created_at foi alterado (não deveria)")
	}
	fmt.Println("   ✅ created_at preservado (não alterado)")

	// Verificar updated_at definido
	if updated.UpdatedAt.IsZero() {
		log.Fatal("updated_at não foi definido")
	}
	fmt.Printf("   ✅ updated_at definido: %s\n", updated.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Verificar que updated_at é posterior a created_at
	if !updated.UpdatedAt.After(updated.CreatedAt) {
		log.Fatal("updated_at não é posterior a created_at")
	}
	fmt.Println("   ✅ updated_at é posterior a created_at")

	// 8. Testar validações
	fmt.Println("8. Testando validações...")

	// Testar conteúdo vazio
	err = service.UpdateObservation(obs.ID, "")
	if err == nil {
		log.Fatal("Update com conteúdo vazio deveria falhar")
	}
	fmt.Println("   ✅ Validação: conteúdo não pode ser vazio")

	// Testar conteúdo muito longo
	longContent := ""
	for i := 0; i < 5001; i++ {
		longContent += "a"
	}
	err = service.UpdateObservation(obs.ID, longContent)
	if err == nil {
		log.Fatal("Update com conteúdo muito longo deveria falhar")
	}
	fmt.Println("   ✅ Validação: conteúdo não pode exceder 5000 caracteres")

	// Testar observação não existente
	err = service.UpdateObservation("id-inexistente", "Conteúdo válido")
	if err == nil {
		log.Fatal("Update de observação inexistente deveria falhar")
	}
	fmt.Println("   ✅ Validação: observação deve existir")

	// 9. Testar repositório diretamente
	fmt.Println("9. Testando repositório diretamente...")

	// Criar observação via repositório
	repoObs := &observation.Observation{
		SessionID: "session-direct-456",
		Content:   "Observação criada diretamente no repositório",
	}

	if err := repo.Save(repoObs); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Observação salva via repositório: ID=%s\n", repoObs.ID)

	// Atualizar via repositório
	repoObs.Content = "Observação atualizada via repositório"
	if err := repo.Update(repoObs); err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✅ Observação atualizada via repositório")

	// Verificar updated_at no banco
	var updatedAt sql.NullTime
	err = db.QueryRow("SELECT updated_at FROM observations WHERE id = ?", repoObs.ID).Scan(&updatedAt)
	if err != nil {
		log.Fatal(err)
	}

	if !updatedAt.Valid {
		log.Fatal("updated_at não foi definido no banco")
	}
	fmt.Printf("   ✅ updated_at persistido no banco: %s\n", updatedAt.Time.Format("2006-01-02 15:04:05"))

	// 10. Resumo
	fmt.Println()
	fmt.Println("🎉 VALIDAÇÃO COMPLETA BEM-SUCEDIDA!")
	fmt.Println("====================================")
	fmt.Println()
	fmt.Println("✅ Todas as camadas validadas:")
	fmt.Println("   1. ✅ Domínio: Struct Observation com campo UpdatedAt")
	fmt.Println("   2. ✅ Infraestrutura: Repositório com método Update")
	fmt.Println("   3. ✅ Aplicação: Serviço com validações completas")
	fmt.Println("   4. ✅ Persistência: Campo updated_at no banco de dados")
	fmt.Println()
	fmt.Println("✅ Validações implementadas:")
	fmt.Println("   - Conteúdo não pode ser vazio")
	fmt.Println("   - Conteúdo máximo 5000 caracteres")
	fmt.Println("   - Observação deve existir")
	fmt.Println("   - created_at preservado")
	fmt.Println("   - updated_at sempre definido na modificação")
	fmt.Println()
	fmt.Println("✅ Fluxo de dados:")
	fmt.Println("   created_at → preservado (imutável)")
	fmt.Println("   updated_at → sempre atualizado na modificação")
	fmt.Println("   updated_at > created_at → garantido")
	fmt.Println()
	fmt.Println("🎯 Próximo passo: Testar interface web em http://localhost:8080")
	fmt.Println("   - Criar uma observação")
	fmt.Println("   - Clicar no ícone de lápis para editar")
	fmt.Println("   - Testar edição e cancelamento")
}
