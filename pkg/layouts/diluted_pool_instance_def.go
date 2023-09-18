package layouts

type DilutedPoolInstanceDef struct {
	UnitsPerStep uint
	Spacing      uint
	NBits        uint
}

func DefaultDilutedPoolInstance() *DilutedPoolInstanceDef {
	return &DilutedPoolInstanceDef{UnitsPerStep: 16, Spacing: 4, NBits: 16}
}
