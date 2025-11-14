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
	EdgeConnections [4][]*Connection `json:"-"` // only applicable for rectangular nodes
	Padding         float32          `json:"padding,omitempty"`
}

type Connection struct {
	Origin         *Node          `json:"origin,omitempty" :"origin" :"origin"`
	Destination    *Node          `json:"destination,omitempty" :"destination" :"destination"`
	OriginPos      utils.LocalPos `json:"origin_pos" :"origin_pos" :"origin_pos"`
	DestinationPos utils.LocalPos `json:"destination_pos" :"destination_pos" :"destination_pos"`
	Angle          float64        `json:"angle,omitempty" :"angle" :"angle"`
	Col            color.NRGBA    `json:"col" :"col" :"col"`
	Thickness      float32        `json:"thickness,omitempty" :"thickness" :"thickness"`
	Type           ConnectionType `json:"type,omitempty" :"type" :"type"`
	EstPos         utils.LocalPos `json:"est_pos" :"est_pos" :"est_pos"`
	EstDim         utils.LocalDim `json:"est_dim" :"est_dim" :"est_dim"`
	EstPadding     float32        `json:"est_padding,omitempty" :"est_padding" :"est_padding"`
	EstWidth       float32        `json:"est_width,omitempty" :"est_width" :"est_width"`
	AlongLineProp  float32        `json:"along_line_prop,omitempty" :"along_line_prop" :"along_line_prop"` // how far along the line is the estimate label? defaults to .5
	Est            float64        `json:"est,omitempty" :"est" :"est"`
	PValue         float64        `json:"p_value,omitempty" :"p_value" :"p_value"`
	CI             [2]float64     `json:"ci,omitempty" :"ci" :"ci"`
	EstText        string         `json:"est_text,omitempty" :"est_text" :"est_text"`
	Bold           bool           `json:"bold,omitempty" :"bold" :"bold"`
	Curvature      float32        `json:"curvature,omitempty" :"curvature" :"curvature"` // only applicable for covariance (curved) connections
	UserDefined    bool           `json:"user_defined,omitempty" :"user_defined" :"user_defined"`
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
