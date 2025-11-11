package model

import (
	"image/color"
	"main/utils"
)

const (
	roundness     = .3
	estPadding    = 2.5
	propAlongLine = .5
)

func InitTestModel() *Model {
	nodeA := &Node{
		Class:       OBSERVED,
		Text:        "khalksfdhlkjsahdflksahfsafd",
		Pos:         utils.LocalPos{100, 160},
		Col:         color.NRGBA{255, 255, 255, 255},
		Thickness:   3.0,
		UserDefined: true,
	}

	nodeB := &Node{
		Class:       LATENT,
		Text:        "Mental Health 1",
		Pos:         utils.LocalPos{500, 200},
		Col:         color.NRGBA{255, 255, 255, 255},
		Thickness:   3.0,
		UserDefined: false,
	}

	nodeC := &Node{
		Class:       OBSERVED,
		Text:        "Migration",
		Pos:         utils.LocalPos{-40, 100},
		Col:         color.NRGBA{255, 255, 255, 255},
		Thickness:   3.0,
		UserDefined: false,
	}

	connectionA := &Connection{
		Origin:        nodeA,
		Destination:   nodeB,
		Col:           color.NRGBA{0, 0, 0, 255},
		Thickness:     2.0,
		Type:          REGRESSION,
		UserDefined:   true,
		Curvature:     roundness,
		EstPadding:    estPadding,
		AlongLineProp: propAlongLine,
	}

	connectionB := &Connection{
		Origin:        nodeC,
		Destination:   nodeA,
		Col:           color.NRGBA{0, 0, 0, 255},
		Thickness:     2.0,
		Type:          COVARIANCE,
		Est:           .01234656213,
		PValue:        .00001,
		UserDefined:   true,
		Curvature:     roundness,
		EstPadding:    estPadding,
		AlongLineProp: propAlongLine,
	}

	connectionC := &Connection{
		Origin:        nodeB,
		Destination:   nodeC,
		Col:           color.NRGBA{0, 0, 0, 255},
		Thickness:     2.0,
		Type:          REGRESSION,
		Est:           .13486,
		PValue:        .049,
		UserDefined:   true,
		Curvature:     roundness,
		EstPadding:    estPadding,
		AlongLineProp: propAlongLine,
	}

	return &Model{
		Nodes:       []*Node{nodeA, nodeB, nodeC},
		Connections: []*Connection{connectionA, connectionB, connectionC},
	}
}
