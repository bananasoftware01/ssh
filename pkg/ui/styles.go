package ui

import (
	"charm.land/lipgloss/v2"
)

var (
	// Exact color palette from landing index.html
	ColorBg      = lipgloss.Color("#000000")
	ColorBg2     = lipgloss.Color("#151719")
	ColorFg      = lipgloss.Color("#ffffff")
	ColorFgDim   = lipgloss.Color("#a0aec0")
	ColorLine    = lipgloss.Color("#3a3e41")
	ColorAccent  = lipgloss.Color("#ff5e00")
	ColorComment = lipgloss.Color("#6e7681")
	ColorCyan    = lipgloss.Color("#38bdf8")
	ColorHover   = lipgloss.Color("#22262b")

	// Header
	HeaderBox = lipgloss.NewStyle().
			Background(ColorBg2).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorLine).
			Padding(0, 1)

	LogoStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	LangActiveStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	LangDimStyle = lipgloss.NewStyle().
			Foreground(ColorLine)

	// Hero
	HeroTitleStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Bold(true)

	CursorStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)

	// Editor & Lines
	LineNumStyle = lipgloss.NewStyle().
			Foreground(ColorLine)

	LineNumActiveStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	CommentStyle = lipgloss.NewStyle().
			Foreground(ColorComment).
			Italic(true)

	ProductNameStyle = lipgloss.NewStyle().
				Foreground(ColorFg).
				Bold(true)

	ProductNameActive = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	ProductUrlStyle = lipgloss.NewStyle().
			Foreground(ColorCyan)

	ProductDescStyle = lipgloss.NewStyle().
				Foreground(ColorComment)

	ArrowStyle = lipgloss.NewStyle().
			Foreground(ColorComment)

	ArrowActiveStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	RowActiveStyle = lipgloss.NewStyle().
			Background(ColorHover)

	// Footer
	FooterBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorLine).
			Foreground(ColorFgDim).
			Padding(0, 1)

	FooterTagStyle = lipgloss.NewStyle().
			Foreground(ColorLine)

	FooterValStyle = lipgloss.NewStyle().
			Foreground(ColorFgDim)

	FooterValAccent = lipgloss.NewStyle().
			Foreground(ColorAccent)
)
