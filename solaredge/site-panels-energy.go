package solaredge

import (
	"fmt"
)

type logicalLayout struct {
	LogicalTree struct {
		Children []struct {
			Children []struct {
				Children []struct {
					Panel struct {
						ID           int    `json:"id"`
						SerialNumber string `json:"serialNumber"`
						DisplayName  string `json:"displayName"`
					} `json:"data"`
				} `json:"children"`
			} `json:"children"`
		} `json:"children"`
	} `json:"logicalTree"`
	ReportersData map[int]struct {
		Energy float64 `json:"unscaledEnergy"`
	} `json:"reportersData"`
}

type PanelData struct {
	ID           int
	SerialNumber string
	DisplayName  string
	Energy       float64
}

func (s *SiteService) PanelsEnergy(siteID string) (*[]PanelData, error) {
	u := fmt.Sprintf("%s/sites/%s/layout/logical", defaultBaseURL, siteID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var panelData logicalLayout
	_, err = s.client.do(req, &panelData)
	if err != nil {
		return nil, err
	}

	var panels []PanelData
	for _, b := range panelData.LogicalTree.Children {
		for _, d := range b.Children {
			for _, f := range d.Children {
				panel := PanelData{
					ID:           f.Panel.ID,
					SerialNumber: f.Panel.SerialNumber,
					DisplayName:  f.Panel.DisplayName,
				}
				// Loop over reporterData, find a matching ID and extract the energy figure.
				for g, h := range panelData.ReportersData {
					if panel.ID == g {
						panel.Energy = h.Energy
						panels = append(panels, panel)
						break
					}
				}
			}
		}
	}

	return &panels, err
}
