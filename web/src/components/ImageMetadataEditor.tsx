import { ChevronDown, ChevronRight } from "lucide-react";
import { useState } from "react";

interface ImageMetadataEditorProps {
  title: string;
  description: string;
  onTitleChange: (title: string) => void;
  onDescriptionChange: (description: string) => void;
}

export function ImageMetadataEditor({
  title,
  description,
  onTitleChange,
  onDescriptionChange,
}: ImageMetadataEditorProps) {
  const [open, setOpen] = useState(false);

  return (
    <div className="border-t pt-2 mt-2">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors w-full text-left"
      >
        {open ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
        {open ? "Hide details" : "Add title & description"}
      </button>
      {open && (
        <div className="space-y-2 mt-2">
          <input
            value={title}
            onChange={(e) => onTitleChange(e.target.value)}
            placeholder="Screenshot title"
            className="flex h-8 w-full rounded-md border border-input bg-transparent px-2 py-1 text-xs shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          />
          <textarea
            value={description}
            onChange={(e) => onDescriptionChange(e.target.value)}
            placeholder="Description of this screenshot"
            rows={2}
            className="flex w-full rounded-md border border-input bg-transparent px-2 py-1 text-xs shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
          />
        </div>
      )}
    </div>
  );
}
