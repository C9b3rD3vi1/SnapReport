export interface ReportInfo {
  title: string;
  project: string;
  company: string;
  author: string;
  version: string;
}

export interface Report extends ReportInfo {
  id: string;
  status: string;
  pdf_path?: string;
  created_at: string;
  updated_at: string;
}

export interface ReportUploadInput {
  id: string;
  title: string;
  description: string;
  notes: string;
  order_index: number;
}

export interface CreateReportRequest extends ReportInfo {
  uploads: ReportUploadInput[];
}
