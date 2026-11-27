package draw

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/foam/theme"
)

func TestFrameSize(t *testing.T) {
	for _, o := range []FrameOpts{
		{},
		{Title: "Results", Right: "18", Footer: "2 048 req"},
		{Title: strings.Repeat("x", 40), Right: strings.Repeat("y", 40)},
		{Scroll: Scroll{Total: 100, Visible: 5, Offset: 50}},
	} {
		lines := strings.Split(Frame("a\nb", 20, 7, o), "\n")
		if len(lines) != 7 {
			t.Fatalf("%+v: %d lines", o, len(lines))
		}
		for _, l := range lines {
			if w := ansi.StringWidth(l); w != 20 {
				t.Fatalf("%+v: line %q is %d wide", o, ansi.Strip(l), w)
			}
		}
	}
}

func TestThumb(t *testing.T) {
	if _, _, ok := (Scroll{Total: 5, Visible: 10}).Thumb(10); ok {
		t.Fatal("no bar when everything fits")
	}
	s, e, ok := Scroll{Total: 100, Visible: 10, Offset: 90}.Thumb(10)
	if !ok || e != 10 || s != 9 {
		t.Fatalf("bottom thumb = %d..%d", s, e)
	}
}

func TestHatchSize(t *testing.T) {
	lines := strings.Split(Hatch(30, 5, "no response"), "\n")
	if len(lines) != 5 || ansi.StringWidth(lines[2]) != 30 || !strings.Contains(ansi.Strip(lines[2]), "no response") {
		t.Fatalf("%q", lines)
	}
}

func TestFillRow(t *testing.T) {
	c := lipgloss.NewCanvas(6, 2)
	c.Compose(lipgloss.NewCompositor(lipgloss.NewLayer(theme.Primary().Render("ab") + "    \nxxxxxx")))
	FillRow(c, 0, theme.Surface())
	for x := range 6 {
		if c.CellAt(x, 0).Style.Bg == nil {
			t.Fatalf("cell %d of the bar has no background", x)
		}
	}
	if c.CellAt(0, 1).Style.Bg != nil {
		t.Fatal("another row was filled")
	}
}

// Wide characters keep their two cells through Dim and FillRow.
func TestWideCharacters(t *testing.T) {
	c := lipgloss.NewCanvas(6, 1)
	c.Compose(lipgloss.NewCompositor(lipgloss.NewLayer("一二三")))
	Dim(c)
	FillRow(c, 0, theme.Surface())
	if got := ansi.Strip(c.Render()); got != "一二三" || ansi.StringWidth(c.Render()) != 6 {
		t.Fatalf("got %q", got)
	}
}

func TestClean(t *testing.T) {
	if got := Clean("a\x1b[2J\u009b\xff\tb"); got != "a·[2J·· b" {
		t.Fatalf("got %q", got)
	}
	if s := "plain é 一"; Clean(s) != s {
		t.Fatal("clean text changed")
	}
}

func TestSelectedKeepsColors(t *testing.T) {
	row := " " + theme.Danger().Render("red") + " plain " + theme.Success().Render("green")
	got := Selected(row, 30, true)
	if ansi.StringWidth(got) != 30 {
		t.Fatalf("width = %d", ansi.StringWidth(got))
	}
	if !strings.Contains(got, theme.Danger().Render("red")) || !strings.Contains(got, theme.Success().Render("green")) {
		t.Fatalf("colors lost: %q", got)
	}
	if ansi.Strip(got) != ansi.Strip(Pad(row, 30)) {
		t.Fatalf("text changed: %q", ansi.Strip(got))
	}
	bg := lipgloss.NewStyle().Background(theme.Selection().Value()).Foreground(theme.Fg().Value()).Render("x")
	bg, _, _ = strings.Cut(bg, "x")
	if strings.Count(got, bg) < 3 {
		t.Fatalf("background not kept after the colored parts: %q", got)
	}
}
