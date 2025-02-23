package services

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"odonto-reports/models"
	"strconv"
)

// Processa o arquivo Excel e retorna os dados estruturados
func ProcessarExcel() (models.Relatorio, error) {
	// Abrindo o arquivo Excel
	file, err := excelize.OpenFile("uploads/relatorio.xlsx")
	if err != nil {
		return models.Relatorio{}, fmt.Errorf("erro ao abrir o arquivo Excel: %v", err)
	}
	defer file.Close()

	// Lendo os dias úteis, corridos e restantes
	diasUteis, _ := lerCelulaComoInt(file, "DIAS_TRABALHO", "C5")
	diasCorridos, _ := lerCelulaComoInt(file, "DIAS_TRABALHO", "C4")
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

		percentReal := (real / meta) * 100
		projecao := (real / float64(diasCorridos)) * float64(diasUteis)
		percentProj := (projecao / meta) * 100

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
