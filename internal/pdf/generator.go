package pdf

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed report.html
var templateContent string

type BlockData struct {
	Type           string // finding, note, warning, tip, important, divider, checklist, statistics
	FindingTitle   string
	FindingDesc    string
	FindingImage   string
	Category       string
	Priority       string
	Severity       string
	Status         string
	Recommendation string
	NoteContent    string
	MessageText    string
	DividerTitle   string
	ChecklistItems []ChecklistItemData
	RenderedHTML   string // Pre-rendered HTML for template insertion
}

type ChecklistItemData struct {
	Text    string
	Checked bool
}

type ReportData struct {
	Title          string
	Project        string
	Company        string
	Author         string
	Version        string
	Date           string
	Classification string
	ReportID       string
	Status         string
	Watermark      string

	Blocks      []BlockData
	Screenshots []ScreenshotData
	Summary     SummaryData
}

type ScreenshotData struct {
	ImagePath      string
	Title          string
	Description    string
	FigureLabel    string
	Category       string
	Priority       string
	Severity       string
	Status         string
	Recommendation string
}

type SummaryData struct {
	TotalFindings   int
	TotalImages     int
	HighCount       int
	MediumCount     int
	LowCount        int
	Categories      []string
	ReadingTime     string
	HasRecommendations bool
}

type Generator struct {
	tmpl    *template.Template
	timeout time.Duration
}

func NewGenerator() (*Generator, error) {
	funcMap := template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"join":  func(elems []string, sep string) string { return strings.Join(elems, sep) },
		"lower": func(s string) string { return strings.ToLower(s) },
		"seq":   func(n int) []int { s := make([]int, n); for i := range s { s[i] = i + 1 }; return s },
	}
	tmpl, err := template.New("report").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("parse report template: %w", err)
	}
	return &Generator{
		tmpl:    tmpl,
		timeout: 60 * time.Second,
	}, nil
}

func (g *Generator) preRenderBlocks(data *ReportData) {
	ctx := RenderContext{Summary: data.Summary}
	for i := range data.Blocks {
		data.Blocks[i].RenderedHTML = RenderBlock(data.Blocks[i], ctx)
	}
}

func (g *Generator) RenderHTML(data ReportData) (string, error) {
	g.preRenderBlocks(&data)
	var buf strings.Builder
	if err := g.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) Generate(data ReportData, outputPath string) error {
	g.preRenderBlocks(&data)

	tempDir, err := os.MkdirTemp("", "snapreport-*")
	if err != nil {
		return fmt.Errorf("create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			slog.Warn("failed to clean up temp directory", "path", tempDir, "error", err)
		}
	}()

	for i, ss := range data.Screenshots {
		ext := filepath.Ext(ss.ImagePath)
		dest := filepath.Join(tempDir, fmt.Sprintf("img_%d%s", i, ext))
		if err := copyFile(ss.ImagePath, dest); err != nil {
			return fmt.Errorf("copy image: %w", err)
		}
		data.Screenshots[i].ImagePath = filepath.Base(dest)
	}

	htmlPath := filepath.Join(tempDir, "index.html")
	f, err := os.Create(htmlPath)
	if err != nil {
		return fmt.Errorf("create html file: %w", err)
	}
	if err := g.tmpl.Execute(f, data); err != nil {
		f.Close()
		return fmt.Errorf("execute template: %w", err)
	}
	f.Close()

	ctx, allocCancel := g.newContext()
	defer allocCancel()

	ctx, timeoutCancel := context.WithTimeout(ctx, g.timeout)
	defer timeoutCancel()

	var pdfBuf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate("file://"+htmlPath),
		chromedp.WaitReady("body"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0.4).
				WithMarginBottom(0.6).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				WithDisplayHeaderFooter(true).
				WithHeaderTemplate(`<div style="font-size:7.5pt;color:#94a3b8;width:100%;text-align:center;padding:0 25px;border-bottom:0.5px solid #e2e8f0;padding-bottom:4pt;">
					<span style="float:left;">`+template.HTMLEscapeString(data.Company)+`</span>
					<span>`+template.HTMLEscapeString(data.Classification)+`</span>
					<span style="float:right;">`+template.HTMLEscapeString(data.Project)+`</span>
				</div>`).
				WithFooterTemplate(`<div style="font-size:7.5pt;color:#94a3b8;width:100%;text-align:center;border-top:0.5px solid #e2e8f0;padding-top:4pt;">
					<span style="float:left;">`+template.HTMLEscapeString(data.Title)+`</span>
					<span class="pageNumber"></span> / <span class="totalPages"></span>
					<span style="float:right;">v`+template.HTMLEscapeString(data.Version)+`</span>
				</div>`).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return fmt.Errorf("generate pdf with chromium: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	return os.WriteFile(outputPath, pdfBuf, 0644)
}

func (g *Generator) newContext() (context.Context, context.CancelFunc) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.NoSandbox,
			chromedp.Headless,
			chromedp.DisableGPU,
			chromedp.Flag("disable-dev-shm-usage", true),
		)...,
	)
	ctx, _ := chromedp.NewContext(allocCtx)
	return ctx, cancel
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer d.Close()

	if _, err := io.Copy(d, s); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}
	return nil
}
