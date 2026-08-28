package ui

import (
	"strings"
	"testing"

	"github.com/bananasoftware/banana-ssh/pkg/data"
)

func TestSplashIntro(t *testing.T) {
	app := NewAppModel("127.0.0.1", "testuser")
	app.Width = 100
	app.Height = 24

	if !app.InSplash {
		t.Fatal("expected initial InSplash to be true")
	}

	splashView := app.View().Content
	if !strings.Contains(splashView, "Banana Software") {
		t.Errorf("expected splash view to contain 'Banana Software'")
	}
	if strings.Contains(splashView, "# sản phẩm") {
		t.Errorf("splash view should not contain products list")
	}

	// Receive SplashEndMsg after 2 seconds
	updatedModel, _ := app.Update(SplashEndMsg{})
	app = updatedModel.(AppModel)
	if app.InSplash {
		t.Errorf("expected InSplash to be false after SplashEndMsg")
	}

	mainView := app.View().Content
	if !strings.Contains(mainView, "# sản phẩm") {
		t.Errorf("expected main view to contain '# sản phẩm' after splash")
	}
}

func TestAppModelSimpleRendering(t *testing.T) {
	app := NewAppModel("127.0.0.1", "testuser")
	app.InSplash = false
	app.Width = 100
	app.Height = 24

	// 1. Initial VI view
	viewVI := app.View().Content
	if !strings.Contains(viewVI, "bananasoftware.net") {
		t.Errorf("expected header 'bananasoftware.net'")
	}
	if !strings.Contains(viewVI, "Banana Software") {
		t.Errorf("expected hero 'Banana Software'")
	}
	if !strings.Contains(viewVI, "█") {
		t.Errorf("expected cursor '█' when CursorBlink is true")
	}
	if !strings.Contains(viewVI, "# sản phẩm") {
		t.Errorf("expected comment '# sản phẩm'")
	}
	if !strings.Contains(viewVI, "TaiHoaDon") {
		t.Errorf("expected product 'TaiHoaDon'")
	}
	if !strings.Contains(viewVI, "Tạo phần mềm mới") {
		t.Errorf("expected 'Tạo phần mềm mới' in view")
	}
	if !strings.Contains(viewVI, "taihoadon.online") {
		t.Errorf("expected url 'taihoadon.online'")
	}
	if !strings.Contains(viewVI, "2803238388") {
		t.Errorf("expected tax id in footer")
	}

	// 2. Language toggle to EN
	app.ToggleLang()
	if app.Lang != data.LangEN {
		t.Errorf("expected lang EN, got %s", app.Lang)
	}
	viewEN := app.View().Content
	if !strings.Contains(viewEN, "# products") {
		t.Errorf("expected EN comment '# products'")
	}
	if !strings.Contains(viewEN, "Add new software") {
		t.Errorf("expected 'Add new software' in EN view")
	}
}

func TestAddNewFeature(t *testing.T) {
	app := NewAppModel("127.0.0.1", "testuser")
	app.InSplash = false

	// Cursor at index 5 ("Tạo phần mềm mới")
	app.Cursor = len(data.Products)
	if app.InputMode {
		t.Error("InputMode should initially be false")
	}

	// Press Enter to open input
	app.InputMode = true
	app.InputArea.Focus()
	app.InputArea.SetValue("App bán hàng tự chốt đơn 0335581402")

	view := app.View().Content
	if !strings.Contains(view, "App bán hàng tự chốt đơn") {
		t.Errorf("expected input value in view when InputMode is true")
	}

	// Submit text
	cmd := sendTelegramCmd(app.InputArea.Value(), app.ClientAddr, app.SSHUser, string(app.Lang))
	if cmd == nil {
		t.Errorf("expected telegram cmd to be created")
	}
}

func TestCursorBlink(t *testing.T) {
	app := NewAppModel("127.0.0.1", "testuser")
	app.InSplash = false
	if !app.CursorBlink {
		t.Errorf("expected initial CursorBlink to be true")
	}

	// Send BlinkMsg
	updatedModel, cmd := app.Update(BlinkMsg{})
	app = updatedModel.(AppModel)
	if app.CursorBlink {
		t.Errorf("expected CursorBlink to toggle to false")
	}
	if cmd == nil {
		t.Errorf("expected blinkCmd to be scheduled")
	}

	// Check view when cursor is hidden
	view := app.View().Content
	if strings.Contains(view, "Banana Software█") {
		t.Errorf("cursor should not be visible when CursorBlink is false")
	}

	// Send second BlinkMsg
	updatedModel, _ = app.Update(BlinkMsg{})
	app = updatedModel.(AppModel)
	if !app.CursorBlink {
		t.Errorf("expected CursorBlink to toggle back to true")
	}
}

func TestCursorNavigation(t *testing.T) {
	app := NewAppModel("127.0.0.1", "testuser")
	app.InSplash = false
	if app.Cursor != 0 {
		t.Errorf("expected initial cursor 0, got %d", app.Cursor)
	}

	// Navigate down to Add new (index 5)
	app.Cursor = len(data.Products)
	if app.Cursor != 5 {
		t.Errorf("expected cursor 5, got %d", app.Cursor)
	}
}

func TestResponsiveLayout(t *testing.T) {
	app := NewAppModel("127.0.0.1", "testuser")
	app.InSplash = false

	// Narrow (50 cols)
	app.Width = 50
	narrowView := app.View().Content
	if narrowView == "" {
		t.Error("expected non-empty narrow view")
	}

	// Wide (140 cols, centered)
	app.Width = 140
	wideView := app.View().Content
	if wideView == "" {
		t.Error("expected non-empty wide view")
	}
	if !strings.Contains(wideView, "Banana Software") {
		t.Error("expected centered view to contain Banana Software")
	}
}
