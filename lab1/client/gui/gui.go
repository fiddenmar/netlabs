package main

import (
	client "../"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	//"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
	//"time"
)

var cl client.Client
var items []string

func appMain(driver gxui.Driver) {
	items := []string{"shit", "crap"}
	cl.Init("User", "127.0.0.1", 34310)

	theme := flags.CreateTheme(driver)

	label := theme.CreateLabel()
	label.SetText("UDP чат")

	adapter := gxui.CreateDefaultAdapter()
	adapter.SetItems(items)

	chatbox := theme.CreateList()
	chatbox.SetAdapter(adapter)

	messageBox := theme.CreateTextBox()

	buttonSend := theme.CreateButton()
	buttonSend.SetText("Send")

	layout := theme.CreateLinearLayout()
	layout.AddChild(label)
	layout.AddChild(chatbox)
	layout.AddChild(messageBox)
	layout.AddChild(buttonSend)
	layout.SetHorizontalAlignment(gxui.AlignLeft)

	window := theme.CreateWindow(800, 600, "Progress bar")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(layout)
	window.OnClose(driver.Terminate)
	go cl.Answer()

	items[1] := "sad"

}

func main() {

	gl.StartDriver(appMain)

}
