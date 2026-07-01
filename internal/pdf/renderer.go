package pdf

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

type RenderContext struct {
	FindingIndex int
	Summary      SummaryData
}

type BlockRenderer interface {
	Render(block BlockData, ctx RenderContext) string
}

var registry map[string]BlockRenderer

func init() {
	registry = map[string]BlockRenderer{
		"finding":    findingRenderer{},
		"note":       noteRenderer{},
		"warning":    calloutRenderer{},
		"tip":        calloutRenderer{},
		"important":  calloutRenderer{},
		"divider":    dividerRenderer{},
		"checklist":  checklistRenderer{},
		"statistics": statisticsRenderer{},
		"markdown":   markdownRenderer{},
		"code":       codeRenderer{},
		"table":      tableRenderer{},
		"image":      imageRenderer{},
		"recommendation": recommendationRenderer{},
	}
}

func RenderBlock(block BlockData, ctx RenderContext) string {
	r, ok := registry[block.Type]
	if !ok {
		return fmt.Sprintf(`<div class="block-error">Unknown block: %s</div>`, template.HTMLEscapeString(block.Type))
	}
	return r.Render(block, ctx)
}

func esc(s string) string {
	return template.HTMLEscapeString(s)
}

// ─── Finding Renderer ───

type findingRenderer struct{}

func (findingRenderer) Render(b BlockData, ctx RenderContext) string {
	ctx.FindingIndex++
	idx := ctx.FindingIndex
	title := b.FindingTitle
	if title == "" {
		title = fmt.Sprintf("Finding %d", idx)
	}

	var html strings.Builder
	fmt.Fprintf(&html, `<div class="finding">`)
	// Header: "Finding 01 — Title"
	fmt.Fprintf(&html, `<div class="finding-head">`)
	fmt.Fprintf(&html, `<span class="finding-num">Finding %02d</span>`, idx)
	fmt.Fprintf(&html, `<h2 class="finding-title">%s</h2>`, esc(title))
	fmt.Fprintf(&html, `</div>`)

	// Metadata panel: 2-column grid
	var metaItems []string
	fields := [][2]string{{"Category", b.Category}, {"Priority", b.Priority}, {"Severity", b.Severity}, {"Status", b.Status}}
	for _, f := range fields {
		if f[1] != "" {
			metaItems = append(metaItems,
				fmt.Sprintf(`<div class="meta-item"><span class="meta-label">%s</span><span class="meta-value">%s</span></div>`,
					esc(f[0]), esc(f[1])))
		}
	}
	if len(metaItems) > 0 {
		html.WriteString(`<div class="finding-meta">`)
		for _, item := range metaItems {
			html.WriteString(item)
		}
		html.WriteString(`</div>`)
	}

	if b.FindingImage != "" {
		fmt.Fprintf(&html, `<div class="image-frame"><img src="%s" alt="%s"></div>`, esc(b.FindingImage), esc(title))
	}
	if b.FindingDesc != "" {
		fmt.Fprintf(&html, `<div class="desc-card"><h3>Description</h3><p>%s</p></div>`, esc(b.FindingDesc))
	}
	if b.Recommendation != "" {
		fmt.Fprintf(&html, `<div class="rec-card"><h3>Recommendation</h3><p>%s</p></div>`, esc(b.Recommendation))
	}
	html.WriteString(`</div>`)
	return html.String()
}

// ─── Note Renderer ───

type noteRenderer struct{}

func (noteRenderer) Render(b BlockData, ctx RenderContext) string {
	content := b.NoteContent
	if content == "" {
		content = b.MessageText // fallback
	}
	return fmt.Sprintf(`<div class="block-note"><p>%s</p></div>`, esc(content))
}

// ─── Callout Renderer (warning / tip / important) ───

type calloutRenderer struct{}

func (calloutRenderer) Render(b BlockData, ctx RenderContext) string {
	label := map[string]string{"warning": "Warning", "tip": "Tip", "important": "Important"}
	l := label[b.Type]
	if l == "" {
		l = b.Type
	}
	return fmt.Sprintf(`<div class="block-message %s"><p class="block-message-label">%s</p><p>%s</p></div>`,
		esc(b.Type), esc(l), esc(b.MessageText))
}

// ─── Divider Renderer ───

type dividerRenderer struct{}

func (dividerRenderer) Render(b BlockData, ctx RenderContext) string {
	var html strings.Builder
	html.WriteString(`<div class="block-divider"><div class="block-divider-line"></div>`)
	if b.DividerTitle != "" {
		fmt.Fprintf(&html, `<span class="block-divider-title">%s</span><div class="block-divider-line"></div>`, esc(b.DividerTitle))
	}
	html.WriteString(`</div>`)
	return html.String()
}

// ─── Checklist Renderer ───

type checklistRenderer struct{}

func (checklistRenderer) Render(b BlockData, ctx RenderContext) string {
	var html strings.Builder
	html.WriteString(`<div class="block-checklist">`)
	for _, item := range b.ChecklistItems {
		cls := ""
		txtCls := ""
		if item.Checked {
			cls = " checked"
			txtCls = " done"
		}
		fmt.Fprintf(&html, `<div class="checklist-item"><div class="checklist-box%s"></div><span class="checklist-text%s">%s</span></div>`, cls, txtCls, esc(item.Text))
	}
	html.WriteString(`</div>`)
	return html.String()
}

