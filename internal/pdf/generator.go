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
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed report.html
var templateContent string

type ReportData struct {
	Title       string
	Project     string
	Company     string
	Author      string
	Version     string
	Date        string
	Screenshots []ScreenshotData
}

type ScreenshotData struct {
	ImagePath   string
	Title       string
	Description string
}

type Generator struct {
	tmpl    *template.Template
	timeout time.Duration
}

func NewGenerator() (*Generator, error) {
	tmpl, err := template.New("report").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("parse report template: %w", err)
	}
	return &Generator{
		tmpl:    tmpl,
		timeout: 60 * time.Second,
	}, nil
}

func (g *Generator) Generate(data ReportData, outputPath string) error {
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
				WithHeaderTemplate(fmt.Sprintf(
					`<div style="font-size:8pt;color:#94a3b8;width:100%%;text-align:center;padding:0 20px;">%s</div>`,
					template.HTMLEscapeString(data.Title),
				)).
				WithFooterTemplate(`<div style="font-size:8pt;color:#94a3b8;width:100%;text-align:center;"><span class="pageNumber"></span> / <span class="totalPages"></span></div>`).
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
