import axios from "axios";
import type { ApiResponse } from "@/types/api";
import type { Upload } from "@/types/upload";
import type { Report, CreateReportRequest } from "@/types/report";

const client = axios.create({
  baseURL: "/api/v1",
  headers: {
    Accept: "application/json",
  },
});

export async function uploadFiles(
  files: File[],
  onProgress?: (percent: number) => void,
): Promise<Upload[]> {
  const form = new FormData();
  files.forEach((f) => form.append("files", f));

  const res = await client.post<ApiResponse<Upload[]>>("/uploads", form, {
    headers: { "Content-Type": "multipart/form-data" },
    onUploadProgress: (e) => {
      if (e.total && onProgress) {
        onProgress(Math.round((e.loaded * 100) / e.total));
      }
    },
  });

  if (!res.data.success || !res.data.data) {
    throw new Error(res.data.error ?? "Upload failed");
  }

  return res.data.data;
}

export async function listUploads(): Promise<Upload[]> {
  const res = await client.get<ApiResponse<Upload[]>>("/uploads");
  if (!res.data.success || !res.data.data) {
    throw new Error(res.data.error ?? "Failed to list uploads");
  }
  return res.data.data;
}

export async function deleteUpload(id: string): Promise<void> {
  const res = await client.delete<ApiResponse<null>>(`/uploads/${id}`);
  if (!res.data.success) {
    throw new Error(res.data.error ?? "Failed to delete upload");
  }
}

export async function createReport(req: CreateReportRequest): Promise<Report> {
  const res = await client.post<ApiResponse<Report>>("/reports", req);
  if (!res.data.success || !res.data.data) {
    throw new Error(res.data.error ?? "Failed to create report");
  }
  return res.data.data;
}

export async function getReport(id: string): Promise<{
  report: Report;
  uploads: Upload[];
}> {
  const res = await client.get<ApiResponse<{ report: Report; uploads: Upload[] }>>(`/reports/${id}`);
  if (!res.data.success || !res.data.data) {
    throw new Error(res.data.error ?? "Failed to get report");
  }
  return res.data.data;
}

export async function listReports(): Promise<Report[]> {
  const res = await client.get<ApiResponse<Report[]>>("/reports");
  if (!res.data.success || !res.data.data) {
    throw new Error(res.data.error ?? "Failed to list reports");
  }
  return res.data.data;
}

export async function deleteReport(id: string): Promise<void> {
  const res = await client.delete<ApiResponse<null>>(`/reports/${id}`);
  if (!res.data.success) {
    throw new Error(res.data.error ?? "Failed to delete report");
  }
}

export interface ReportTemplate {
  id: string;
  name: string;
  description: string;
  category: string;
  default_title: string;
  classification: string;
  categories: string[];
  priorities: string[];
  statuses: string[];
  default_severity: string;
}

export async function listTemplates(): Promise<ReportTemplate[]> {
  const res = await client.get<ApiResponse<ReportTemplate[]>>("/templates");
  if (!res.data.success || !res.data.data) {
    throw new Error(res.data.error ?? "Failed to list templates");
  }
  return res.data.data;
}

export async function getTemplate(id: string): Promise<ReportTemplate> {
  const res = await client.get<ApiResponse<ReportTemplate>>(`/templates/${id}`);
  if (!res.data.success || !res.data.data) {
    throw new Error(res.data.error ?? "Failed to get template");
  }
  return res.data.data;
}

import type { Block } from "@/types/block";

export async function getBlocks(reportId: string): Promise<Block[]> {
  const res = await client.get<ApiResponse<Block[]>>(`/reports/${reportId}/blocks`);
  if (!res.data.success) throw new Error(res.data.error ?? "Failed to get blocks");
  return res.data.data ?? [];
}

export async function createBlock(reportId: string, type: string, content: string, position: number): Promise<Block> {
  const res = await client.post<ApiResponse<Block>>(`/reports/${reportId}/blocks`, { type, content, position });
  if (!res.data.success || !res.data.data) throw new Error(res.data.error ?? "Failed to create block");
  return res.data.data;
}

export async function updateBlock(reportId: string, blockId: string, updates: Record<string, unknown>): Promise<Block> {
  const res = await client.patch<ApiResponse<Block>>(`/reports/${reportId}/blocks/${blockId}`, updates);
  if (!res.data.success || !res.data.data) throw new Error(res.data.error ?? "Failed to update block");
  return res.data.data;
}

export async function deleteBlock(reportId: string, blockId: string): Promise<void> {
  await client.delete(`/reports/${reportId}/blocks/${blockId}`);
}

export async function reorderBlocks(reportId: string, blockIds: string[]): Promise<void> {
  const res = await client.post(`/reports/${reportId}/blocks/reorder`, { block_ids: blockIds });
  if (!res.data.success) throw new Error("Failed to reorder");
}

export async function previewReport(req: CreateReportRequest): Promise<string> {
  const res = await client.post("/reports/preview", req, {
    responseType: "text",
  });
  return res.data;
}
