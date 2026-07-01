export type BlockType =
  | "finding"
  | "note"
  | "warning"
  | "tip"
  | "important"
  | "divider"
  | "checklist"
  | "statistics";

export interface Block {
  id: string;
  report_id: string;
  type: BlockType;
  position: number;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface FindingContent {
  screenshot_id?: string;
  title?: string;
  description?: string;
  category?: string;
  priority?: string;
  severity?: string;
  status?: string;
  recommendation?: string;
  notes?: string;
}

export interface NoteContent {
  html: string;
}

export interface WarningContent {
  text: string;
}

export interface TipContent {
  text: string;
}

export interface ImportantContent {
  text: string;
}

export interface DividerContent {
  title: string;
}

export interface ChecklistItem {
  text: string;
  checked: boolean;
}

export interface ChecklistContent {
  items: ChecklistItem[];
}
