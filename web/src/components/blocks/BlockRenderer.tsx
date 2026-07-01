import { GripVertical, Plus, Trash2, AlertTriangle, Lightbulb, Info, Minus, CheckSquare, FileText, Type } from "lucide-react";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { cn } from "@/utils/cn";
import type { Block, BlockType, FindingContent, NoteContent, WarningContent, TipContent, ImportantContent, DividerContent, ChecklistContent, ChecklistItem } from "@/types/block";

interface BlockRendererProps {
  block: Block;
  uploads: { id: string; filename: string; thumbnail_url?: string }[];
  onUpdate: (id: string, content: string) => void;
  onDelete: (id: string) => void;
}

const blockIcons: Record<BlockType, React.ReactNode> = {
  finding: <FileText className="h-4 w-4" />,
  note: <Type className="h-4 w-4" />,
  warning: <AlertTriangle className="h-4 w-4" />,
  tip: <Lightbulb className="h-4 w-4" />,
  important: <Info className="h-4 w-4" />,
  divider: <Minus className="h-4 w-4" />,
  checklist: <CheckSquare className="h-4 w-4" />,
  statistics: <FileText className="h-4 w-4" />,
};

export function BlockRenderer({ block, uploads, onUpdate, onDelete }: BlockRendererProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: block.id });

  const style = { transform: CSS.Transform.toString(transform), transition };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        "rounded-lg border bg-card overflow-hidden",
        isDragging && "opacity-50 shadow-lg",
        block.type === "warning" && "border-l-4 border-l-amber-500",
        block.type === "tip" && "border-l-4 border-l-blue-500",
        block.type === "important" && "border-l-4 border-l-red-500",
      )}
    >
      <div className="flex items-center gap-2 px-3 py-1.5 border-b bg-muted/30 text-xs text-muted-foreground">
        <button {...attributes} {...listeners} className="cursor-grab hover:text-foreground" aria-label="Drag to reorder">
          <GripVertical className="h-3.5 w-3.5" />
        </button>
        <span className="flex items-center gap-1 font-medium">
          {blockIcons[block.type]}
          {block.type.charAt(0).toUpperCase() + block.type.slice(1)}
        </span>
        <div className="ml-auto flex gap-1">
          <button onClick={() => onDelete(block.id)} className="hover:text-destructive" aria-label="Delete block">
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
      <div className="p-4">
        <BlockContent block={block} uploads={uploads} onUpdate={onUpdate} />
      </div>
    </div>
  );
}

function BlockContent({ block, uploads, onUpdate }: BlockRendererProps) {
  switch (block.type) {
    case "finding": return <FindingBlock block={block} uploads={uploads} onUpdate={onUpdate} />;
    case "note": return <NoteBlock block={block} onUpdate={onUpdate} />;
    case "warning": return <MessageBlock block={block} variant="warning" onUpdate={onUpdate} />;
    case "tip": return <MessageBlock block={block} variant="tip" onUpdate={onUpdate} />;
    case "important": return <MessageBlock block={block} variant="important" onUpdate={onUpdate} />;
    case "divider": return <DividerBlock block={block} onUpdate={onUpdate} />;
    case "checklist": return <ChecklistBlock block={block} onUpdate={onUpdate} />;
    case "statistics": return <StatisticsBlock block={block} />;
    default: return <div className="text-sm text-muted-foreground">Unknown block type</div>;
  }
}

function parseContent<T>(block: Block, fallback: T): T {
  try { return JSON.parse(block.content) as T; } catch { return fallback; }
}

function FindingBlock({ block, uploads, onUpdate }: BlockRendererProps) {
  const content = parseContent<FindingContent>(block, {});
  const upload = uploads.find((u) => u.id === content.screenshot_id);

  const update = (partial: Partial<FindingContent>) => {
    onUpdate(block.id, JSON.stringify({ ...content, ...partial }));
  };

  return (
    <div className="space-y-3">
      {upload && (
        <img
          src={upload.thumbnail_url || `/uploads/${upload.filename}`}
          alt={content.title || ""}
          className="w-full max-h-40 object-contain rounded border bg-muted"
        />
      )}
      <input
        value={content.title || ""}
        onChange={(e) => update({ title: e.target.value })}
        placeholder="Finding title"
        className="w-full text-sm font-medium border-0 border-b border-input bg-transparent pb-1 focus:outline-none focus:border-primary"
      />
      <textarea
        value={content.description || ""}
        onChange={(e) => update({ description: e.target.value })}
        placeholder="Description"
        rows={2}
        className="w-full text-xs border-0 bg-transparent resize-none focus:outline-none text-muted-foreground"
      />
      <div className="flex flex-wrap gap-2">
        <select value={content.priority || ""} onChange={(e) => update({ priority: e.target.value })} className="text-xs border rounded px-2 py-1 bg-background">
          <option value="">Priority</option>
          <option>Critical</option><option>High</option><option>Medium</option><option>Low</option>
        </select>
        <select value={content.status || ""} onChange={(e) => update({ status: e.target.value })} className="text-xs border rounded px-2 py-1 bg-background">
          <option value="">Status</option>
          <option>Open</option><option>Resolved</option><option>Closed</option>
        </select>
      </div>
    </div>
  );
}

