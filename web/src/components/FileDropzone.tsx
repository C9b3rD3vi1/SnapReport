import { useCallback, useRef, useState } from "react";
import { Upload as UploadIcon } from "lucide-react";
import { cn } from "@/utils/cn";

interface FileDropzoneProps {
  onFilesSelected: (files: File[]) => void;
  disabled?: boolean;
}

export function FileDropzone({ onFilesSelected, disabled }: FileDropzoneProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = useState(false);

  const handleFiles = useCallback(
    (files: FileList | File[]) => {
      const valid = Array.from(files).filter((f) =>
        ["image/png", "image/jpeg", "image/webp"].includes(f.type),
      );
      if (valid.length > 0) onFilesSelected(valid);
    },
    [onFilesSelected],
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setDragging(false);
      if (disabled) return;
      handleFiles(e.dataTransfer.files);
    },
    [disabled, handleFiles],
  );

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files) handleFiles(e.target.files);
      if (inputRef.current) inputRef.current.value = "";
    },
    [handleFiles],
  );

  const handlePaste = useCallback(
    (e: React.ClipboardEvent) => {
      const items = e.clipboardData?.items;
      if (!items) return;
      const files: File[] = [];
      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.kind === "file" && item.type.startsWith("image/")) {
          files.push(item.getAsFile()!);
        }
      }
      if (files.length > 0) {
        e.preventDefault();
        onFilesSelected(files);
      }
    },
    [onFilesSelected],
  );

  return (
    <div
      onDragOver={(e) => { e.preventDefault(); setDragging(true); }}
      onDragLeave={() => setDragging(false)}
      onDrop={handleDrop}
      onPaste={handlePaste}
      onClick={() => inputRef.current?.click()}
      onKeyDown={(e) => { if (e.key === "Enter" || e.key === " ") inputRef.current?.click(); }}
      role="button"
      tabIndex={0}
      aria-label="Upload screenshots. Drop files, click to browse, or paste from clipboard."
      className={cn(
        "border-2 border-dashed rounded-xl p-12 text-center cursor-pointer transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
        dragging ? "border-primary bg-primary/5" : "border-muted-foreground/25 hover:border-muted-foreground/50",
        disabled && "opacity-50 pointer-events-none",
      )}
    >
      <input
        ref={inputRef}
        type="file"
        accept="image/png,image/jpeg,image/webp"
        multiple
        className="hidden"
        onChange={handleChange}
        disabled={disabled}
        aria-hidden="true"
      />
      <UploadIcon className="mx-auto h-12 w-12 text-muted-foreground/50" />
      <p className="mt-4 text-sm font-medium">
        Drop screenshots here or click to browse
      </p>
      <p className="mt-1 text-xs text-muted-foreground">
        PNG, JPG, WEBP &middot; Max 10 MB each &middot; Paste from clipboard
      </p>
    </div>
  );
}
