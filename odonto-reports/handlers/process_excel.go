package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"odonto-reports/models"
	"odonto-reports/services"
	"strconv"
	"sync"
)

// Estrutura para armazenar os dados extraídos
type ReportData struct {
	DiasUteis    int
	DiasCorridos int
	DiasFaltam   int
	Pilares      map[string]PilarData
}

// Estrutura para armazenar os dados de cada Pilar
type PilarData struct {
	Real        float64
	Meta        float64
	PercentReal float64
	Projecao    float64
	PercentProj float64
}

// Função para processar o Excel
func ProcessExcel(filePath string) (*ReportData, error) {
	log.Println("Abrindo o arquivo Excel:", filePath)

	// Tenta abrir o arquivo
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println("Erro ao abrir o arquivo Excel:", err)
		return nil, fmt.Errorf("erro ao abrir o arquivo Excel: %v", err)
	}
	defer f.Close()

	// Supondo que "f" seja seu arquivo Excel (workbook)
	sheets := f.GetSheetList()
	log.Println("Planilhas encontradas:", sheets)
	log.Println("Arquivo Excel aberto com sucesso!")

	// -------- Ler os dados da aba "DIAS_TRABALHO" --------
	diasUteis, err := f.GetCellValue("DIAS_TRABALHO", "A2")
	if err != nil {
		log.Fatalf("Erro ao ler A2: %v", err)
	}
	if diasUteis == "" {
		log.Fatal("A célula A2 está vazia!")
	}

	diasCorridos, err := f.GetCellValue("DIAS_TRABALHO", "B2")
	if err != nil {
		log.Fatalf("Erro ao ler B2: %v", err)
	}
	if diasCorridos == "" {
		log.Fatal("A célula B2 está vazia!")
	}

	diasFaltam, err := f.GetCellValue("DIAS_TRABALHO", "C2")
	if err != nil {
		log.Fatalf("Erro ao ler C2: %v", err)
	}
	if diasFaltam == "" {
		log.Fatal("A célula C2 está vazia!")
	}

	log.Println("Valores lidos:", diasUteis, diasCorridos, diasFaltam)

	diasCorridosInt, err := strconv.Atoi(diasCorridos)
	if err != nil {
		return nil, fmt.Errorf("erro ao retornar diasCorridos")

	}

	diasUteisInt, err := strconv.Atoi(diasUteis)
	if err != nil {
		return nil, fmt.Errorf("erro ao retornar diasUteis")

	}

	diasFaltamInt, err := strconv.Atoi(diasFaltam)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter diasFaltam: %v", err)
	}

	// -------- Ler os dados da aba "CONTROLE" --------
	rows, err := f.GetRows("CONTROLE")

	if err != nil {
		log.Fatalf("Erro ao obter as linhas da aba CONTROLE: %v", err)
	}
	log.Printf("Total de linhas na aba CONTROLE: %d", len(rows))

	for i, row := range rows {
		log.Printf("Linha %d tem %d colunas: %+v", i, len(row), row)
	}
	// Verifique se há linhas antes de acessá-las
	if len(rows) == 0 {
		log.Println("Erro: Nenhuma linha encontrada no arquivo Excel")
		return nil, fmt.Errorf("o arquivo Excel não contém dados")
	}

	if len(rows[0]) < 2 { // Supondo que cada linha deve ter pelo menos 2 colunas
		log.Println("Erro: O arquivo Excel não tem colunas suficientes")
		return nil, fmt.Errorf("o arquivo Excel não contém colunas suficientes")
	}

	if err != nil {
		log.Println("Erro ao ler a aba CONTROLE:", err)
		return nil, err
	}

	// Criar mapa para armazenar os dados dos pilares
	pilares := make(map[string]PilarData)

	// Percorrer as linhas (começando da linha 2 para ignorar cabeçalhos)
	for i := 1; i < len(rows); i++ {
		pilar := rows[i][0] // Nome do Pilar (ex: ORTO, CREDIARIO, etc.)

		real, _ := parseFloat(rows[i][1]) // Coluna "Real"
		meta, _ := parseFloat(rows[i][2]) // Coluna "Meta"

		// Cálculo de Percentual Real/Meta
		percentReal := (real / meta) * 100

		// Cálculo da Projeção
		projecao := (real / float64(diasCorridosInt)) * float64(diasUteisInt)

		// Cálculo de Percentual Projeção/Meta
		percentProj := (projecao / meta) * 100

		// Adicionar ao mapa
		pilares[pilar] = PilarData{
			Real:        real,
			Meta:        meta,
			PercentReal: percentReal,
			Projecao:    projecao,
			PercentProj: percentProj,
		}
	}

	// Retornar os dados extraídos
	return &ReportData{
		DiasUteis:    diasUteisInt,
		DiasCorridos: diasCorridosInt,
		DiasFaltam:   diasFaltamInt,
		Pilares:      pilares,
	}, nil
}

// Função auxiliar para converter string em float
func parseFloat(value string) (float64, error) {
	if value == "" {
		return 0, nil
	}
	var num float64
	_, err := fmt.Sscanf(value, "%f", &num)
	return num, err
}

var (
	relatorioGerado models.Relatorio // Armazena os dados processados
	mu              sync.Mutex       // Mutex para evitar concorrência
)

// Processa o Excel e salva os dados globalmente
func ProcessExcelHandler(w http.ResponseWriter, r *http.Request) {
	dadosProcessados, err := services.ProcessarExcel()
	if err != nil {
		http.Error(w, "Erro ao processar o Excel", http.StatusInternalServerError)
		return
	}

	// Salvar os dados no cache global
	mu.Lock()
	relatorioGerado = dadosProcessados
	mu.Unlock()

	// Retorna os dados processados
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Arquivo processado com sucesso",
		"data":    dadosProcessados,
	})
}
