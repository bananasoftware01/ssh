package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bananasoftware/banana-ssh/pkg/data"
)

type BlinkMsg struct{}
type SplashEndMsg struct{}

func blinkCmd() tea.Cmd {
	return tea.Tick(530*time.Millisecond, func(t time.Time) tea.Msg {
		return BlinkMsg{}
	})
}

func splashTimerCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return SplashEndMsg{}
	})
}
type AppModel struct {
	Cursor      int
	Lang        data.Lang
	Width       int
	Height      int
	Quitting    bool
	StatusMsg   string
	CursorBlink bool
	InSplash    bool
	InputMode   bool
	InputArea   textinput.Model
	ClientAddr  string
	SSHUser     string
}

func NewAppModel(clientAddr, sshUser string) AppModel {
	ti := textinput.New()
	ti.Placeholder = "Nhập mô tả ý tưởng phần mềm & SĐT/Zalo của bạn..."
	ti.CharLimit = 200
	ti.SetWidth(56)

	return AppModel{
		Cursor:      0,
		Lang:        data.LangVI,
		Width:       80,
		Height:      24,
		CursorBlink: true,
		InSplash:    true,
		InputArea:   ti,
		ClientAddr:  clientAddr,
		SSHUser:     sshUser,
	}
}

type TelegramSentMsg struct{}

func sendTelegramCmd(text, clientAddr, sshUser, lang string) tea.Cmd {
	return func() tea.Msg {
		botToken := "8891045799:AAFSwsk-1bE9oKfYtTr2YTgRAchcigmOhTc"
		chatID := "-5326514969"

		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			loc = time.UTC
		}
		timeStr := time.Now().In(loc).Format("02/01/2006 15:04:05 MST")

		userTag := sshUser
		if userTag == "" {
			userTag = "anonymous"
		}
		addrTag := clientAddr
		if addrTag == "" {
			addrTag = "local-tui"
		}

		msg := fmt.Sprintf(
			"🚀 <b>YÊU CẦU TẠO PHẦN MỀM MỚI (SSH TERMINAL)</b>\n\n"+
				"💡 <b>Nội dung yêu cầu / SĐT / Zalo:</b>\n<blockquote>%s</blockquote>\n\n"+
				"🌐 <b>Ngôn ngữ:</b> %s\n"+
				"💻 <b>SSH User:</b> <code>%s</code>\n"+
				"📍 <b>Client Address:</b> <code>%s</code>\n"+
				"📡 <b>Nguồn:</b> <code>ssh bananasoftware.net</code>\n"+
				"⏰ <b>Thời gian:</b> %s",
			html.EscapeString(text),
			html.EscapeString(lang),
			html.EscapeString(userTag),
			html.EscapeString(addrTag),
			html.EscapeString(timeStr),
		)

		payload, _ := json.Marshal(map[string]string{
			"chat_id":    chatID,
			"text":       msg,
			"parse_mode": "HTML",
		})

		client := &http.Client{Timeout: 8 * time.Second}
		_, _ = client.Post(
			fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken),
			"application/json",
			bytes.NewBuffer(payload),
		)
		return TelegramSentMsg{}
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
		blinkCmd(),
		splashTimerCmd(),
	)
}

func (m *AppModel) ToggleLang() {
	if m.Lang == data.LangVI {
		m.Lang = data.LangEN
		m.InputArea.Placeholder = "Describe software requirement & your Phone/Zalo..."
	} else {
		m.Lang = data.LangVI
		m.InputArea.Placeholder = "Nhập mô tả ý tưởng phần mềm & SĐT/Zalo của bạn..."
	}
}


