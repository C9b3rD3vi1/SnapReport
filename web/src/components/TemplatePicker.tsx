import { useEffect, useState } from "react";
import { FileText, Layout } from "lucide-react";
import { Button } from "@/components/ui/button";
import { listTemplates } from "@/services/api";
import type { ReportTemplate } from "@/services/api";
import { cn } from "@/utils/cn";

interface TemplatePickerProps {
  onSelect: (template: ReportTemplate) => void;
}

export function TemplatePicker({ onSelect }: TemplatePickerProps) {
  const [templates, setTemplates] = useState<ReportTemplate[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listTemplates()
      .then(setTemplates)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="text-center py-8">
        <p className="text-sm text-muted-foreground">Loading templates...</p>
      </div>
    );
  }

  if (templates.length === 0) return null;

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-semibold flex items-center gap-2">
        <Layout className="h-4 w-4" />
        Report Templates
      </h3>
      <p className="text-xs text-muted-foreground mb-3">
        Choose a template to pre-configure categories, priorities, and defaults
      </p>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {templates.map((t) => (
          <button
            key={t.id}
            onClick={() => onSelect(t)}
            className={cn(
              "text-left p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors",
              "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
            )}
          >
            <div className="flex items-start gap-3">
              <FileText className="h-5 w-5 text-muted-foreground shrink-0 mt-0.5" />
              <div>
                <p className="text-sm font-medium">{t.name}</p>
                <p className="text-xs text-muted-foreground mt-0.5 line-clamp-2">{t.description}</p>
                <div className="flex gap-1.5 mt-2 flex-wrap">
                  <span className="text-[10px] px-1.5 py-0.5 rounded bg-muted text-muted-foreground">
                    {t.categories.length} categories
                  </span>
                  <span className="text-[10px] px-1.5 py-0.5 rounded bg-muted text-muted-foreground">
                    {t.classification}
                  </span>
                </div>
              </div>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}
