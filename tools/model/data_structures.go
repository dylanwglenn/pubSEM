package model

import (
	"image/color"
	"main/utils"

	"gioui.org/font"
)

type ParamType int

const (
	OBSERVED ParamType = iota
	LATENT
	INTERCEPT //TODO: add support for intercepts (triangle nodes)
)

type ConnectionType int

const (
	STRAIGHT ConnectionType = iota
	CURVED
	CIRCULAR
)

type Node struct {
	Class           ParamType        `json:"class,omitempty"`
	Pos             utils.LocalPos   `json:"pos"`
	Dim             utils.LocalDim   `json:"dim"`
	Col             color.NRGBA      `json:"col"`
	VarName         string           `json:"var_name,omitempty"`
	Text            string           `json:"text,omitempty"`
	TextWidth       float32          `json:"text_width,omitempty"`
	Bold            bool             `json:"bold,omitempty"`
	Thickness       float32          `json:"thickness,omitempty"`
	UserDefined     bool             `json:"user_defined,omitempty"`
	Visible         bool             `json:"visible,omitempty"`
	EdgeConnections [4][]*Connection `json:"-"` // only applicable for rectangular nodes
	Padding         float32          `json:"padding,omitempty"`
}

type Connection struct {
	Origin         *Node          `json:"origin,omitempty"`
	Destination    *Node          `json:"destination,omitempty"`
	OriginPos      utils.LocalPos `json:"origin_pos"`
	DestinationPos utils.LocalPos `json:"destination_pos"`
	RefPos         utils.LocalPos `json:"ref_pos"`        // only applicable for circular connections
	VarianceAngle  float64        `json:"variance_angle"` // only applicable for circular connections
	Angle          float64        `json:"angle,omitempty"`
	Col            color.NRGBA    `json:"col"`
	Thickness      float32        `json:"thickness,omitempty"`
	Type           ConnectionType `json:"type,omitempty"`
	EstPos         utils.LocalPos `json:"est_pos"`
	EstDim         utils.LocalDim `json:"est_dim"`
	EstPadding     float32        `json:"est_padding,omitempty"`
	EstWidth       float32        `json:"est_width,omitempty"`
	AlongLineProp  float32        `json:"along_line_prop,omitempty"`
	Est            float64        `json:"est,omitempty"`
	PValue         float64        `json:"p_value,omitempty"`
	CI             [2]float64     `json:"ci,omitempty"`
	EstText        string         `json:"est_text,omitempty"`
	Bold           bool           `json:"bold,omitempty"`
	Curvature      float32        `json:"curvature,omitempty"`
	UserDefined    bool           `json:"user_defined,omitempty"`
}

type FontSettings struct {
	Family string        `json:"family,omitempty"`
	Size   float32       `json:"size,omitempty"`
	Face   font.FontFace `json:"-"`
}

type Model struct {
	Nodes         []*Node                  `json:"nodes,omitempty"`
	Connections   []*Connection            `json:"connections,omitempty"`
	Network       map[*Node][]*Node        `json:"-"`
	Font          FontSettings             `json:"font"`
	CoeffDisplay  utils.CoefficientDisplay `json:"coeff_display,omitempty"`
	ViewGenerated bool                     `json:"view_generated,omitempty"`
}
