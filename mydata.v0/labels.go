package mydata

// Label represents a label definition as returned from the API
type Label struct {
	Type string `json:"type"`
	Name string `json:"label"`
}

// Labels is just a collection of Label structures
type Labels []Label

// Compartments returns the compartment labels
func (l Labels) Compartments() Labels {
	lbls := Labels{}
	for _, lbl := range l {
		if lbl.Type == "compartment" {
			lbls = append(lbls, lbl)
		}
	}

	return lbls
}

// Irrigators returns the irrigation system labels
func (l Labels) Irrigators() Labels {
	lbls := Labels{}
	for _, lbl := range l {
		if lbl.Type == "irrigator" {
			lbls = append(lbls, lbl)
		}
	}

	return lbls
}

// Monitors returns the monitor labels
func (l Labels) Monitors() Labels {
	lbls := Labels{}
	for _, lbl := range l {
		if lbl.Type == "monitor" {
			lbls = append(lbls, lbl)
		}
	}

	return lbls
}

// GrowRooms returns the grow room labels
func (l Labels) GrowRooms() Labels {
	lbls := Labels{}
	for _, lbl := range l {
		if lbl.Type == "growroom" {
			lbls = append(lbls, lbl)
		}
	}

	return lbls
}

// Fields returns the field labels
func (l Labels) Fields() Labels {
	lbls := Labels{}
	for _, lbl := range l {
		if lbl.Type == "field" {
			lbls = append(lbls, lbl)
		}
	}

	return lbls
}
