import { useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { cn } from "@/utils/cn";
import type { ReportInfo } from "@/types/report";

const schema = z.object({
  title: z.string().min(1, "Report title is required"),
  project: z.string().optional().default(""),
  company: z.string().optional().default(""),
  author: z.string().optional().default(""),
  version: z.string().optional().default(""),
});

interface ReportInfoFormProps {
  onChange: (data: ReportInfo) => void;
  className?: string;
}

export function ReportInfoForm({ onChange, className }: ReportInfoFormProps) {
  const {
    register,
    watch,
    formState: { errors },
  } = useForm<ReportInfo>({
    resolver: zodResolver(schema),
    defaultValues: {
      title: "",
      project: "",
      company: "",
      author: "",
      version: "",
    },
  });

  const prevRef = useRef("");

  useEffect(() => {
    const sub = watch((values) => {
      const json = JSON.stringify(values);
      if (json !== prevRef.current) {
        prevRef.current = json;
        onChange(values as ReportInfo);
      }
    });
    return () => sub.unsubscribe();
  }, [watch, onChange]);

  return (
    <div className={cn("space-y-4", className)}>
      <div>
        <label htmlFor="title" className="text-sm font-medium">
          Report Title <span className="text-destructive">*</span>
        </label>
        <input
          id="title"
          {...register("title")}
          className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          placeholder="Enter report title"
        />
        {errors.title && (
          <p className="text-xs text-destructive mt-1">{errors.title.message}</p>
        )}
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label htmlFor="project" className="text-sm font-medium">
            Project Name
          </label>
          <input
            id="project"
            {...register("project")}
            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            placeholder="e.g. SnapReport v2"
          />
        </div>
        <div>
          <label htmlFor="company" className="text-sm font-medium">
            Company
          </label>
          <input
            id="company"
            {...register("company")}
            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            placeholder="e.g. ACME Inc."
          />
        </div>
        <div>
          <label htmlFor="author" className="text-sm font-medium">
            Author
          </label>
          <input
            id="author"
            {...register("author")}
            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            placeholder="Your name"
          />
        </div>
        <div>
          <label htmlFor="version" className="text-sm font-medium">
            Version
          </label>
          <input
            id="version"
            {...register("version")}
            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            placeholder="e.g. 1.0.0"
          />
        </div>
      </div>
    </div>
  );
}
