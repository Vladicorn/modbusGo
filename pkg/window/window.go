// SPDX-License-Identifier: Unlicense OR MIT

package window

// A Gio program that demonstrates Gio widgets. See https://gioui.org for more information.

import (
	"Modbus/pkg/modbus"
	"context"
	"fmt"
	"gioui.org/font/gofont"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

type iconAndTextButton struct {
	theme  *material.Theme
	button *widget.Clickable
	icon   *widget.Icon
	word   string
}

var cancel = make(chan bool)
var OutModbus = make(chan uint16)

func CreateWindow() {

	editor.SetText(longText)
	ic, err := widget.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}
	icon = ic

	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())

	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				if checkbox.Changed() {
					if checkbox.Value {
						//	fmt.Println("Hi")
						transformTime = e.Now
					} else {
						transformTime = time.Time{}
					}
				}
				kitchen(gtx, th)
				e.Frame(gtx.Ops)
			}

		}
	}
}

var (
	ctxPar         = context.Background()
	ctx, cancelCtx = context.WithCancel(ctxPar)
)

var (
	editor     = new(widget.Editor)
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	lineEditorIp = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	lineEditorPort = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	lineEditorSlaveId = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	lineEditorAdr = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	lineEditorQuantity = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}

	button            = new(widget.Clickable)
	buttonStart       = new(widget.Clickable)
	buttonFinish      = new(widget.Clickable)
	greenButton       = new(widget.Clickable)
	iconTextButton    = new(widget.Clickable)
	iconButton        = new(widget.Clickable)
	flatBtn           = new(widget.Clickable)
	disableBtn        = new(widget.Clickable)
	radioButtonsGroup = new(widget.Enum)
	list              = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	progress            = float32(0)
	progressIncrementer chan float32
	green               = true
	topLabel            = "Hello, Gio"
	topLabelState       = new(widget.Selectable)
	icon                *widget.Icon
	checkbox            = new(widget.Bool)
	swtch               = new(widget.Bool)
	transformTime       time.Time
	float               = new(widget.Float)
	disableBut          bool
)

type (
	D = layout.Dimensions
	C = layout.Context
)

func (b iconAndTextButton) Layout(gtx layout.Context) layout.Dimensions {
	return material.ButtonLayout(b.theme, b.button).Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)
			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					var d D
					if b.icon != nil {
						size := gtx.Dp(unit.Dp(56)) - 2*gtx.Dp(unit.Dp(16))
						gtx.Constraints = layout.Exact(image.Pt(size, size))
						d = b.icon.Layout(gtx, b.theme.ContrastFg)
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					l := material.Body1(b.theme, b.word)
					l.Color = b.theme.Palette.ContrastFg
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}

//Прорисовка
func kitchen(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var conTest modbus.ConnectionModbus

	for _, e := range lineEditor.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			topLabel = e.Text
			lineEditor.SetText("")
		}
	}
	widgets := []layout.Widget{
		//Заглавная
		func(gtx C) D {
			l := material.H5(th, "Apeyron Modbus")
			l.State = topLabelState
			return l.Layout(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(200))
			return material.Editor(th, editor, "Hint").Layout(gtx)
		},
		//Строка ввода
		func(gtx C) D {
			e := material.Editor(th, lineEditorIp, "Введите ваш IP")

			e.Font.Style = text.Italic
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},

		func(gtx C) D {
			e := material.Editor(th, lineEditorPort, "Введите ваш порт")

			e.Font.Style = text.Italic
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},

		func(gtx C) D {
			e := material.Editor(th, lineEditorAdr, "Введите адрес")

			e.Font.Style = text.Italic
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},

		func(gtx C) D {
			e := material.Editor(th, lineEditorQuantity, "Введите размер чтения")

			e.Font.Style = text.Italic
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},

		func(gtx C) D {
			e := material.Editor(th, lineEditorSlaveId, "Введите slaveId")

			e.Font.Style = text.Italic
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},

		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(8))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						for buttonStart.Clicked() {
							if conTest.IpAdr != "" {
								disableBut = !disableBut
								go conTest.Connect(cancel, OutModbus)
							}

						}
						if disableBut {
							gtx = gtx.Disabled()
						}
						return material.Button(th, buttonStart, "Начать").Layout(gtx)
					})
				}),
			)
		},

		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(8))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						for buttonFinish.Clicked() {
							cancel <- true
							disableBut = !disableBut
						}
						return material.Button(th, buttonFinish, "Стоп").Layout(gtx)
					})
				}),
			)
		},
	}
	select {
	case result := <-OutModbus:
		test := []layout.Widget{
			func(gtx C) D {
				resText := fmt.Sprintf("%v", result)
				l := material.H5(th, resText)
				l.State = topLabelState
				return l.Layout(gtx)
			}}
		widgets = append(widgets, test[0])
	default:
	}

	conTest.Port = lineEditorPort.Text()
	conTest.IpAdr = lineEditorIp.Text()
	conTest.SlaveId = lineEditorSlaveId.Text()
	conTest.Quantity = lineEditorQuantity.Text()
	conTest.Adr = lineEditorAdr.Text()

	return material.List(th, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})

}

const longText = `Программа для работы с ModbusTCP устройствами`
