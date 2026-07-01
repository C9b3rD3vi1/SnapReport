import { useCallback, useEffect, useRef, useState } from "react";
import { FileDropzone } from "@/components/FileDropzone";
import { ImagePreview } from "@/components/ImagePreview";
import { Button } from "@/components/ui/button";
import { uploadFiles, listUploads, deleteUpload, createReport } from "@/services/api";
import type { UploadFile, Upload } from "@/types/upload";
import { Loader2 } from "lucide-react";

interface UploadPageProps {
  onContinue?: () => void;
}

function autoTitle(filename: string): string {
  const withoutExt = filename.replace(/\.[^.]+$/, "");
  return withoutExt
    .replace(/[-_]+/g, " ")
    .replace(/\b\w/g, (c) => c.toUpperCase());
}

export function UploadPage({ onContinue }: UploadPageProps) {
  const [files, setFiles] = useState<UploadFile[]>([]);
  const [deleting, setDeleting] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const counterRef = useRef(0);

  useEffect(() => {
    listUploads()
      .then((uploads) => {
        setFiles(
          uploads.map((u) => ({
            id: u.id,
            file: new File([], u.original_name),
            preview: u.thumbnail_url || `/uploads/${u.filename}`,
            progress: 100,
            upload: u,
            error: null,
          })),
        );
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const handleFilesSelected = useCallback(
    async (newFiles: File[]) => {
      const mapped: UploadFile[] = newFiles.map((file) => ({
        id: `pending-${counterRef.current++}`,
        file,
        preview: URL.createObjectURL(file),
        progress: 0,
        upload: null,
        error: null,
      }));

      setFiles((prev) => [...prev, ...mapped]);

      try {
        const uploads = await uploadFiles(newFiles);
        setFiles((prev) =>
          prev.map((f) => {
            const match = uploads.find(
              (u) => u.original_name === f.file.name && u.size === f.file.size,
            );
            if (match) {
              URL.revokeObjectURL(f.preview);
              return {
                ...f,
                id: match.id,
                preview: match.thumbnail_url || `/uploads/${match.filename}`,
                progress: 100,
                upload: { ...match, title: match.title || autoTitle(match.original_name) },
              };
            }
            return f;
          }),
        );
      } catch (err) {
        setFiles((prev) =>
          prev.map((f) =>
            f.upload === null
              ? { ...f, error: err instanceof Error ? err.message : "Upload failed" }
              : f,
          ),
        );
      }
    },
    [],
  );

  const handleDelete = useCallback(async (id: string) => {
    setDeleting((prev) => new Set(prev).add(id));
    try {
      await deleteUpload(id);
      setFiles((prev) => prev.filter((f) => f.id !== id));
    } catch {
      // Error handled silently
    } finally {
      setDeleting((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  }, []);

  const handleQuickGenerate = useCallback(async () => {
    setGenerating(true);
    try {
      const report = await createReport({
        title: "Quick Report",
        project: "",
        company: "",
        author: "",
        version: "",
        uploads: files
          .filter((f) => f.upload)
          .map((f, i) => ({
            id: f.upload!.id,
            title: autoTitle(f.upload!.original_name),
            description: "",
            notes: "",
            order_index: i,
          })),
      });
      window.open(`/api/v1/reports/${report.id}/download`, "_blank");
    } catch {
      // Silently handle
    } finally {
      setGenerating(false);
    }
  }, [files]);

  const hasUploads = files.some((f) => f.upload);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[40vh]">
        <p className="text-muted-foreground">Loading uploads...</p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Upload Screenshots</h2>
        <p className="text-muted-foreground mt-1">
          Drag and drop your screenshots or click to browse
        </p>
      </div>

      <FileDropzone onFilesSelected={handleFilesSelected} disabled={files.some(f => f.progress > 0 && f.progress < 100)} />

      {files.length === 0 && !loading && (
        <div className="text-center py-12">
          <p className="text-muted-foreground">No screenshots uploaded yet</p>
        </div>
      )}

      <ImagePreview files={files} onDelete={handleDelete} deleting={deleting} />

      {hasUploads && (
        <div className="flex justify-end gap-4">
          <Button variant="outline" onClick={handleQuickGenerate} disabled={generating}>
            {generating && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
            Generate Quick PDF
          </Button>
          <Button onClick={onContinue}>
            Continue to Report
          </Button>
        </div>
      )}
    </div>
  );
}
