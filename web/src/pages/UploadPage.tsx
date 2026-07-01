import { useCallback, useEffect, useRef, useState } from "react";
import { FileDropzone } from "@/components/FileDropzone";
import { ImagePreview } from "@/components/ImagePreview";
import { Button } from "@/components/ui/button";
import { uploadFiles, listUploads, deleteUpload } from "@/services/api";
import type { UploadFile, Upload } from "@/types/upload";

interface UploadPageProps {
  onContinue?: () => void;
}

export function UploadPage({ onContinue }: UploadPageProps) {
  const [files, setFiles] = useState<UploadFile[]>([]);
  const [deleting, setDeleting] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);
  const counterRef = useRef(0);

  useEffect(() => {
    listUploads()
      .then((uploads) => {
        setFiles(
          uploads.map((u) => ({
            id: u.id,
            file: new File([], u.original_name),
            preview: `/uploads/${u.filename}`,
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
                preview: `/uploads/${match.filename}`,
                progress: 100,
                upload: match,
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

      <FileDropzone onFilesSelected={handleFilesSelected} disabled={files.some(f => f.progress < 100 && f.progress > 0)} />

      {files.length === 0 && !loading && (
        <div className="text-center py-12">
          <p className="text-muted-foreground">No screenshots uploaded yet</p>
        </div>
      )}

      <ImagePreview files={files} onDelete={handleDelete} deleting={deleting} />

      {files.length > 0 && (
        <div className="flex justify-end">
          <Button onClick={onContinue}>
            Continue to Report
          </Button>
        </div>
      )}
    </div>
  );
}
