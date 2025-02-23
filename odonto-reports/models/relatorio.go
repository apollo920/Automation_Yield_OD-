package models

type Relatorio struct {
	DiasCorridos int
	DiasFaltam   int
	DiasUteis    int
	Pilares      map[string]struct {
		Real        float64
		Meta        float64
		PercentReal float64
		Projecao    float64
		PercentProj float64
	}
}
