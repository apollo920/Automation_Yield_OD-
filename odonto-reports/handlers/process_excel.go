package handlers

import (
	"fmt"
	"log"
	"github.com/xuri/excelize/v2"
)

// Estrutura para armazenar os dados extraídos
type ReportData struct {
	DiasUteis       int
	DiasCorridos    int
	DiasFaltam      int
	Pilares         map[string]PilarData
}

// Estrutura para armazenar os dados de cada Pilar
type PilarData struct {
	Real          float64
	Meta          float64
	PercentReal   float64
	Projecao      float64
	PercentProj   float64
}

// Função para processar o Excel
func ProcessExcel(filePath string) (*ReportData, error) {
	// Abre o arquivo Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println("Erro ao abrir o arquivo Excel:", err)
		return nil, err
	}
	defer f.Close()

	// -------- Ler os dados da aba "DIAS_TRABALHO" --------
	diasUteis, err := f.GetCellValue("DIAS_TRABALHO", "C2")
	diasCorridos, err := f.GetCellValue("DIAS_TRABALHO", "C3")
	diasFaltam, err := f.GetCellValue("DIAS_TRABALHO", "C4")

	if err != nil {
		log.Println("Erro ao ler os dados de DIAS_TRABALHO:", err)
		return nil, err
	}

	diasCorridosInt, err := strconv.Atoi(diasCorridos)
	if err != nil {
    return fmt.Errorf("erro ao converter diasCorridos: %v", err)
}

	diasUteisInt, err := strconv.Atoi(diasUteis)
	if err != nil {
	return fmt.Errorf("erro ao converter diasUteis: %v", err)
}

	// -------- Ler os dados da aba "CONTROLE" --------
	rows, err := f.GetRows("CONTROLE")
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
		DiasUteis:    int(diasUteis),
		DiasCorridos: int(diasCorridos),
		DiasFaltam:   int(diasFaltam),
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
