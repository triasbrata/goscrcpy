package gui

import (
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ebitick"
	"github.com/triasbrata/gadb"
	"github.com/triasbrata/goscrcpy/internals/slices"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type WindowSettings struct {
	ui            *ebitenui.UI
	listDevices   *widget.List
	timeTick      *ebitick.TimerSystem
	connectDevice func(devices *gadb.Device)
}

func (g *WindowSettings) SetConnectDevice(cb func(device *gadb.Device)) {
	g.connectDevice = cb
}

// Update implements ebiten.Game.
func (g *WindowSettings) UpdateSerialList(devices []*gadb.Device) error {
	existings := make(map[string]*gadb.Device)
	entries, safe := g.listDevices.Entries().([]*gadb.Device)
	if safe {
		existings = slices.Entries(entries, func(d *gadb.Device) (string, *gadb.Device) {
			return d.Serial(), d
		})
	}

	for _, d := range devices {
		if _, safe := existings[d.Serial()]; !safe {
			g.listDevices.AddEntry(d)
			// g.listDevices.SelectFocused()
		}
	}
	return nil

}
func (g *WindowSettings) Cron(cb func()) {
	g.timeTick.After(5*time.Second, cb)
}

func (g *WindowSettings) Update() error {
	g.ui.Update()
	g.timeTick.Update()

	if list, ok := g.ui.GetFocusedWidget().(*widget.List); ok {
		//Test that you can call Click on the focused widget.
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			list.FocusPrevious()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			list.FocusNext()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			list.SelectFocused()
		}
	}

	return nil
}

func (g *WindowSettings) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *WindowSettings) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func NewWindow() ebiten.Game {

	window := &WindowSettings{
		timeTick: ebitick.NewTimerSystem(),
	}
	window.initWindowButton()
	return window

}
func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     150,
		Hinting: font.HintingFull,
	}), nil
}

func (window *WindowSettings) initWindowButton() {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(12)
	// Create array of list entries

	entries := make([]any, 0)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Construct a list. This is one of the more complicated widgets to use since
	// it is composed of multiple widget types
	window.listDevices = widget.NewList(
		// Set how wide the list should be
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
				Padding:            widget.NewInsetsSimple(50),
			}),
		)),
		// Set the entries in the list
		widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerOpts(
			// Set the background images/color for the list
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		widget.ListOpts.SliderOpts(
			// Set the background images/color for the background of the slider track
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}, buttonImage),
			widget.SliderOpts.MinHandleSize(5),
			// Set how wide the track should be
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
		// Hide the horizontal slider
		widget.ListOpts.HideHorizontalSlider(),
		// Set the font for the list options
		widget.ListOpts.EntryFontFace(face),
		// Set the colors for the list
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{R: 0, G: 255, B: 0, A: 255},     // Foreground color for the unfocused selected entry
			Unselected:                 color.NRGBA{R: 254, G: 255, B: 255, A: 255}, // Foreground color for the unfocused unselected entry
			SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255}, // Background color for the unfocused selected entry
			SelectingBackground:        color.NRGBA{R: 130, G: 130, B: 130, A: 255}, // Background color for the unfocused being selected entry
			SelectingFocusedBackground: color.NRGBA{R: 130, G: 140, B: 170, A: 255}, // Background color for the focused being selected entry
			SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255}, // Background color for the focused selected entry
			FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255}, // Background color for the focused unselected entry
			DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Foreground color for the disabled unselected entry
			DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Foreground color for the disabled selected entry
			DisabledSelectedBackground: color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Background color for the disabled selected entry
		}),
		// This required function returns the string displayed in the list
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(*gadb.Device).Serial()
		}),
		// Padding for each entry
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		// Text position for each entry
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		// This handler defines what function to run when a list item is selected.
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(*gadb.Device)
			window.connectDevice(entry)
		}),
		// This option will select the entry as it is focused
		// widget.ListOpts.SelectFocus(),

		// This option will disable default keys (up and down)
		//widget.ListOpts.DisableDefaultKeys(true),
	)

	// Add list to the root container
	rootContainer.AddChild(window.listDevices)
	window.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}
