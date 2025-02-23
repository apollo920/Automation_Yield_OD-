package handlers

import (
	"encoding/json"
	"net/http"
	"odonto-reports/services"
)

// Retorna um relatório formatado com base nos dados processados
func GerarRelatorioHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Verifica se há dados processados
	if relatorioGerado.DiasUteis == 0 {
		http.Error(w, "Nenhum dado processado ainda", http.StatusBadRequest)
		return
	}

	// Gera o relatório
	relatorio := services.GerarRelatorio(relatorioGerado)

	// Retorna o relatório
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"relatorio": relatorio})
}
