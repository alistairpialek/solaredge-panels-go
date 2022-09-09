package solaredge

type PanelData struct {
	ID           int
	SerialNumber string
	DisplayName  string
	Energy       float64
}

func (s *SitesService) PanelsEnergy() (PanelData, error) {
	p := PanelData{}
	return p, nil
}
