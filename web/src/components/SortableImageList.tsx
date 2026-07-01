import {
  DndContext,
  type DragEndEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  useSortable,
  rectSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical } from "lucide-react";
import { cn } from "@/utils/cn";
import { ImageMetadataEditor } from "./ImageMetadataEditor";
import type { Upload } from "@/types/upload";

interface SortableImageListProps {
  uploads: Upload[];
  onReorder: (uploads: Upload[]) => void;
  onUpdateMetadata: (id: string, title: string, description: string) => void;
}

interface SortableItemProps {
  upload: Upload;
  onUpdateMetadata: (id: string, title: string, description: string) => void;
}

function SortableItem({ upload, onUpdateMetadata }: SortableItemProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: upload.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        "rounded-lg border bg-card overflow-hidden",
        isDragging && "opacity-50 shadow-lg",
      )}
    >
      <div className="relative aspect-video bg-muted">
        <img
          src={`/uploads/${upload.filename}`}
          alt={upload.original_name}
          className="w-full h-full object-cover"
        />
        <button
          type="button"
          className="absolute top-2 left-2 h-7 w-7 flex items-center justify-center rounded-md bg-background/80 text-muted-foreground hover:text-foreground cursor-grab active:cursor-grabbing transition-colors"
          {...attributes}
          {...listeners}
        >
          <GripVertical className="h-4 w-4" />
        </button>
      </div>
      <div className="p-3">
        <p className="text-xs font-medium truncate">{upload.original_name}</p>
        <ImageMetadataEditor
          title={upload.title}
          description={upload.description}
          onTitleChange={(title) => onUpdateMetadata(upload.id, title, upload.description)}
          onDescriptionChange={(description) => onUpdateMetadata(upload.id, upload.title, description)}
        />
      </div>
    </div>
  );
}

export function SortableImageList({
  uploads,
  onReorder,
  onUpdateMetadata,
}: SortableImageListProps) {
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 8 },
    }),
  );

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const oldIndex = uploads.findIndex((u) => u.id === active.id);
    const newIndex = uploads.findIndex((u) => u.id === over.id);

    const reordered = [...uploads];
    const [removed] = reordered.splice(oldIndex, 1);
    reordered.splice(newIndex, 0, removed);

    onReorder(reordered);
  }

  return (
    <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
      <SortableContext items={uploads.map((u) => u.id)} strategy={rectSortingStrategy}>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {uploads.map((u) => (
            <SortableItem
              key={u.id}
              upload={u}
              onUpdateMetadata={onUpdateMetadata}
            />
          ))}
        </div>
      </SortableContext>
    </DndContext>
  );
}
