package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type mainMenu struct {
	*fyne.MainMenu

	sysMenu *fyne.Menu
	quit    *fyne.MenuItem
}

type systemData struct {
	folderURIs   []fyne.URI
	fullscreen   bool
	uriIndex     int
	windowWidth  float64
	windowHeight float64
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Vishal's Photos")
	//myWindow.SetFixedSize(true)
	//myWindow.SetFullScreen(true)
	sysData := systemData{}
	sysData.fullscreen = false
	sysData.windowWidth = 1280.0
	sysData.windowHeight = 720.0

	var mm mainMenu
	mm.quit = fyne.NewMenuItem("Quit", myWindow.Close)
	var fileOpenDialog = fyne.NewMenuItem("Open File", func() {
		var file_dialog = dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			//image := canvas.NewImageFromResource(theme.FyneLogo())
			image := canvas.NewImageFromURI(uc.URI())
			// image := canvas.NewImageFromImage(src)
			// image := canvas.NewImageFromReader(reader, name)
			// image := canvas.NewImageFromFile(fileName)
			image.FillMode = canvas.ImageFillContain
			//myWindow.Resize(fyne.NewSize(image.MinSize().Width, image.MinSize().Height))
			myWindow.SetContent(image)
		}, myWindow)

		file_dialog.SetFilter(
			storage.NewExtensionFileFilter([]string{".jpg", ".png"}))
		file_dialog.Show()
	})
	var folderOpenDialog = fyne.NewMenuItem("Open Folder", func() {
		var folder_dialog = dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
			//sysData.folderURIs, _ = lu.List()
			dat, _ := lu.List()
			var dtx []fyne.URI
			for i := 0; i < len(dat); i++ {
				if contains([]string{".jpg", ".png"}, dat[i].Extension()) {
					dtx = append(dtx, dat[i])
				}
			}
			sysData.folderURIs = dtx
			folderImageLoad(myWindow, sysData, 0)
		}, myWindow)
		folder_dialog.Show()
	})
	var fileMenu = fyne.NewMenu("File", fileOpenDialog, folderOpenDialog, fyne.NewMenuItemSeparator(), mm.quit)

	// about button
	content := widget.NewCard("Vishal's Photos", "created by github@VishalRamki", widget.NewRichTextFromMarkdown(`Vishal's Photos is a simple way to view all your photos in a folder.
	It is built on Fyne.
	The source code is available @ https://github.com/VishalRamki/vishals-photos
	https://vishalramkissoon.com`))
	aboutDialog := dialog.NewCustom("", "Close", content, myWindow)
	aboutdialogMenu := fyne.NewMenuItem("About", func() {
		aboutDialog.Show()
	})
	var aboutMenu = fyne.NewMenu("About", aboutdialogMenu)

	// settings menu
	var fullscreenOptionItemMenu = fyne.NewMenuItem("Fullscreen Mode", func() {
		sysData.fullscreen = !sysData.fullscreen
		myWindow.SetFullScreen(sysData.fullscreen)

		if len(sysData.folderURIs) > 0 {
			folderImageLoad(myWindow, sysData, sysData.uriIndex)
		}
		// reset to 720p
		if !sysData.fullscreen {

			myWindow.Resize(fyne.NewSize(float32(sysData.windowWidth), float32(sysData.windowHeight)))
		}
	})
	var settingsMenu = fyne.NewMenu("Setting", fullscreenOptionItemMenu)

	var mainmm = fyne.NewMainMenu(fileMenu, settingsMenu, aboutMenu)
	myWindow.Resize(fyne.NewSize(1280, 720))
	myWindow.SetMainMenu(mainmm)

	// build layout

	text4 := canvas.NewText("Use File > Open File, to open a single file, and File > Open Folder, to open all the files in the folder.", color.Black)
	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())

	myWindow.SetContent(centered)
	myWindow.ShowAndRun()
}

func folderImageLoad(mainWindow fyne.Window, d systemData, index int) {
	d.uriIndex = index
	if index < 0 || index >= len(d.folderURIs) {
		return
	}

	image := canvas.NewImageFromURI(d.folderURIs[index])
	image.FillMode = canvas.ImageFillContain
	renderWidth := d.windowWidth
	renderHeight := d.windowHeight

	if mainWindow.Canvas().Size().Width > float32(d.windowWidth) {
		renderWidth = float64(mainWindow.Canvas().Size().Width - 10.0)
	}
	if mainWindow.Canvas().Size().Height > float32(d.windowHeight) {
		renderHeight = float64(mainWindow.Canvas().Size().Height - 80.0)
	}
	image.SetMinSize(fyne.NewSize(float32(renderWidth), float32(renderHeight)))
	image.ScaleMode = canvas.ImageScaleFastest
	prevImage := widget.NewButton("<< Previous", func() {
		folderImageLoad(mainWindow, d, index-1)
	})
	nextImage := widget.NewButton("Next >>", func() {
		folderImageLoad(mainWindow, d, index+1)
	})
	fileName := canvas.NewText(d.folderURIs[index].Name(), color.Black)
	content := container.New(layout.NewHBoxLayout(), prevImage, layout.NewSpacer(), fileName, layout.NewSpacer(), nextImage)
	centered := container.New(layout.NewHBoxLayout(), image)
	//grid := container.New(layout.NewGridLayout(1), centered, content)
	vbox := container.New(layout.NewVBoxLayout(), centered, content)
	mainWindow.SetContent(vbox)
}

// utils
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
