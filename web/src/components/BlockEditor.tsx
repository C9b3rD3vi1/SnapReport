import { useCallback, useState } from "react";
import {
  DndContext, type DragEndEvent, PointerSensor, useSensor, useSensors,
} from "@dnd-kit/core";
import {
  SortableContext, rectSortingStrategy, arrayMove,
} from "@dnd-kit/sortable";
import { Plus, FileText, Type, AlertTriangle, Lightbulb, Info, Minus, CheckSquare } from "lucide-react";
import { BlockRenderer } from "@/components/blocks/BlockRenderer";
import type { Block, BlockType } from "@/types/block";
import { cn } from "@/utils/cn";

interface BlockEditorProps {
  blocks: Block[];
  uploads: { id: string; filename: string; thumbnail_url?: string }[];
  onReorder: (blocks: Block[]) => void;
  onUpdate: (id: string, content: string) => void;
  onDelete: (id: string) => void;
  onAdd: (type: BlockType) => void;
}

const blockTypes: { type: BlockType; label: string; icon: React.ReactNode; description: string }[] = [
  { type: "finding", label: "Finding", icon: <FileText className="h-4 w-4" />, description: "Screenshot with metadata" },
  { type: "note", label: "Note", icon: <Type className="h-4 w-4" />, description: "Rich text section" },
  { type: "warning", label: "Warning", icon: <AlertTriangle className="h-4 w-4" />, description: "Important warning" },
  { type: "tip", label: "Tip", icon: <Lightbulb className="h-4 w-4" />, description: "Best practice" },
  { type: "important", label: "Important", icon: <Info className="h-4 w-4" />, description: "Critical notice" },
  { type: "divider", label: "Divider", icon: <Minus className="h-4 w-4" />, description: "Section separator" },
  { type: "checklist", label: "Checklist", icon: <CheckSquare className="h-4 w-4" />, description: "Checkable items" },
  { type: "statistics", label: "Statistics", icon: <FileText className="h-4 w-4" />, description: "Auto-generated stats" },
];

export function BlockEditor({ blocks, uploads, onReorder, onUpdate, onDelete, onAdd }: BlockEditorProps) {
  const [showAdd, setShowAdd] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
  );

  const handleDragEnd = useCallback((event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const oldIndex = blocks.findIndex((b) => b.id === active.id);
    const newIndex = blocks.findIndex((b) => b.id === over.id);
    onReorder(arrayMove(blocks, oldIndex, newIndex));
  }, [blocks, onReorder]);

  return (
    <div className="space-y-3">
      <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
        <SortableContext items={blocks.map((b) => b.id)} strategy={rectSortingStrategy}>
          {blocks.map((block) => (
            <BlockRenderer
              key={block.id}
              block={block}
              uploads={uploads}
              onUpdate={onUpdate}
              onDelete={onDelete}
            />
          ))}
        </SortableContext>
      </DndContext>

      <div className="relative">
        <button
          onClick={() => setShowAdd(!showAdd)}
          className="w-full py-3 border-2 border-dashed border-border rounded-lg text-sm text-muted-foreground hover:border-primary hover:text-primary transition-colors flex items-center justify-center gap-2"
        >
          <Plus className="h-4 w-4" />
          Add Block
        </button>

        {showAdd && (
          <div className="mt-2 grid grid-cols-2 sm:grid-cols-4 gap-2 p-3 rounded-lg border bg-card shadow-lg">
            {blockTypes.map((bt) => (
              <button
                key={bt.type}
                onClick={() => { onAdd(bt.type); setShowAdd(false); }}
                className={cn(
                  "flex items-center gap-2 p-2 rounded-md text-xs text-left hover:bg-accent transition-colors",
                  "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
                )}
              >
                {bt.icon}
                <div>
                  <p className="font-medium">{bt.label}</p>
                  <p className="text-[10px] text-muted-foreground">{bt.description}</p>
                </div>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
