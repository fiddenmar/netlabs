package main

import (
	client "../"
	"github.com/conformal/gotk3/gtk"
	"log"
	"os"
	"strconv"
)

func printMessages(msgView *gtk.TextView, msgs chan string) {
	for {
		msg:=<-msgs
		buffer, err := msgView.GetBuffer()
		if err != nil {
	        log.Fatal("Unable to load buffer:", err)
	    }
		start, end := buffer.GetBounds()
		text, err := buffer.GetText(start, end, true)
		if err != nil {
	        log.Fatal("Unable to save buffer as string:", err)
	    }
		buffer.SetText(text+"\n"+msg)
	}
}

func main() {
    gtk.Init(nil)

    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window:", err)
    }
    win.SetTitle("UDP client v0.0.1")
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })

    var c client.Client
    answers := make(chan string)

    grid, err := gtk.GridNew()
    if err != nil {
    	log.Fatal("Unable to create grid:", err)
    }

    messageHistory, err := gtk.TextViewNew()
    if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	grid.Attach(messageHistory, 0, 0, 4, 1)

    messageEntry, err:= gtk.EntryNew()
    if err != nil {
    	log.Fatal("Unable to create entry:", err)
    }
    grid.Attach(messageEntry, 0, 1, 1, 1)

    privateEntry, err:=gtk.EntryNew()
    if err != nil {
    	log.Fatal("Unable to create entry:", err)
    }
    grid.Attach(privateEntry, 1, 1, 1, 1)

    sendButton, err := gtk.ButtonNewWithLabel("Send")
    if err != nil {
        log.Fatal("Unable to create button:", err)
    }
    sendButton.Connect("clicked", func(btn *gtk.Button){
    	lbl, _ := btn.GetLabel()
    	if lbl!="Send" {
    		return
    	}
    	log.Print(lbl)
    	msg, _ := messageEntry.GetText()
    	log.Print(msg)
    	c.Message(msg)
    	})
    grid.Attach(sendButton, 0, 2, 1, 1)

    privateButton, err := gtk.ButtonNewWithLabel("Private")
    if err != nil {
        log.Fatal("Unable to create button:", err)
    }
    privateButton.Connect("clicked", func(btn *gtk.Button){
    	lbl, _ := btn.GetLabel()
    	if lbl!="Private" {
    		return
    	}
    	log.Print(lbl)
    	private, _ := privateEntry.GetText()
    	log.Print(private)
    	if private!="" {
    		msg, _ := messageEntry.GetText()
    		log.Print(msg)
    		c.Private(private, msg)
    	}
    	})
    grid.Attach(privateButton, 1, 2, 1, 1)

    listButton, err := gtk.ButtonNewWithLabel("List")
    if err != nil {
        log.Fatal("Unable to create button:", err)
    }
    listButton.Connect("clicked", func(btn *gtk.Button){
    	lbl, _ := btn.GetLabel()
    	if lbl!="List" {
    		return
    	}
    	log.Print(lbl)
    	c.List()
    	log.Print(lbl)
    	})
    grid.Attach(listButton, 2, 1, 1, 1)

    leaveButton, err := gtk.ButtonNewWithLabel("Leave")
    if err != nil {
        log.Fatal("Unable to create button:", err)
    }
    leaveButton.Connect("clicked", func(btn *gtk.Button){
    	lbl, _ := btn.GetLabel()
    	if lbl!="Leave" {
    		return
    	}
    	log.Print(lbl)
    	c.Leave()
	    os.Exit(0)
    	})
    grid.Attach(leaveButton, 3, 1, 1, 1)
    win.Add(grid)
    // Set the default window size.
    win.SetDefaultSize(400, 600)

    // Recursively show all widgets contained in this window.
    win.ShowAll()

	port,_:=strconv.Atoi(os.Args[3])
	c.Init(os.Args[1], os.Args[2], port, answers)
	go printMessages(messageHistory, answers)
	go c.Answer()
	c.Register()
	gtk.Main()
}