func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SplashEndMsg:
		m.InSplash = false
		return m, nil

	case BlinkMsg:
		if !m.Quitting {
			m.CursorBlink = !m.CursorBlink
			return m, blinkCmd()
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		if m.InSplash {
			if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "Q" {
				m.Quitting = true
				return m, tea.Quit
			}
			m.InSplash = false
			return m, nil
		}

		if m.InputMode {
			switch msg.String() {
			case "ctrl+c":
				m.Quitting = true
				return m, tea.Quit
			case "esc":
				m.InputMode = false
				m.InputArea.Reset()
				return m, nil
			case "enter":
				text := strings.TrimSpace(m.InputArea.Value())
				if text != "" {
					m.InputMode = false
					m.InputArea.Reset()
					if m.Lang == data.LangVI {
						m.StatusMsg = "✔ Đã gửi yêu cầu tạo phần mềm đến kỹ sư Banana Software!"
					} else {
						m.StatusMsg = "✔ Request sent to Banana Software engineers!"
					}
					return m, sendTelegramCmd(text, m.ClientAddr, m.SSHUser, string(m.Lang))
				}
				m.InputMode = false
				return m, nil
			}

			var cmd tea.Cmd
			m.InputArea, cmd = m.InputArea.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "Q", "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "l", "L":
			m.ToggleLang()
			m.StatusMsg = ""

		case "up", "k", "K":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(data.Products) // 0..5 (5 is "Tạo phần mềm mới")
			}
			m.StatusMsg = ""

		case "down", "j", "J":
			if m.Cursor < len(data.Products) {
				m.Cursor++
			} else {
				m.Cursor = 0
			}
			m.StatusMsg = ""

		case "enter":
			if m.Cursor == len(data.Products) {
				// Open inline input
				m.InputMode = true
				m.InputArea.Focus()
				m.StatusMsg = ""
				return m, nil
			} else if m.Cursor >= 0 && m.Cursor < len(data.Products) {
				p := data.Products[m.Cursor]
				m.StatusMsg = fmt.Sprintf("➜ %s (%s)", p.Name, p.URL)
			}
		}
	}

	return m, nil
}

