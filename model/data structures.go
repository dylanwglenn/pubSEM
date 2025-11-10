package model

import (
	"image/color"
	"main/utils"
)

type ParamType int

const (
	OBSERVED ParamType = iota
	LATENT
	INTERCEPT //TODO: add support for intercepts (triangle nodes)
)

type ConnectionType int

const (
	REGRESSION ConnectionType = iota
	COVARIANCE
)

type Node struct {
	Class           ParamType
	Pos             utils.LocalPos
	Dim             utils.LocalDim
	Col             color.NRGBA
	Text            string
	Thickness       float32
	UserDefined     bool
	EdgeConnections [4][]*Connection // only applicable for observed (rectangular) nodes
	Padding         float32
}

type Connection struct {
	Origin         *Node
	Destination    *Node
	OriginPos      utils.LocalPos
	DestinationPos utils.LocalPos
	Angle          float64
	Col            color.NRGBA
	Thickness      float32
	Type           ConnectionType
	EstPos         utils.LocalPos
	EstDim         utils.LocalDim
	Est            float64
	PValue         float64
	CI             [2]float64
	EstText        string
	Curvature      float32 // only applicable for covariance (curved) connections
	UserDefined    bool
}

type Model struct {
	Nodes       []*Node
	Connections []*Connection
}
