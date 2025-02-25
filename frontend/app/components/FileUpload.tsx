"use client";

import { useState } from "react";
import { CloudArrowUpIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";

interface FileUploadProps {
  onFileUploaded: (success: boolean) => void;
}

export default function FileUpload({ onFileUploaded }: FileUploadProps) {
  const [file, setFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string>("");

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const droppedFile = e.dataTransfer.files[0];
      if (droppedFile.type === "text/plain") {
        setFile(droppedFile);
        setError("");
      } else {
        setError("Please upload a text file (.txt)");
      }
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];
      if (selectedFile.type === "text/plain") {
        setFile(selectedFile);
        setError("");
      } else {
        setError("Please upload a text file (.txt)");
      }
    }
  };

  const handleUpload = async () => {
    if (!file) return;

    setIsUploading(true);
    setError("");

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Upload failed");
      }

      onFileUploaded(true);
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (err) {
      setError("Failed to upload file. Please try again.");
      onFileUploaded(false);
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto p-6">
      <div
        className={classNames(
          "border-2 border-dashed rounded-lg p-8",
          "transition-colors duration-200",
          {
            "border-blue-500 bg-blue-50": isDragging,
            "border-gray-300 hover:border-gray-400": !isDragging,
          }
        )}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <CloudArrowUpIcon className="mx-auto h-12 w-12 text-gray-400" />
        <div className="mt-4 text-center">
          <label className="cursor-pointer">
            <span className="text-blue-600 hover:text-blue-700">
              Click to upload
            </span>
            <input
              type="file"
              className="hidden"
              onChange={handleFileChange}
              accept=".txt"
            />
          </label>
          <span className="text-gray-500"> or drag and drop</span>
        </div>
        <p className="mt-2 text-center text-gray-500 text-sm">
          Text files (.txt) up to 10MB
        </p>

        {error && (
          <p className="mt-2 text-center text-red-500 text-sm">{error}</p>
        )}

        {file && !error && (
          <div className="mt-4 text-center">
            <p className="text-sm text-gray-600">Selected file: {file.name}</p>
            <button
              onClick={handleUpload}
              disabled={isUploading}
              className="mt-2 btn-primary"
            >
              {isUploading ? "Uploading..." : "Upload File"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
