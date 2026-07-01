import { useCallback, useEffect, useState } from "react";
import { FileText, Download, Trash2, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ConfirmDialog";
import { listReports, deleteReport } from "@/services/api";
import type { Report } from "@/types/report";

interface ReportHistoryPageProps {
  onViewReport: (id: string) => void;
  onBack: () => void;
}

export function ReportHistoryPage({ onViewReport, onBack }: ReportHistoryPageProps) {
  const [reports, setReports] = useState<Report[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchReports = useCallback(async () => {
    setLoading(true);
    try {
      const data = await listReports();
      setReports(data);
    } catch {
      setError("Failed to load reports");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  const handleDelete = useCallback(async () => {
    if (!deleteId) return;
    try {
      await deleteReport(deleteId);
      setReports((prev) => prev.filter((r) => r.id !== deleteId));
    } catch {
      setError("Failed to delete report");
    } finally {
      setDeleteId(null);
    }
  }, [deleteId]);

  const filtered = reports.filter(
    (r) =>
      !search || r.title.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Report History</h2>
          <p className="text-muted-foreground mt-1">
            View and download previously generated reports
          </p>
        </div>
        <Button variant="outline" onClick={onBack}>
          Back to Uploads
        </Button>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search reports by title..."
          className="flex h-9 w-full rounded-md border border-input bg-transparent pl-9 pr-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          aria-label="Search reports"
        />
      </div>

      {loading && (
        <div className="flex items-center justify-center min-h-[20vh]">
          <p className="text-muted-foreground">Loading reports...</p>
        </div>
      )}

      {error && (
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4">
          <p className="text-sm text-destructive">{error}</p>
        </div>
      )}

      {!loading && !error && filtered.length === 0 && (
        <div className="text-center py-16">
          <FileText className="mx-auto h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">
            {search ? "No reports match your search" : "No reports yet"}
          </p>
          {!search && (
            <p className="text-xs text-muted-foreground mt-1">
              Upload screenshots and generate a report to get started
            </p>
          )}
        </div>
      )}

      {!loading && filtered.length > 0 && (
        <div className="space-y-3">
          {filtered.map((r) => (
            <div
              key={r.id}
              className="flex items-center gap-4 rounded-lg border bg-card p-4"
            >
              <FileText className="h-8 w-8 text-muted-foreground shrink-0" />
              <div className="flex-1 min-w-0">
                <p className="font-medium truncate">{r.title}</p>
                <p className="text-xs text-muted-foreground">
                  {new Date(r.created_at).toLocaleDateString()} &middot; {r.status}
                </p>
              </div>
              <div className="flex gap-2 shrink-0">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onViewReport(r.id)}
                  aria-label={`View report ${r.title}`}
                >
                  View
                </Button>
                {r.status === "completed" && (
                  <Button
                    variant="ghost"
                    size="sm"
                    asChild
                  >
                    <a href={`/api/v1/reports/${r.id}/download`} download aria-label={`Download ${r.title}`}>
                      <Download className="h-4 w-4" />
                    </a>
                  </Button>
                )}
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setDeleteId(r.id)}
                  aria-label={`Delete report ${r.title}`}
                >
                  <Trash2 className="h-4 w-4 text-destructive" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}

      <ConfirmDialog
        open={deleteId !== null}
        title="Delete Report"
        message="Are you sure you want to delete this report? This action cannot be undone."
        confirmLabel="Delete"
        onConfirm={handleDelete}
        onCancel={() => setDeleteId(null)}
      />
    </div>
  );
}
