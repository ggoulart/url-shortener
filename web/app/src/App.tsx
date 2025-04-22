// src/App.tsx
import {useState} from "react";

function App() {
    const [url, setUrl] = useState("");
    const [shortened, setShortened] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setShortened("");
        setLoading(true);
        try {
            const res = await fetch("http://localhost:8080/api/v1/shorten", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({"longUrl": url}),
            });

            if (!res.ok) throw new Error("Failed to shorten URL");

            const data = await res.json();
            setShortened(data.shortUrl);
        } catch (err) {
            setError((err as Error).message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex flex-col items-center justify-center bg-[#f8f9fa] px-4">
            <div className="w-full max-w-xl">
                <form onSubmit={handleSubmit} className="flex flex-col items-center gap-4">
                    <h1 className="text-3xl font-semibold text-gray-800 mb-6">
                        URL Shortener
                    </h1>
                    <input
                        type="url"
                        placeholder="Paste a URL here"
                        className="w-full px-5 py-3 rounded-full border border-gray-300 shadow focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        required
                    />
                    <button
                        type="submit"
                        className="bg-[#4285f4] text-white text-sm px-6 py-2 rounded hover:bg-[#357ae8] transition-colors duration-200"
                        disabled={loading}
                    >
                        {loading ? "Shortening..." : "Shorten URL"}
                    </button>
                </form>

                {shortened && (
                    <div className="mt-6 text-center text-green-700">
                        <p className="text-md">Shortened URL:</p>
                        <a
                            href={shortened}
                            target={shortened}
                            rel="noopener noreferrer"
                            className="text-blue-600 underline break-all"
                        >
                            {shortened}
                        </a>
                    </div>
                )}

                {error && (
                    <p className="mt-6 text-center text-red-600">{error}</p>
                )}
            </div>
        </div>
    );
}

export default App;
