export interface Upload {
  id: string;
  filename: string;
  original_name: string;
  mime_type: string;
  size: number;
  title: string;
  description: string;
  notes?: string;
  order_index: number;
  created_at: string;
}

export interface UploadFile {
  id: string;
  file: File;
  preview: string;
  progress: number;
  upload: Upload | null;
  error: string | null;
}
