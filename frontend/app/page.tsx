"use client";
import { useState } from "react";

export default function Home() {
  const [text, setText] = useState("");
  const [result, setResult] = useState("");
  const [loading, setLoading] = useState(false);

  const handleDump = async () => {
    setLoading(true);
    setResult("");

    // 1. send `text` to your Go backend's /dump endpoint using fetch (POST)
    await fetch(`${process.env.NEXT_PUBLIC_API_URL}/dump`, {
      method: "POST",
      body: text,
    });

    // 2. after that, wait a few seconds (or just move to checking /latest)
    // 3. call /latest, get the result
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/latest`, {
      method: "GET",
    });
    // Extract the body content as text
    const textData = await response.text();
    // 4. setResult(...) with what you got back
    setResult(textData);
    setLoading(false);

    // 5. setLoading(false)
  };

  return (
    <main style={{ padding: "2rem" }}>
      <h1>Focus Dashboard</h1>
      <textarea
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="Dump everything on your mind..."
        rows={5}
        style={{ width: "100%" }}
      />
      <button onClick={handleDump} disabled={loading || !text.trim()}>
        {loading ? "Processing..." : "Dump it"}
      </button>
      <p>{result}</p>
    </main>
  );
}