function NoteBlock({ block, onUpdate }: { block: Block; onUpdate: (id: string, content: string) => void }) {
  const content = parseContent<NoteContent>(block, { html: "" });
  return (
    <textarea
      value={content.html}
      onChange={(e) => onUpdate(block.id, JSON.stringify({ html: e.target.value }))}
      placeholder="Write your notes here..."
      rows={4}
      className="w-full text-sm border-0 bg-transparent resize-y focus:outline-none"
    />
  );
}

function MessageBlock({ block, variant }: { block: Block; onUpdate: (id: string, content: string) => void; variant: "warning" | "tip" | "important" }) {
  const content = parseContent<WarningContent>(block, { text: "" });
  const icons = { warning: <AlertTriangle className="h-5 w-5 text-amber-500" />, tip: <Lightbulb className="h-5 w-5 text-blue-500" />, important: <Info className="h-5 w-5 text-red-500" /> };
  const colors = { warning: "text-amber-800 bg-amber-50", tip: "text-blue-800 bg-blue-50", important: "text-red-800 bg-red-50" };
  return (
    <div className={cn("flex gap-3 p-3 rounded-lg text-sm", colors[variant])}>
      <div className="mt-0.5">{icons[variant]}</div>
      <input
        value={content.text}
        onChange={(e) => onUpdate(block.id, JSON.stringify({ text: e.target.value }))}
        placeholder={variant === "warning" ? "Enter warning text..." : variant === "tip" ? "Enter tip..." : "Enter important notice..."}
        className="flex-1 bg-transparent border-0 focus:outline-none placeholder-opacity-50"
      />
    </div>
  );
}

function DividerBlock({ block, onUpdate }: { block: Block; onUpdate: (id: string, content: string) => void }) {
  const content = parseContent<DividerContent>(block, { title: "" });
  return (
    <div className="flex items-center gap-3">
      <div className="flex-1 h-px bg-border" />
      <input
        value={content.title}
        onChange={(e) => onUpdate(block.id, JSON.stringify({ title: e.target.value }))}
        placeholder="Section title (optional)"
        className="text-xs font-medium text-center border-0 bg-transparent focus:outline-none text-muted-foreground"
      />
      <div className="flex-1 h-px bg-border" />
    </div>
  );
}

function ChecklistBlock({ block, onUpdate }: { block: Block; onUpdate: (id: string, content: string) => void }) {
  const content = parseContent<ChecklistContent>(block, { items: [] });

  const toggle = (idx: number) => {
    const items = content.items.map((item, i) => i === idx ? { ...item, checked: !item.checked } : item);
    onUpdate(block.id, JSON.stringify({ items }));
  };

  const addItem = () => {
    const items = [...content.items, { text: "", checked: false }];
    onUpdate(block.id, JSON.stringify({ items }));
  };

  const updateText = (idx: number, text: string) => {
    const items = content.items.map((item, i) => i === idx ? { ...item, text } : item);
    onUpdate(block.id, JSON.stringify({ items }));
  };

  return (
    <div className="space-y-1">
      {content.items.map((item, i) => (
        <label key={i} className="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" checked={item.checked} onChange={() => toggle(i)} className="rounded" />
          <input
            value={item.text}
            onChange={(e) => updateText(i, e.target.value)}
            placeholder={`Item ${i + 1}`}
            className={cn("flex-1 border-0 bg-transparent text-sm focus:outline-none", item.checked && "line-through text-muted-foreground")}
          />
        </label>
      ))}
      <button onClick={addItem} className="text-xs text-muted-foreground hover:text-foreground mt-1">+ Add item</button>
    </div>
  );
}

function StatisticsBlock({ block }: { block: Block }) {
  return (
    <div className="text-center py-4 text-sm text-muted-foreground">
      Report statistics will be auto-generated here
    </div>
  );
}
