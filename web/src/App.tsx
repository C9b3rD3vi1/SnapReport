import { useState } from "react";
import { UploadPage } from "@/pages/UploadPage";
import { ReportEditorPage } from "@/pages/ReportEditorPage";

type Page = "upload" | "editor";

function App() {
  const [page, setPage] = useState<Page>("upload");

  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b">
        <div className="container flex h-16 items-center px-4">
          <h1 className="text-xl font-semibold">SnapReport</h1>
          <span className="text-sm text-muted-foreground ml-2">
            Screenshot to PDF
          </span>
        </div>
      </header>
      <main className="flex-1 container px-4 py-8">
        {page === "upload" && (
          <UploadPage onContinue={() => setPage("editor")} />
        )}
        {page === "editor" && (
          <ReportEditorPage onBack={() => setPage("upload")} />
        )}
      </main>
      <footer className="border-t py-4">
        <div className="container text-center text-sm text-muted-foreground px-4">
          SnapReport &mdash; Screenshot to PDF Report Generator
        </div>
      </footer>
    </div>
  );
}

export default App;
