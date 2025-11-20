package read_write

import (
	"encoding/json"
	"image/color"
	"log"
	"main/model"
	"main/utils"
	"math"
	"math/rand"
	"os"
	"path/filepath"
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

func ModelFromJSON(dir, projectName string) *model.Model {
	m := new(model.Model)
	// Check for existing project
	projPath := filepath.Join(dir, projectName+".json")
	loadedProj := false
	mExisting, err := LoadProject(projPath)
	if err == nil {
		m = mExisting
		loadedProj = true
	}

	tempPath := filepath.Join(dir, "temp.json")
	rows := readJSON(tempPath)

	// translate data to model type
	varMap := make(map[string]*model.Node)

	// account for existing nodes from project
	for _, n := range m.Nodes {
		n.Visible = false
		varMap[n.VarName] = n
	}

	connections := make([]*model.Connection, 0)
	randMag := float32(2000)
	var i int
	for _, row := range rows {
		if !(row.Op == "=~" || row.Op == "~~" || row.Op == "~") {
			continue
		}

		lhs, ok := varMap[row.Lhs]
		if !ok {
			lhs = new(model.Node)
			pos := utils.LocalPos{X: (rand.Float32() - .5) * randMag, Y: (rand.Float32() - .5) * randMag}
			lhs.Pos = utils.SnapToGrid(pos, 20)
			lhs.Col = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			i++
		}
		rhs, ok := varMap[row.Rhs]
		if !ok {
			rhs = new(model.Node)
			pos := utils.LocalPos{X: (rand.Float32() - .5) * randMag, Y: (rand.Float32() - .5) * randMag}
			rhs.Pos = utils.SnapToGrid(pos, 20)
			rhs.Col = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			i++
		}

		c := new(model.Connection)

		lhs.Visible = true
		rhs.Visible = true

		// set var names (init to be same as text for now)
		lhs.VarName = row.Lhs
		rhs.VarName = row.Rhs

		if lhs.Text == "" {
			lhs.Text = row.Lhs
		}

		if rhs.Text == "" {
			rhs.Text = row.Rhs
		}

		// Set estimate values
		c.Est = row.Est
		c.PValue = row.PValue
		c.CI = [2]float64{row.CiLower, row.CiUpper}

		// define connection and node types
		switch row.Op {
		case "=~":
			lhs.Class = model.LATENT
			c.Type = model.STRAIGHT
			c.Origin = lhs
			c.Destination = rhs
		case "~~":
			if lhs.VarName == rhs.VarName {
				c.Type = model.CIRCULAR
			} else {
				c.Type = model.CURVED
			}
			c.Origin = lhs
			c.Destination = rhs
		case "~":
			c.Type = model.STRAIGHT
			c.Origin = rhs
			c.Destination = lhs
		default:
			continue
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

		c.EstPadding = estPadding

		c.Curvature = roundness
		c.AlongLineProp = propAlongLine
		c.Col = color.NRGBA{A: 255}

		//assign nodes to map
		varMap[row.Lhs] = lhs
		varMap[row.Rhs] = rhs

		// check of connection already exists. Match label placements
		if mExisting != nil {
			for _, cExisting := range mExisting.Connections {
				if c.Origin.VarName == cExisting.Origin.VarName && c.Destination.VarName == cExisting.Destination.VarName && c.Type == cExisting.Type {
					c.AlongLineProp = cExisting.AlongLineProp
					c.VarianceAngle = cExisting.VarianceAngle
					c.Curvature = cExisting.Curvature
				}
			}
		}

		connections = append(connections, c)

	}

	if mExisting != nil {
		m.CoeffDisplay = mExisting.CoeffDisplay
		m.Font = mExisting.Font
		m.PxPerDp = mExisting.PxPerDp
	} else {
		m.Font = model.FontSettings{
			Family: "sans",
			Size:   16,
			Faces:  utils.LoadSansFontFace(),
		}
		m.CoeffDisplay = utils.STAR
	}

	m.Connections = connections
	m.Nodes = utils.MapValsToSlice(varMap)
	m.Network = CalculateNodeNetwork(connections)

	if !loadedProj {
		forceDirectNodes(m)
	}

	return m
}

func CalculateNodeNetwork(connections []*model.Connection) map[*model.Node][]*model.Node {
	res := make(map[*model.Node][]*model.Node)
	for _, c := range connections {
		if res[c.Origin] == nil {
			res[c.Origin] = make([]*model.Node, 0)
			res[c.Origin] = append(res[c.Origin], c.Origin)
		}
		if res[c.Destination] == nil {
			res[c.Destination] = make([]*model.Node, 0)
			res[c.Destination] = append(res[c.Destination], c.Destination)
		}

		if c.Destination.Visible {
			res[c.Origin] = append(res[c.Origin], c.Destination)
		}
		if c.Origin.Visible {
			res[c.Destination] = append(res[c.Destination], c.Origin)
		}
	}
	return res
}

func forceDirectNodes(m *model.Model) {
	// see https://i11www.iti.kit.edu/_media/teaching/winter2016/graphvis/graphvis-ws16-v6.pdf
	var idealSpringLength float64 = 10
	var repelForce float32 = .8
	var attractForce float32 = .7
	var temperature float32 = 150.0 // Initial max movement, "cools" over time
	var coolingFactor float32 = 0.99

	funcRepel := func(a, b utils.LocalPos) utils.LocalPos {
		dist := utils.DistLoc(a, b)
		vec := utils.UnitVector(a, b)
		return vec.Mul(-repelForce / (dist * dist))
	}

	funAttract := func(a, b utils.LocalPos) utils.LocalPos {
		dist := utils.DistLoc(a, b)
		vec := utils.UnitVector(a, b)
		return vec.Mul(attractForce * float32(math.Log(float64(dist)/idealSpringLength)))
	}

	// sore the net displacement for each node per iteration
	displacements := make(map[*model.Node]utils.LocalPos)
	for i := 0; i < 10; i++ {
		// Every node repels every other node
		for j := 0; j < len(m.Nodes); j++ {
			for k := j + 1; k < len(m.Nodes); k++ {
				u := m.Nodes[j]
				v := m.Nodes[k]

				// Calculate the repulsive force vector
				repelVec := funcRepel(u.Pos, v.Pos)

				// Apply force to both nodes (Newton's 3rd Law)
				displacements[u] = displacements[u].Add(repelVec)
				displacements[v] = displacements[v].Sub(repelVec) // Equal and opposite
			}
		}

		// Calculate Attractive Forces (Edges-Only)
		for _, conn := range m.Connections {
			u := conn.Origin
			v := conn.Destination

			if u == nil || v == nil {
				continue // Skip invalid connections
			}

			if u == v {
				continue
			}

			// Calculate the attractive/repulsive spring force vector
			attractVec := funAttract(u.Pos, v.Pos)

			// Apply force to both nodes
			displacements[u] = displacements[u].Add(attractVec)
			displacements[v] = displacements[v].Sub(attractVec) // Equal and opposite
		}

		// Update Node Positions
		// Apply the calculated displacements, capped by the current temperature
		for _, node := range m.Nodes {
			disp := displacements[node]
			dispMag := utils.DistLoc(displacements[node], utils.LocalPos{0, 0})

			// Don't let the node move further than the current temperature
			if dispMag > temperature {
				// Scale the displacement vector down to match the temperature
				disp = disp.Div(dispMag).Mul(temperature)
			}

			// Apply the final displacement to the node's position
			node.Pos = node.Pos.Add(disp)
		}

		// 5. Cool the temperature
		// This makes the layout stabilize over time
		temperature *= coolingFactor
	}
}

func readJSON(path string) []DataRow {
	// read the json file
	var rows []DataRow

	data, err := os.ReadFile(path)
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
