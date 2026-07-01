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
