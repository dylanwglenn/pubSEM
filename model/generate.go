package model

import (
	"image/color"
	"main/utils"
)

func InitTestModel() *Model {

	nodeA := &Node{
		Class:       OBSERVED,
		Pos:         utils.LocalPos{100, 100},
		Dim:         utils.LocalDim{50, 50},
		Col:         color.NRGBA{0, 255, 0, 255},
		Thickness:   2.0,
		UserDefined: true,
	}

	nodeB := &Node{
		Class:       OBSERVED,
		Pos:         utils.LocalPos{500, 200},
		Dim:         utils.LocalDim{70, 50},
		Col:         color.NRGBA{0, 0, 255, 255},
		Thickness:   2.0,
		UserDefined: false,
	}

	nodeC := &Node{
		Class:       OBSERVED,
		Pos:         utils.LocalPos{-50, 100},
		Dim:         utils.LocalDim{70, 50},
		Col:         color.NRGBA{0, 0, 255, 255},
		Thickness:   2.0,
		UserDefined: false,
	}

	connectionA := &Connection{
		Origin:      nodeA,
		Destination: nodeB,
		Col:         color.NRGBA{0, 0, 0, 255},
		Thickness:   2.0,
		Type:        REGRESSION,
		UserDefined: true,
	}

	connectionB := &Connection{
		Origin:      nodeC,
		Destination: nodeA,
		Col:         color.NRGBA{0, 0, 0, 255},
		Thickness:   2.0,
		Type:        COVARIANCE,
		UserDefined: true,
	}

	return &Model{
		Nodes:       []*Node{nodeA, nodeB, nodeC},
		Connections: []*Connection{connectionA, connectionB},
	}
}
