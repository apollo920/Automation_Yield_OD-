package services

import (
	"fmt"
	"odonto-reports/models"
)

// Gera um texto resumido do relatório
func GerarRelatorio(dados models.Relatorio) string {
	relatorio := fmt.Sprintf(
		"📊 *Relatório Financeiro*\n\n"+
			"📅 *Dias Corridos:* %d\n"+
			"📅 *Dias Restantes:* %d\n"+
			"📅 *Dias Úteis:* %d\n\n"+
			"💰 *Pilares:* \n",
		dados.DiasCorridos, dados.DiasFaltam, dados.DiasUteis,
	)

	for nome, pilar := range dados.Pilares {
		relatorio += fmt.Sprintf(
			"🔹 *%s*\n"+
				"   ✅ Real: R$ %.2f\n"+
				"   🎯 Meta: R$ %.2f\n"+
				"   📈 Percentual Real: %.2f%%\n"+
				"   🔮 Projeção: R$ %.2f\n"+
				"   📊 Percentual Projeção: %.2f%%\n\n",
			nome, pilar.Real, pilar.Meta, pilar.PercentReal, pilar.Projecao, pilar.PercentProj,
		)
	}

	return relatorio
}
