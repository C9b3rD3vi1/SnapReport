import { useEffect, useState } from "react";
import { ArrowLeft, Download, FileText, Loader2, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ConfirmDialog";
import { getReport, deleteReport } from "@/services/api";
import type { Report, ReportInfo } from "@/types/report";
import type { Upload } from "@/types/upload";

interface ReportDetailPageProps {
  reportId: string;
  onBack: () => void;
  onDeleted: () => void;
}

export function ReportDetailPage({ reportId, onBack, onDeleted }: ReportDetailPageProps) {
  const [report, setReport] = useState<Report | null>(null);
  const [uploads, setUploads] = useState<Upload[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleting, setDeleting] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getReport(reportId)
      .then((data) => {
        setReport(data.report);
        setUploads(data.uploads);
      })
      .catch(() => setError("Failed to load report"))
      .finally(() => setLoading(false));
  }, [reportId]);

  const handleDelete = async () => {
    setDeleting(true);
    try {
      await deleteReport(reportId);
      onDeleted();
    } catch {
      setError("Failed to delete report");
    } finally {
      setDeleting(false);
      setShowDelete(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[40vh]">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !report) {
    return (
      <div className="text-center py-16">
        <p className="text-destructive">{error || "Report not found"}</p>
        <Button variant="outline" className="mt-4" onClick={onBack}>Go Back</Button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={onBack}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h2 className="text-2xl font-bold tracking-tight">{report.title}</h2>
          <p className="text-muted-foreground text-sm mt-1">
            {new Date(report.created_at).toLocaleDateString()} &middot; {uploads.length} screenshot{uploads.length !== 1 ? "s" : ""} &middot; {report.status}
          </p>
        </div>
      </div>

      <div className="flex gap-4">
        {report.status === "completed" && (
          <Button asChild>
            <a href={`/api/v1/reports/${report.id}/download`} download>
              <Download className="h-4 w-4 mr-2" />
              Download PDF
            </a>
          </Button>
        )}
        <Button variant="destructive" onClick={() => setShowDelete(true)} disabled={deleting}>
          {deleting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
          <Trash2 className="h-4 w-4 mr-2" />
          Delete Report
        </Button>
      </div>

      {uploads.length > 0 && (
        <div>
          <h3 className="text-lg font-semibold mb-4">Screenshots ({uploads.length})</h3>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
            {uploads.map((u) => (
              <div key={u.id} className="rounded-lg border bg-card overflow-hidden">
                <div className="aspect-video bg-muted">
                  <img
                    src={u.thumbnail_url || `/uploads/${u.filename}`}
                    alt={u.title || u.original_name}
                    className="w-full h-full object-cover"
                  />
                </div>
                <div className="p-2">
                  <p className="text-xs font-medium truncate">{u.title || u.original_name}</p>
                  {u.description && (
                    <p className="text-xs text-muted-foreground truncate mt-1">{u.description}</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <ConfirmDialog
        open={showDelete}
        title="Delete Report"
        message="Are you sure you want to delete this report? The PDF and all associated data will be permanently removed."
        confirmLabel="Delete"
        onConfirm={handleDelete}
        onCancel={() => setShowDelete(false)}
      />
    </div>
  );
}
