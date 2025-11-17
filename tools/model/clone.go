package model

func (m *Model) Clone() *Model {
	if m == nil {
		return nil
	}

	// Create node mapping to maintain pointer relationships
	nodeMap := make(map[string]*Node)

	// Deep copy all nodes first
	newNodes := make([]*Node, len(m.Nodes))
	for i, n := range m.Nodes {
		newNode := n.Clone()
		newNodes[i] = newNode
		nodeMap[n.VarName] = newNode
	}

	// Deep copy connections, remapping node pointers
	newConnections := make([]*Connection, len(m.Connections))
	for i, c := range m.Connections {
		newConnections[i] = c.Clone(nodeMap)
	}

	return &Model{
		Nodes:         newNodes,
		Connections:   newConnections,
		Font:          m.Font,
		CoeffDisplay:  m.CoeffDisplay,
		ViewGenerated: m.ViewGenerated,
		PxPerDp:       m.PxPerDp,
	}
}

// deepCopy creates a deep copy of a Node
func (n *Node) Clone() *Node {
	return &Node{
		Class:       n.Class,
		Pos:         n.Pos,
		Dim:         n.Dim,
		Col:         n.Col,
		VarName:     n.VarName,
		Text:        n.Text,
		TextWidth:   n.TextWidth,
		Bold:        n.Bold,
		Thickness:   n.Thickness,
		UserDefined: n.UserDefined,
		Visible:     n.Visible,
		Padding:     n.Padding,
	}
}

// deepCopy creates a deep copy of a Connection, remapping node pointers
func (c *Connection) Clone(nodeMap map[string]*Node) *Connection {
	newConn := &Connection{
		OriginPos:      c.OriginPos,
		DestinationPos: c.DestinationPos,
		RefPos:         c.RefPos,
		VarianceAngle:  c.VarianceAngle,
		Angle:          c.Angle,
		Col:            c.Col,
		Thickness:      c.Thickness,
		Type:           c.Type,
		EstPos:         c.EstPos,
		EstDim:         c.EstDim,
		EstPadding:     c.EstPadding,
		EstWidth:       c.EstWidth,
		AlongLineProp:  c.AlongLineProp,
		Est:            c.Est,
		PValue:         c.PValue,
		CI:             c.CI,
		EstText:        c.EstText,
		Bold:           c.Bold,
		Curvature:      c.Curvature,
		UserDefined:    c.UserDefined,
	}

	// Remap node pointers to the new copies
	newConn.Origin = nodeMap[c.Origin.VarName]
	newConn.Destination = nodeMap[c.Destination.VarName]

	return newConn
}
