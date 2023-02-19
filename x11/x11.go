package x11

import (
	"fmt"
	"strings"
	"time"

	x11 "github.com/linuxdeepin/go-x11-client"
	ewmh "github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
)

// TODO: Our system doesn't work when the windows have the same name, which is a
// clear issue. We need a way to distinguish windows better than the name.
// /////////////////////////////////////////////////////////////////////////////

// TODO: Is what has to be passed to `SetActiveWindow` so we will need to
// save this type of data in some way or another `xproto.Window`

type X11 struct {
	Client  *x11.Conn
	Windows []Window

	// TODO: Maybe just cache the active window name so we do simple name
	// comparison, but this leads to a bug where two windows with the same name
	// are considered the name window
	ActiveWindowTitle     string
	ActiveWindowChangedAt time.Time
}

func ConnectToX11() *x11.Conn {
	client, err := x11.NewConn()
	if err != nil {
		panic(err)
	}
	return client
}

// TODO: If we can move some of these to be methods of Window struct, it would
// be better organized but there will be obvious limitations we have to work
// through
func (x *X11) HasActiveWindowChanged() bool {
	return x.ActiveWindowTitle != x.ActiveWindow().Title
}

// TODO: WE could manage the windows and switch between without Alt+Tab which
// would be far better

func (x *X11) ActiveWindow() Window {
	// TODO: This returns an x.Window object which can get all sorts of
	// information beyond just the name, like the PID. We shouldn't need a second
	// call at all to get the title of the window, thats absurdist.
	activeWindow, err := ewmh.GetActiveWindow(x.Client).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetActiveWindow(x.Client)...):", err)
		return UndefinedWindow
	}

	fmt.Printf("active_window: %v\n", activeWindow)

	// TODO: Do we actually need to do GetWMName? Shouldn't we actually do the
	// GetWindowInfo thing so we get it and much more information we could cache

	activeWindowTitle, err := ewmh.GetWMName(x.Client, activeWindow).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetWMName(x.Client, windowName)...):", err)
		return UndefinedWindow
	}

	cachedWindow := Window{
		Title: activeWindowTitle,
	}

	// TODO: Switch case to determine the window type, this will be useful for
	// simplifying automation. Needs to also detect Tor/Firefox/etc
	//   * Switch case so we can load Browser if chromium/firefox/tor etc
	// TODO: This would also fail to correctly identify browser window because for
	// example a terminal window is in the chromium or firefox folder

	downcasedTitle := strings.ToLower(cachedWindow.Title)
	switch {
	case strings.HasSuffix(downcasedTitle, "chromium"):
		cachedWindow.Type = Browser
	case strings.Contains(downcasedTitle, "firefox-esr"):
		cachedWindow.Type = Browser
	case strings.HasPrefix(downcasedTitle, "user@host:"):
		cachedWindow.Type = Terminal
	default:
		cachedWindow.Type = UndefinedType
	}

	return cachedWindow
}

func (x *X11) InitActiveWindow() Window {
	activeWindow := x.ActiveWindow()
	x.ActiveWindowTitle = activeWindow.Title
	x.ActiveWindowChangedAt = time.Now()
	return activeWindow
}

func (x *X11) CacheActiveWindow() Window {
	activeWindow := x.ActiveWindow()
	x.ActiveWindowTitle = x.ActiveWindow().Title
	x.ActiveWindowChangedAt = time.Now()
	return activeWindow
}
