"use client";

import { useState } from "react";
import FileUpload from "./components/FileUpload";
import Chat from "./components/Chat";

export default function Home() {
  const [isFileUploaded, setIsFileUploaded] = useState(false);

  const handleFileUploaded = (success: boolean) => {
    setIsFileUploaded(success);
  };

  return (
    <main className="min-h-screen bg-gray-50 px-4">
      <div className="container mx-auto">
        <h1 className="text-3xl font-bold text-center py-6 text-gray-800">
          BoWatt RAG Chat Assistant
        </h1>

        {!isFileUploaded ? (
          <FileUpload onFileUploaded={handleFileUploaded} />
        ) : (
          <Chat />
        )}
      </div>
    </main>
  );
}
