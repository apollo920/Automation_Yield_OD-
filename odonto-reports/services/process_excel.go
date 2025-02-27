package services

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"odonto-reports/models"
	"strconv"
	"log"

)

// Processa o arquivo Excel e retorna os dados estruturados
func ProcessarExcel(filepath string) (models.Relatorio, error) {
	// Abrindo o arquivo Excel
	file, err := excelize.OpenFile(filepath)
	if err != nil {
		return models.Relatorio{}, fmt.Errorf("erro ao abrir o arquivo Excel: %v", err)
	}
	defer file.Close()

	// Lendo os dias úteis, corridos e restantes
	diasUteis, _ := lerCelulaComoInt(file, "DIAS_TRABALHO", "A2")
	diasCorridos, _ := lerCelulaComoInt(file, "DIAS_TRABALHO", "B2")
	diasFaltam := diasUteis - diasCorridos

	// Lendo os valores da aba CONTROLE
	pilares := map[string]struct {
		Real        float64
		Meta        float64
		PercentReal float64
		Projecao    float64
		PercentProj float64
	}{}

	linhas, err := file.GetRows("CONTROLE")
	if err != nil {
		return models.Relatorio{}, fmt.Errorf("erro ao ler a aba CONTROLE: %v", err)
	}

	// Log para depuração
	log.Println("Linhas lidas da aba CONTROLE:")
	for i, linha := range linhas {
		log.Printf("Linha %d: %v", i, linha)
	}	

	for _, linha := range linhas[1:] { // Pulando o cabeçalho
		if len(linha) < 5 {
			continue
		}
	
		nomePilar := linha[0]
		realStr := linha[1]
		metaStr := linha[2]
	
		// Log para depuração
		log.Printf("Processando linha: Nome=%s, Real=%s, Meta=%s", nomePilar, realStr, metaStr)
	
		// Converte os valores para float64
		real, err := strconv.ParseFloat(realStr, 64)
		if err != nil {
			log.Printf("Erro ao converter Real (%s) para float64: %v", realStr, err)
			real = 0
		}
	
		meta, err := strconv.ParseFloat(metaStr, 64)
		if err != nil {
			log.Printf("Erro ao converter Meta (%s) para float64: %v", metaStr, err)
			meta = 0
		}
	
		// Cálculos
		var percentReal, projecao, percentProj float64
		if meta != 0 {
			percentReal = (real / meta) * 100
			projecao = (real / float64(diasCorridos)) * float64(diasUteis)
			percentProj = (projecao / meta) * 100
		}
	
		
		// Adiciona os dados ao mapa de pilares
		pilares[nomePilar] = struct {
			Real        float64
			Meta        float64
			PercentReal float64
			Projecao    float64
			PercentProj float64
		}{
			Real:        real,
			Meta:        meta,
			PercentReal: percentReal,
			Projecao:    projecao,
			PercentProj: percentProj,
		}
	}

	// Retornando os dados estruturados
	return models.Relatorio{
		DiasCorridos: diasCorridos,
		DiasFaltam:   diasFaltam,
		DiasUteis:    diasUteis,
		Pilares:      pilares,
	}, nil

}
// Função auxiliar para ler valores inteiros de células
func lerCelulaComoInt(file *excelize.File, aba string, celula string) (int, error) {
	valor, err := file.GetCellValue(aba, celula)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(valor)
	if err != nil {
		return 0, err
	}
	return num, nil
}
