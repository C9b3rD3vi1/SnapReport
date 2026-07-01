import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { FileDropzone } from "./FileDropzone";

describe("FileDropzone", () => {
  it("renders the dropzone text", () => {
    render(<FileDropzone onFilesSelected={() => {}} />);
    expect(screen.getByText(/Drop screenshots here/i)).toBeInTheDocument();
  });

  it("shows accepted file types", () => {
    render(<FileDropzone onFilesSelected={() => {}} />);
    expect(screen.getByText(/PNG|JPG|WEBP/i)).toBeInTheDocument();
  });

  it("applies disabled styles when disabled", () => {
    render(<FileDropzone onFilesSelected={() => {}} disabled />);
    const text = screen.getByText(/Drop screenshots here/i);
    const container = text.parentElement!;
    expect(container.className).toContain("pointer-events-none");
  });

  it("renders without crashing", () => {
    const fn = vi.fn();
    const { container } = render(<FileDropzone onFilesSelected={fn} />);
    expect(container.querySelector('input[type="file"]')).toBeInTheDocument();
  });
});