func (m AppModel) View() tea.View {
	if m.Quitting {
		var bye string
		if m.Lang == data.LangVI {
			bye = "Cảm ơn bạn đã ghé thăm Banana Software!\n"
		} else {
			bye = "Thank you for visiting Banana Software!\n"
		}
		v := tea.NewView(bye)
		return v
	}
	var cursorChar string
	if m.CursorBlink {
		cursorChar = CursorStyle.Render("█")
	} else {
		cursorChar = " "
	}

	if m.InSplash {
		hero := HeroTitleStyle.Render("Banana Software") + cursorChar
		splash := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, hero)
		v := tea.NewView(splash)
		v.AltScreen = true
		return v
	}

	maxW := 106
	w := m.Width - 4
	if w > maxW {
		w = maxW
	}
	if w < 50 {
		w = 50
	}

	// 1. Header (Fixed at top, with leading empty line)
	logo := LogoStyle.Render("bananasoftware.net")
	var langSwitch string
	if m.Lang == data.LangVI {
		langSwitch = LangActiveStyle.Render("VI") + " " + LangDimStyle.Render("/") + " " + LangDimStyle.Render("EN")
	} else {
		langSwitch = LangDimStyle.Render("VI") + " " + LangDimStyle.Render("/") + " " + LangActiveStyle.Render("EN")
	}

	leftHeader := " " + logo
	rightHeader := langSwitch
	innerHeaderWidth := (w - 2) - HeaderBox.GetHorizontalFrameSize()
	headerGap := innerHeaderWidth - lipgloss.Width(leftHeader) - lipgloss.Width(rightHeader)
	if headerGap < 1 {
		headerGap = 1
	}
	headerSpacer := lipgloss.NewStyle().Width(headerGap).Render("")
	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, leftHeader, headerSpacer, rightHeader)
	headerContent := "\n" + HeaderBox.Width(w - 2).Render(headerRow) + "\n"
	// 2. Main Body (Hero + Editor)
	var mainSb strings.Builder
	hero := HeroTitleStyle.Render("Banana Software") + cursorChar
	mainSb.WriteString("  " + hero)
	mainSb.WriteString("\n\n")
	commentText := "# sản phẩm"
	if m.Lang == data.LangEN {
		commentText = "# products"
	}
	line1Num := LineNumStyle.Render("  1 │ ")
	line1Body := CommentStyle.Render(commentText)
	mainSb.WriteString(line1Num + line1Body)
	mainSb.WriteString("\n")

	maxNameLen := 14
	maxUrlLen := 18

	for i, p := range data.Products {
		isSelected := m.Cursor == i
		lineNum := fmt.Sprintf("%3d │ ", i+2)

		var numStyled string
		var nameStyled string
		var urlStyled string
		var descStyled string
		var arrowStyled string

		desc := p.DescVI
		if m.Lang == data.LangEN {
			desc = p.DescEN
		}

		if isSelected {
			numStyled = LineNumActiveStyle.Render(lineNum)
			nameStyled = ProductNameActive.Render(fmt.Sprintf("%-*s", maxNameLen, p.Name))
			urlStyled = lipgloss.NewStyle().Foreground(ColorAccent).Underline(true).Hyperlink(p.URL).Render(fmt.Sprintf("%-*s", maxUrlLen, p.DisplayURL))
			descStyled = lipgloss.NewStyle().Foreground(ColorFg).Render("// " + desc)
			arrowStyled = ArrowActiveStyle.Render("↗")
		} else {
			numStyled = LineNumStyle.Render(lineNum)
			nameStyled = ProductNameStyle.Render(fmt.Sprintf("%-*s", maxNameLen, p.Name))
			urlStyled = ProductUrlStyle.Hyperlink(p.URL).Render(fmt.Sprintf("%-*s", maxUrlLen, p.DisplayURL))
			descStyled = ProductDescStyle.Render("// " + desc)
			arrowStyled = ArrowStyle.Render("↗")
		}

		availableWidth := w - 6 - maxNameLen - maxUrlLen - 12
		if availableWidth > 15 {
			descRunes := []rune(desc)
			if len(descRunes) > availableWidth {
				desc = string(descRunes[:availableWidth-3]) + "..."
			}
			if isSelected {
				descStyled = lipgloss.NewStyle().Foreground(ColorFg).Render("// " + desc)
			} else {
				descStyled = ProductDescStyle.Render("// " + desc)
			}
			row := fmt.Sprintf("%s %s %s  %s", nameStyled, urlStyled, descStyled, arrowStyled)
			if isSelected {
				mainSb.WriteString(numStyled + RowActiveStyle.Render(" "+row+" "))
			} else {
				mainSb.WriteString(numStyled + " " + row)
			}
		} else {
			row := fmt.Sprintf("%s %s  %s", nameStyled, urlStyled, arrowStyled)
			if isSelected {
				mainSb.WriteString(numStyled + RowActiveStyle.Render(" "+row+" "))
			} else {
				mainSb.WriteString(numStyled + " " + row)
			}
		}
		mainSb.WriteString("\n")
	}

	// Line 7: "Tạo phần mềm mới" (Add new software)
	isAddSelected := m.Cursor == len(data.Products)
	line7Num := LineNumStyle.Render("  7 │ ")
	var addName string
	var addDesc string
	if m.Lang == data.LangVI {
		addName = "Tạo phần mềm mới"
		addDesc = "// Đặt hàng tính năng hoặc phần mềm theo yêu cầu"
	} else {
		addName = "Add new software"
		addDesc = "// Request custom feature or software build"
	}

	var addRow string
	if isAddSelected {
		line7Num = LineNumActiveStyle.Render("  7 │ ")
		addNameStyled := ProductNameActive.Render(fmt.Sprintf("%-*s", maxNameLen+maxUrlLen+1, addName))
		addDescStyled := lipgloss.NewStyle().Foreground(ColorFg).Render(addDesc)
		addArrow := ArrowActiveStyle.Render("+")
		addRow = fmt.Sprintf("%s %s  %s", addNameStyled, addDescStyled, addArrow)
		mainSb.WriteString(line7Num + RowActiveStyle.Render(" "+addRow+" ") + "\n")
	} else {
		addNameStyled := ProductNameStyle.Render(fmt.Sprintf("%-*s", maxNameLen+maxUrlLen+1, addName))
		addDescStyled := ProductDescStyle.Render(addDesc)
		addArrow := ArrowStyle.Render("+")
		addRow = fmt.Sprintf("%s %s  %s", addNameStyled, addDescStyled, addArrow)
		mainSb.WriteString(line7Num + " " + addRow + "\n")
	}

	// Line 8: Inline Input (when active)
	if m.InputMode {
		m.InputArea.SetWidth(w - 38)
		line8Num := LineNumStyle.Render("  8 │ ")
		promptPrefix := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("❯ ")
		inputHints := CommentStyle.Render("  [Enter: Gửi | Esc: Đóng]")
		if m.Lang == data.LangEN {
			inputHints = CommentStyle.Render("  [Enter: Submit | Esc: Close]")
		}
		mainSb.WriteString(line8Num + " " + promptPrefix + m.InputArea.View() + inputHints + "\n")
		mainSb.WriteString(LineNumStyle.Render("  9 │ "))
	} else {
		mainSb.WriteString(LineNumStyle.Render("  8 │ "))
	}

	if m.StatusMsg != "" {
		mainSb.WriteString("\n\n")
		statusColor := ColorAccent
		if strings.HasPrefix(m.StatusMsg, "✔") {
			statusColor = lipgloss.Color("#10b981")
		}
		mainSb.WriteString(lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render("  " + m.StatusMsg))
	}
	mainBody := mainSb.String() + "\n"
	// 3. Footer (Fixed at bottom)
	var compName string
	if m.Lang == data.LangVI {
		compName = data.Company.NameVI
	} else {
		compName = data.Company.NameEN
	}

	footerLeft := fmt.Sprintf("%s  │  %s %s  │  %s %s",
		FooterValStyle.Render(compName),
		FooterTagStyle.Render("MST:"),
		FooterValStyle.Render(data.Company.TaxID),
		FooterTagStyle.Render("SĐT/Zalo:"),
		FooterValAccent.Render(data.Company.Phone),
	)
	footerRight := FooterTagStyle.Render("© " + data.Company.Year + " Banana Software")

	innerFooterWidth := (w - 2) - FooterBox.GetHorizontalFrameSize()
	footerGap := innerFooterWidth - lipgloss.Width(footerLeft) - lipgloss.Width(footerRight)
	if footerGap < 1 {
		footerGap = 1
	}
	footerSpacer := lipgloss.NewStyle().Width(footerGap).Render("")
	footerRow := lipgloss.JoinHorizontal(lipgloss.Top, footerLeft, footerSpacer, footerRight)
	hints := "  " + FooterTagStyle.Render("↑/↓ / j/k: Chọn  │  Enter: Mở link  │  l: Đổi ngôn ngữ  │  q / Esc: Thoát")
	footerContent := FooterBox.Width(w - 2).Render(footerRow) + "\n" + hints

	// Vertical centering calculation
	headerLines := strings.Count(headerContent, "\n")
	footerLines := strings.Count(footerContent, "\n") + 1
	mainLines := strings.Count(mainBody, "\n")

	targetH := m.Height
	if targetH < 20 {
		targetH = 20
	}

	middleSpace := targetH - headerLines - footerLines
	topPad := (middleSpace - mainLines) / 2
	if topPad < 1 {
		topPad = 1
	}
	bottomPad := middleSpace - mainLines - topPad
	if bottomPad < 1 {
		bottomPad = 1
	}

	var fullSb strings.Builder
	fullSb.WriteString(headerContent)
	fullSb.WriteString(strings.Repeat("\n", topPad))
	fullSb.WriteString(mainBody)
	fullSb.WriteString(strings.Repeat("\n", bottomPad))
	fullSb.WriteString(footerContent)

	fullContent := fullSb.String()
	leftMargin := (m.Width - w) / 2
	if leftMargin > 0 {
		fullContent = lipgloss.NewStyle().MarginLeft(leftMargin).Render(fullContent)
	}

	v := tea.NewView(fullContent)
	v.AltScreen = true
	return v
}
