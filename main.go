package main

import (
	"image"
	"log"
	"main/model"
	"main/utils"
	"os"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const (
	startingWidth    int = 1200
	startingHeight   int = 800
	roundness            = .3
	editorVertOffset     = 30
)

var (
	leftClickTag  = new(int)
	rightClickTag = new(int)
	fontSize      = float32(16)
	fontFace      font.FontFace
)

// EditContext contains the current editor state
type EditContext struct {
	viewportCenter   utils.LocalPos
	windowSize       utils.GlobalDim
	scaleFactor      float32
	snapGridSize     float32
	nodeDragOffset   utils.LocalPos // The offset of the cursor from the center of a node when clicking
	panClickPos      utils.LocalPos
	panOffset        utils.LocalPos
	selectedNode     *model.Node
	editingSelection interface{}
}

func main() {
	m := model.InitTestModel()
	ec := InitEditContext()
	widgets := InitWidgets(m)
	th := material.NewTheme()

	//fontFace = utils.LoadCousineFontFace()[0] // monospaced font
	fontFace = utils.LoadSansFontFace()[0]

	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Pub SEM"))
		w.Option(app.Size(unit.Dp(startingWidth), unit.Dp(startingHeight)))
		if err := loop(w, th, m, ec, widgets); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, th *material.Theme, m *model.Model, ec *EditContext, widgets ModelWidgets) error {
	ops := new(op.Ops)

	// listen for events in the window.
	for {
		// detect what type of event
		switch e := w.Event().(type) {

		// this is sent when the application should re-render.
		case app.FrameEvent:
			// gtx is used to pass around rendering and event information.
			gtx := app.NewContext(ops, e)

			ec.windowSize = utils.GlobalDim{W: gtx.Constraints.Max.X, H: gtx.Constraints.Max.Y}

			// handle scrolling to zoom
			Scroll(ops, gtx, ec)

			// draw the model
			DrawModel(ops, gtx, m, ec)

			//editorLayout := app.NewContext(ops, e)
			if ec.editingSelection != nil {
				switch s := ec.editingSelection.(type) {
				case *model.Node:
					topNodePos := s.Pos.Sub(utils.LocalPos{Y: s.Dim.H / 2})
					posOffset := topNodePos.Sub(utils.LocalPos{Y: editorVertOffset})
					widgets.DrawNodeEditor(ops, gtx, th, s, posOffset, ec)
				case *model.Connection:
				}
			}

			// if not clicking a node, panning is available
			LeftClick(ops, gtx, m, ec)

			RightClick(ops, gtx, m, ec, widgets)

			// complete the frame event
			e.Frame(gtx.Ops)

		// this is sent when the application is closed
		case app.DestroyEvent:
			return e.Err
		}
	}
}

func LeftClick(ops *op.Ops, gtx layout.Context, m *model.Model, ec *EditContext) {
	// Register for pan events on the entire window
	event.Op(ops, leftClickTag)

	// Process pan events
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: leftClickTag,
			Kinds:  pointer.Press | pointer.Drag | pointer.Release,
		})
		if !ok {
			break
		}

		switch evt := ev.(type) {
		case pointer.Event:
			switch evt.Kind {
			case pointer.Press:
				// Only respond to left mouse button
				if evt.Buttons != pointer.ButtonPrimary {
					continue
				}

				// check if clicking a node
				for _, n := range m.Nodes {
					if ec.selectedNode != nil {
						break
					}

					rect := utils.MakeRect(
						n.Pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
						n.Dim.ToGlobal(ec.scaleFactor),
					)

					switch n.Class {
					case model.OBSERVED:
						if utils.WithinRect(evt.Position.Round(), rect) {
							ec.selectedNode = n
						}
					case model.LATENT:
						if utils.WithinEllipse(evt.Position.Round(), rect) {
							ec.selectedNode = n
						}
					case model.INTERCEPT:
						//todo: handle intercept
					}
				}

				if n := ec.selectedNode; n != nil { // if clicking a node...
					ec.nodeDragOffset = utils.ToLocalPos(evt.Position).Div(ec.scaleFactor).Sub(n.Pos)
				} else { // if not clicking a node, then setup pan
					ec.panClickPos = utils.ToLocalPos(evt.Position)
					ec.panOffset = ec.viewportCenter
				}
			case pointer.Drag:
				// Only respond to left mouse button
				if evt.Buttons != pointer.ButtonPrimary {
					continue
				}
				if n := ec.selectedNode; n != nil { // if dragging a node...
					newPos := utils.ToLocalPos(evt.Position).Div(ec.scaleFactor).Sub(ec.nodeDragOffset)
					n.Pos = utils.SnapToGrid(newPos, ec.snapGridSize)
				} else { // if not dragging a node, then pan
					panDelta := utils.ToLocalPos(evt.Position).Sub(ec.panClickPos).Div(ec.scaleFactor)
					ec.viewportCenter = ec.panOffset.Add(panDelta)
					// add pan cursor hand
					pointer.CursorGrab.Add(ops)
				}
			case pointer.Release:
				ec.selectedNode = nil
				pointer.CursorDefault.Add(ops)
			default:
			}
		}
	}
}

