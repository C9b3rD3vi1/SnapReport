import { Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { cn } from "@/utils/cn";
import type { UploadFile } from "@/types/upload";

interface ImagePreviewProps {
  files: UploadFile[];
  onDelete: (id: string) => void;
  deleting: Set<string>;
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function ImagePreview({ files, onDelete, deleting }: ImagePreviewProps) {
  if (files.length === 0) return null;

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
      {files.map((f) => (
        <div
          key={f.id}
          className={cn(
            "group relative rounded-lg border overflow-hidden bg-card",
            f.error && "border-destructive",
          )}
        >
          <div className="aspect-video relative bg-muted">
            <img
              src={f.preview}
              alt={f.file.name}
              className="w-full h-full object-cover"
            />
            {f.upload && (
              <Button
                variant="destructive"
                size="icon"
                className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity h-8 w-8"
                onClick={() => onDelete(f.id)}
                disabled={deleting.has(f.id)}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            )}
          </div>
          <div className="p-2 space-y-1">
            <p className="text-xs font-medium truncate">{f.file.name}</p>
            <p className="text-xs text-muted-foreground">{formatSize(f.file.size)}</p>
            {f.progress < 100 && !f.error && (
              <Progress value={f.progress} className="h-1" />
            )}
            {f.error && (
              <p className="text-xs text-destructive truncate">{f.error}</p>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
