const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const uploadFile = async (file: File): Promise<boolean> => {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(`${API_BASE_URL}/upload`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    throw new Error("Upload failed");
  }

  return true;
};

export const queryDocument = async (text: string): Promise<string> => {
  const response = await fetch(`${API_BASE_URL}/query`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ text }),
  });

  if (!response.ok) {
    throw new Error("Query failed");
  }

  const data = await response.json();
  return data.response;
};
