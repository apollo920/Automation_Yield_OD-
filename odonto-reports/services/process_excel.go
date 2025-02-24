package services

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"odonto-reports/models"
	"strconv"
	"math"
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

	for _, linha := range linhas[1:] { // Pulando o cabeçalho
		if len(linha) < 5 {
			continue
		}
	
		nomePilar := linha[0]
		real, _ := strconv.ParseFloat(linha[1], 64)
		meta, _ := strconv.ParseFloat(linha[2], 64)
	
		// Verifica se meta é zero para evitar divisão por zero
		var percentReal, projecao, percentProj float64
		if meta != 0 {
			percentReal = (real / meta) * 100
			projecao = (real / float64(diasCorridos)) * float64(diasUteis)
			percentProj = (projecao / meta) * 100
		} else {
			// Define valores padrão ou trata o caso de meta zero
			percentReal = 0
			projecao = 0
			percentProj = 0
		}
	
		pilares[nomePilar] = struct {
			Real        float64
			Meta        float64
			PercentReal float64
			Projecao    float64
			PercentProj float64
		}{
			Real:        replaceNaNInf(real),
			Meta:        replaceNaNInf(meta),
			PercentReal: replaceNaNInf(percentReal),
			Projecao:    replaceNaNInf(projecao),
			PercentProj: replaceNaNInf(percentProj),
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

func replaceNaNInf(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return value
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
