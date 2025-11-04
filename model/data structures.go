package model

import (
	"image/color"
	"main/utils"
)

type ParamType int

const (
	OBSERVED ParamType = iota
	LATENT
)

type ConnectionType int

const (
	REGRESSION ConnectionType = iota
	COVARIANCE
)

type Node struct {
	Class       ParamType
	Pos         utils.LocalPos
	Dim         utils.LocalDim
	Col         color.NRGBA
	Text        string
	Thickness   float32
	UserDefined bool
}

type Connection struct {
	Origin      *Node
	Destination *Node
	Col         color.NRGBA
	Thickness   float32
	Type        ConnectionType
	Text        string
	UserDefined bool
}

type Model struct {
	Nodes       []*Node
	Connections []*Connection
}
