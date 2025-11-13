package read_write

import (
	"encoding/json"
	"image/color"
	"log"
	"main/model"
	"main/utils"
	"os"
)

const (
	roundness     = .3
	estPadding    = 2.5
	propAlongLine = .5
)

type DataRow struct {
	Lhs     string  `json:"lhs"`
	Op      string  `json:"op"`
	Rhs     string  `json:"rhs"`
	User    int     `json:"user"`
	Group   int     `json:"group"`
	Est     float64 `json:"est"`
	Label   string  `json:"label"`
	PValue  float64 `json:"pvalue"`
	CiLower float64 `json:"ci_lower"`
	CiUpper float64 `json:"ci_upper"`
}

func ModelFromJSON(dir string) *model.Model {
	// TODO: Check for existing layout

	rows := readJSON(dir)

	// translate data to model type
	varMap := make(map[string]*model.Node)
	connections := make([]*model.Connection, 0, len(rows))
	var i int
	for _, row := range rows {
		lhs, ok := varMap[row.Lhs]
		if !ok {
			lhs = new(model.Node)
			lhs.Pos = utils.LocalPos{X: float32(i * 100)}
			i++
		}
		rhs, ok := varMap[row.Rhs]
		if !ok {
			rhs = new(model.Node)
			rhs.Pos = utils.LocalPos{X: float32(i * 100)}
			i++
		}

		c := new(model.Connection)

		// set var names (init to be same as text for now)
		lhs.VarName = row.Lhs
		rhs.VarName = row.Rhs
		lhs.Text = row.Lhs
		rhs.Text = row.Rhs

		// Set estimate values
		c.Est = row.Est
		c.PValue = row.PValue
		c.CI = [2]float64{row.CiLower, row.CiUpper}

		// define connection and node types
		skip := false
		switch row.Op {
		case "=~":
			lhs.Class = model.LATENT
			c.Type = model.STRAIGHT
			c.Origin = lhs
			c.Destination = rhs
		case "~~":
			c.Type = model.CURVED
			c.Origin = lhs
			c.Destination = rhs
		case "~":
			c.Type = model.STRAIGHT
			c.Origin = lhs
			c.Destination = rhs
		default:
			skip = true
		}

		if row.User == 1 {
			c.UserDefined = true
		}

		if lhs.Class != model.LATENT {
			lhs.Class = model.OBSERVED
		}

		// set thickness
		lhs.Thickness = 3.0
		rhs.Thickness = 3.0
		c.Thickness = 2.0

		lhs.Padding = estPadding
		rhs.Padding = estPadding

		c.Curvature = roundness
		c.AlongLineProp = propAlongLine
		c.Col = color.NRGBA{A: 255}

		//assign nodes to map
		varMap[row.Lhs] = lhs
		varMap[row.Rhs] = rhs
		// assign connection to array
		if !skip {
			connections = append(connections, c)
		}
	}

	m := new(model.Model)
	m.Font = model.FontSettings{
		Family: "sans",
		Size:   16,
		Face:   utils.LoadSansFontFace()[0],
	}
	m.CoeffDisplay = utils.STAR
	m.Connections = connections
	m.Nodes = utils.MapValsToSlice(varMap)

	return m
}

func readJSON(dir string) []DataRow {
	// read the json file
	var rows []DataRow

	data, err := os.ReadFile(dir + "/temp.json")
	//data, err := os.ReadFile("test.json") // testing
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &rows)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}
