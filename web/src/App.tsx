import { useState } from "react";
import { UploadPage } from "@/pages/UploadPage";
import { ReportEditorPage } from "@/pages/ReportEditorPage";
import { ReportHistoryPage } from "@/pages/ReportHistoryPage";
import { ReportDetailPage } from "@/pages/ReportDetailPage";
import { ToastProvider } from "@/components/ui/toast";
import { ErrorBoundary } from "@/components/ErrorBoundary";

type Page = "upload" | "editor" | "history" | "report-detail";

function App() {
  const [page, setPage] = useState<Page>("upload");
  const [reportId, setReportId] = useState<string | null>(null);

  return (
    <ErrorBoundary>
      <ToastProvider>
        <div className="min-h-screen flex flex-col">
          <header className="border-b">
            <div className="container flex h-16 items-center px-4">
              <h1
                className="text-xl font-semibold cursor-pointer"
                onClick={() => setPage("upload")}
              >
                SnapReport
              </h1>
              <span className="text-sm text-muted-foreground ml-2">
                Technical Report Builder
              </span>
              <nav className="ml-auto flex gap-4 text-sm">
                <button
                  onClick={() => setPage("upload")}
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  Upload
                </button>
                <button
                  onClick={() => setPage("history")}
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  Reports
                </button>
              </nav>
            </div>
          </header>
          <main className="flex-1 container px-4 py-8">
            {page === "upload" && (
              <UploadPage onContinue={() => setPage("editor")} />
            )}
            {page === "editor" && (
              <ReportEditorPage
                onBack={() => setPage("upload")}
                onViewReport={(id) => { setReportId(id); setPage("report-detail"); }}
              />
            )}
            {page === "history" && (
              <ReportHistoryPage
                onViewReport={(id) => { setReportId(id); setPage("report-detail"); }}
                onBack={() => setPage("upload")}
              />
            )}
            {page === "report-detail" && reportId && (
              <ReportDetailPage
                reportId={reportId}
                onBack={() => setPage("history")}
                onDeleted={() => setPage("history")}
              />
            )}
          </main>
          <footer className="border-t py-4">
            <div className="container text-center text-sm text-muted-foreground px-4">
              SnapReport &mdash; Technical Report Builder
            </div>
          </footer>
        </div>
      </ToastProvider>
    </ErrorBoundary>
  );
}

export default App;
