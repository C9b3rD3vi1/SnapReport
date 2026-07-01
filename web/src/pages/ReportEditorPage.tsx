import { useCallback, useEffect, useRef, useState } from "react";
import { ArrowLeft, Eye, FileText, Loader2, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ReportInfoForm } from "@/components/ReportInfoForm";
import { SortableImageList } from "@/components/SortableImageList";
import { TemplatePicker } from "@/components/TemplatePicker";
import { listUploads, createReport, previewReport } from "@/services/api";
import type { ReportTemplate } from "@/services/api";
import type { Upload } from "@/types/upload";
import type { ReportInfo, Report } from "@/types/report";

interface ReportEditorPageProps {
  onBack: () => void;
  onViewReport: (id: string) => void;
}

export function ReportEditorPage({ onBack, onViewReport }: ReportEditorPageProps) {
  const [uploads, setUploads] = useState<Upload[]>([]);
  const [reportInfo, setReportInfo] = useState<ReportInfo>({
    title: "", project: "", company: "", author: "", version: "",
  });
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [result, setResult] = useState<Report | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [previewHtml, setPreviewHtml] = useState<string | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const previewTimer = useRef<ReturnType<typeof setTimeout>>();

  useEffect(() => {
    listUploads()
      .then(setUploads)
      .catch(() => setError("Failed to load uploads"))
      .finally(() => setLoading(false));
  }, []);

  const debouncedPreview = useCallback(() => {
    if (previewTimer.current) clearTimeout(previewTimer.current);
    previewTimer.current = setTimeout(async () => {
      const valid = uploads.filter((u) => u.id);
      if (!reportInfo.title || valid.length === 0) return;
      setPreviewLoading(true);
      try {
        const html = await previewReport({
          title: reportInfo.title,
          project: reportInfo.project,
          company: reportInfo.company,
          author: reportInfo.author,
          version: reportInfo.version,
          uploads: valid.map((u, i) => ({
            id: u.id,
            title: u.title,
            description: u.description,
            notes: u.notes ?? "",
            category: u.category ?? "",
            priority: u.priority ?? "",
            severity: u.severity ?? "",
            status: u.status ?? "",
            recommendation: "",
            order_index: i,
          })),
        });
        setPreviewHtml(html);
      } catch {
        // Preview silently fails
      } finally {
        setPreviewLoading(false);
      }
    }, 800);
  }, [reportInfo, uploads]);

  useEffect(() => {
    debouncedPreview();
  }, [debouncedPreview]);

  const handleTemplateSelect = useCallback((t: ReportTemplate) => {
    setReportInfo((prev) => ({
      ...prev,
      title: t.default_title,
    }));
    setUploads((prev) =>
      prev.map((u) => ({
        ...u,
        category: t.categories[0] ?? u.category,
        priority: t.priorities[1] ?? u.priority,
        severity: t.default_severity ?? u.severity,
        status: t.statuses[0] ?? u.status,
      })),
    );
  }, []);

  const handleReorder = useCallback((reordered: Upload[]) => {
    setUploads(reordered.map((u, i) => ({ ...u, order_index: i })));
  }, []);

  const handleUpdateMetadata = useCallback(
    (id: string, title: string, description: string) => {
      setUploads((prev) =>
        prev.map((u) => (u.id === id ? { ...u, title, description } : u)),
      );
    },
    [],
  );

  const handleGenerate = useCallback(async () => {
    if (!reportInfo.title) { setError("Report title is required"); return; }
    if (uploads.length === 0) { setError("At least one upload is required"); return; }

    setGenerating(true);
    setError(null);
    try {
      const report = await createReport({
        ...reportInfo,
        uploads: uploads.map((u, i) => ({
          id: u.id,
          title: u.title,
          description: u.description,
          notes: u.notes ?? "",
          category: u.category ?? "",
          priority: u.priority ?? "",
          severity: u.severity ?? "",
          status: u.status ?? "",
          recommendation: u.recommendation ?? "",
          order_index: i,
        })),
      });
      setResult(report);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create report");
    } finally {
      setGenerating(false);
    }
  }, [reportInfo, uploads]);

  if (loading) {
    return <div className="flex items-center justify-center min-h-[40vh]"><p className="text-muted-foreground">Loading uploads...</p></div>;
  }

  if (result) {
    return (
      <div className="max-w-lg mx-auto text-center space-y-6 py-12">
        <div className="rounded-full bg-primary/10 p-4 w-16 h-16 mx-auto flex items-center justify-center">
          <FileText className="h-8 w-8 text-primary" />
        </div>
        <h2 className="text-2xl font-bold">Report Created</h2>
        <p className="text-muted-foreground">Your report has been created with {uploads.length} screenshot{uploads.length !== 1 ? "s" : ""}.</p>
        <div className="bg-muted rounded-lg p-4 text-left space-y-1 text-sm">
          <p><strong>ID:</strong> {result.id}</p>
          <p><strong>Title:</strong> {result.title}</p>
          <p><strong>Status:</strong> {result.status}</p>
        </div>
        <div className="flex gap-4 justify-center flex-wrap">
          <Button variant="outline" onClick={onBack}>Back to Uploads</Button>
          <Button onClick={() => onViewReport(result.id)}><FileText className="h-4 w-4 mr-2" />View Report</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={onBack}><ArrowLeft className="h-5 w-5" /></Button>
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Edit Report</h2>
          <p className="text-muted-foreground mt-1">Reorder screenshots, add descriptions, and enter report information</p>
        </div>
      </div>

      <TemplatePicker onSelect={handleTemplateSelect} />

      <div className="grid grid-cols-1 xl:grid-cols-5 gap-6">
        <div className="xl:col-span-3 space-y-6">
          <div>
            <h3 className="text-lg font-semibold mb-4">Screenshots ({uploads.length})</h3>
            <SortableImageList
              uploads={uploads}
              onReorder={handleReorder}
              onUpdateMetadata={handleUpdateMetadata}
            />
          </div>
        </div>

        <div className="xl:col-span-2 space-y-6">
          <div className="rounded-lg border bg-card p-6">
            <h3 className="text-lg font-semibold mb-4">Report Information</h3>
            <ReportInfoForm onChange={setReportInfo} />
          </div>

          <div className="rounded-lg border bg-card overflow-hidden">
            <div className="flex items-center justify-between p-3 border-b bg-muted/50">
              <h3 className="text-sm font-semibold flex items-center gap-2">
                <Eye className="h-4 w-4" />
                Preview
              </h3>
              {previewHtml && (
                <Button variant="ghost" size="sm" onClick={() => setPreviewHtml(null)}>
                  <X className="h-3 w-3" />
                </Button>
              )}
            </div>
            {previewLoading && (
              <div className="flex items-center justify-center p-8">
                <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
              </div>
            )}
            {previewHtml && !previewLoading && (
              <iframe
                srcDoc={previewHtml}
                className="w-full h-[500px] border-0"
                title="Report Preview"
                sandbox="allow-same-origin"
              />
            )}
            {!previewHtml && !previewLoading && (
              <div className="text-center p-8 text-sm text-muted-foreground">
                Enter a report title and add screenshots to see a live preview
              </div>
            )}
          </div>

          {error && (
            <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4">
              <p className="text-sm text-destructive">{error}</p>
            </div>
          )}

          <Button className="w-full" size="lg" onClick={handleGenerate} disabled={generating}>
            {generating && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
            {generating ? "Creating Report..." : "Generate Report"}
          </Button>
        </div>
      </div>
    </div>
  );
}
