package models

import "time"

type BlockType string

const (
	BlockFinding    BlockType = "finding"
	BlockNote       BlockType = "note"
	BlockWarning    BlockType = "warning"
	BlockTip        BlockType = "tip"
	BlockImportant  BlockType = "important"
	BlockDivider    BlockType = "divider"
	BlockChecklist  BlockType = "checklist"
	BlockStatistics BlockType = "statistics"
)

type Block struct {
	ID        string    `json:"id"`
	ReportID  string    `json:"report_id"`
	Type      BlockType `json:"type"`
	Position  int       `json:"position"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FindingContent struct {
	ScreenshotID   string `json:"screenshot_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	Priority       string `json:"priority"`
	Severity       string `json:"severity"`
	Status         string `json:"status"`
	Recommendation string `json:"recommendation"`
	Notes          string `json:"notes"`
}

type NoteContent struct {
	HTML string `json:"html"`
}

type WarningContent struct {
	Text string `json:"text"`
}

type TipContent struct {
	Text string `json:"text"`
}

type ImportantContent struct {
	Text string `json:"text"`
}

type DividerContent struct {
	Title string `json:"title"`
}

type ChecklistItem struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
}

type ChecklistContent struct {
	Items []ChecklistItem `json:"items"`
}
