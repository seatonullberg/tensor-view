package main

import (
	"strconv"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

func main() {

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create a blue sphere and add it to the scene
	geom := geometry.NewSphere(1, 32, 16)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// Testing my stuff
	var editFields []*gui.Edit
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			var s string
			if i == j {
				s = "1"
			} else {
				s = "0"
			}
			ed := gui.NewEdit(40, "")
			ed.SetText(s)
			ed.SetFontSize(16)
			ed.SetPosition(float32(i*50+5), float32(j*30+5))
			editFields = append(editFields, ed)
			scene.Add(ed)
		}
	}
	btn := gui.NewButton("Update")
	btn.SetPosition(47.5, 100)
	btn.SetSize(50, 20)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		// fmt.Println(editFields[0].Text())
		m := func() *math32.Matrix3 {
			var entries []float32
			for _, ed := range editFields {
				f, err := strconv.ParseFloat(ed.Text(), 32)
				if err != nil {
					panic("bad entry")
				}
				entries = append(entries, float32(f))
			}
			m := math32.NewMatrix3()
			m.FromArray(entries, 0)
			return m
		}
		mesh.Dispose()                       // dispose of existing mesh
		geom = geometry.NewSphere(1, 32, 16) // reset from a sphere
		mesh = graphic.NewMesh(geom, mat)
		scene.Add(mesh)
		// do the transformation
		geom.OperateOnVertices(func(vertex *math32.Vector3) bool {
			vertex.ApplyMatrix3(m())
			return false
		})
	})
	scene.Add(btn)

	// Create and add lights to the scene
	//scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	//pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	scene.Add(light.NewAmbient(math32.NewColor("white"), 0.8))
	pointLight := light.NewPoint(math32.NewColor("white"), 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(10))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