// ─── Statistics Renderer ───

type statisticsRenderer struct{}

func (statisticsRenderer) Render(b BlockData, ctx RenderContext) string {
	s := ctx.Summary
	return fmt.Sprintf(`<div class="stats-grid" style="margin-top:12pt;">
		<div class="stat-card"><h3>Total Findings</h3><p>%d</p></div>
		<div class="stat-card"><h3>Screenshots</h3><p>%d</p></div>
	</div>`, s.TotalFindings, s.TotalImages)
}

// ─── Markdown Renderer ───

type markdownRenderer struct{}

func (markdownRenderer) Render(b BlockData, ctx RenderContext) string {
	// Parse content as JSON {html: "..."} or use raw content
	var c struct {
		HTML string `json:"html"`
		Text string `json:"text"`
	}
	json.Unmarshal([]byte(b.MessageText), &c)
	content := c.HTML
	if content == "" {
		content = c.Text
	}
	if content == "" {
		content = b.MessageText
	}
	return fmt.Sprintf(`<div class="block-note"><p>%s</p></div>`, esc(content))
}

// ─── Code Renderer ───

type codeRenderer struct{}

func (codeRenderer) Render(b BlockData, ctx RenderContext) string {
	var c struct {
		Code     string `json:"code"`
		Language string `json:"language"`
	}
	json.Unmarshal([]byte(b.NoteContent), &c)
	code := c.Code
	if code == "" {
		code = b.FindingDesc
	}
	return fmt.Sprintf(`<div class="block-code"><pre><code>%s</code></pre></div>`, esc(code))
}

// ─── Table Renderer ───

type tableRenderer struct{}

func (tableRenderer) Render(b BlockData, ctx RenderContext) string {
	var c struct {
		Headers []string   `json:"headers"`
		Rows    [][]string `json:"rows"`
	}
	json.Unmarshal([]byte(b.ChecklistItems[0].Text), &c) // placeholder
	var html strings.Builder
	html.WriteString(`<table class="block-table"><thead><tr>`)
	for _, h := range c.Headers {
		fmt.Fprintf(&html, "<th>%s</th>", esc(h))
	}
	html.WriteString(`</tr></thead><tbody>`)
	for _, row := range c.Rows {
		html.WriteString(`<tr>`)
		for _, cell := range row {
			fmt.Fprintf(&html, "<td>%s</td>", esc(cell))
		}
		html.WriteString(`</tr>`)
	}
	html.WriteString(`</tbody></table>`)
	return html.String()
}

// ─── Image Renderer ───

type imageRenderer struct{}

func (imageRenderer) Render(b BlockData, ctx RenderContext) string {
	var c struct {
		Src     string `json:"src"`
		Caption string `json:"caption"`
	}
	json.Unmarshal([]byte(b.FindingTitle), &c)
	return fmt.Sprintf(`<div class="block-image"><div class="image-frame"><img src="%s" alt="%s"></div>%s</div>`,
		esc(c.Src), esc(c.Caption),
		func() string { if c.Caption != "" { return fmt.Sprintf(`<p class="caption" style="text-align:center;margin-top:4pt;">%s</p>`, esc(c.Caption)) }; return "" }())
}

// ─── Recommendation Renderer ───

type recommendationRenderer struct{}

func (recommendationRenderer) Render(b BlockData, ctx RenderContext) string {
	var c struct {
		Issue    string `json:"issue"`
		Action   string `json:"action"`
		Outcome  string `json:"outcome"`
		Priority string `json:"priority"`
	}
	json.Unmarshal([]byte(b.Recommendation), &c)
	return fmt.Sprintf(`<div class="rec-card"><h3>Recommendation</h3><p>%s</p></div>`,
		esc(c.Action))
}

// ─── Parse block content JSON into BlockData ───

func ParseBlockContent(content string) BlockData {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return BlockData{}
	}
	var b BlockData
	if v, ok := raw["title"].(string); ok { b.FindingTitle = v }
	if v, ok := raw["description"].(string); ok { b.FindingDesc = v }
	if v, ok := raw["text"].(string); ok { b.MessageText = v }
	if v, ok := raw["html"].(string); ok { b.NoteContent = v }
	if v, ok := raw["code"].(string); ok { b.NoteContent = v }
	if v, ok := raw["caption"].(string); ok { b.FindingTitle = v }
	if v, ok := raw["recommendation"].(string); ok { b.Recommendation = v }
	if v, ok := raw["category"].(string); ok { b.Category = v }
	if v, ok := raw["priority"].(string); ok { b.Priority = v }
	if v, ok := raw["severity"].(string); ok { b.Severity = v }
	if v, ok := raw["status"].(string); ok { b.Status = v }
	if v, ok := raw["title"].(string); ok { b.DividerTitle = v }
	if items, ok := raw["items"].([]interface{}); ok {
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				ci := ChecklistItemData{}
				if v, ok := m["text"].(string); ok { ci.Text = v }
				if v, ok := m["checked"].(bool); ok { ci.Checked = v }
				b.ChecklistItems = append(b.ChecklistItems, ci)
			}
		}
	}
	return b
}