func RightClick(ops *op.Ops, gtx layout.Context, m *model.Model, ec *EditContext, widgets ModelWidgets) {
	// Register for pan events on the entire window
	event.Op(ops, rightClickTag)

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: rightClickTag,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}

		switch evt := ev.(type) {
		case pointer.Event:
			switch evt.Kind {
			case pointer.Press:
				// Only respond to right mouse button
				if evt.Buttons != pointer.ButtonSecondary {
					continue
				}

				for _, c := range m.Connections {
					tolerance := float32(5)
					samples := 10
					if WithinConnection(evt.Position.Round(), c, ec, tolerance, samples) {
						if ec.editingSelection != c {
							ec.editingSelection = c
						} else {
							ec.editingSelection = nil
						}
					}
				}

				for _, n := range m.Nodes {
					rect := utils.MakeRect(
						n.Pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
						n.Dim.ToGlobal(ec.scaleFactor),
					)

					switch n.Class {
					case model.OBSERVED:
						if utils.WithinRect(evt.Position.Round(), rect) {
							if ec.editingSelection != n {
								ec.editingSelection = n
							} else {
								ec.editingSelection = nil
							}
						}
					case model.LATENT:
						if utils.WithinEllipse(evt.Position.Round(), rect) {
							if ec.editingSelection != n {
								ec.editingSelection = n
							} else {
								ec.editingSelection = nil
							}
						}
					}
				}

			default:
			}
		}
	}
}

func Scroll(ops *op.Ops, gtx layout.Context, ec *EditContext) {
	// Register for scroll events on the entire window
	event.Op(ops, ops)

	// Process scroll events
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:  ops,
			Kinds:   pointer.Scroll,
			ScrollY: pointer.ScrollRange{Min: -100, Max: 100},
		})
		if !ok {
			break
		}

		switch evt := ev.(type) {
		case pointer.Event:
			if evt.Kind == pointer.Scroll {
				// Adjust scale factor based on scroll direction
				var zoomSpeed float32 = 0.01
				ec.scaleFactor -= evt.Scroll.Y * zoomSpeed

				// Clamp scale factor to reasonable bounds
				if ec.scaleFactor < 0.5 {
					ec.scaleFactor = 0.5
				} else if ec.scaleFactor > 5.0 {
					ec.scaleFactor = 5.0
				}
			}
		}
	}
}

func InitEditContext() *EditContext {
	ec := new(EditContext)
	ec.scaleFactor = 1.0
	ec.snapGridSize = 20.0
	return ec
}

func WithinConnection(pos image.Point, c *model.Connection, ec *EditContext, tolerance float32, samples int) bool {
	posA := c.OriginPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize)
	posB := c.DestinationPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize)

	// Add thickness/2 to tolerance for better hit detection
	hitRadius := tolerance + (c.Thickness*ec.scaleFactor)/2

	if c.Type == model.COVARIANCE {
		return utils.WithinArc(pos, posA, posB, roundness, hitRadius, c.Curvature, samples)
	}
	return utils.WithinLine(pos, posA, posB, hitRadius)
}
