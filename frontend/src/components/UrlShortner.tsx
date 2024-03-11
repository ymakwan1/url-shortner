import React, { useState } from 'react';
import '../animation.css';

const UrlShortener: React.FC = () => {
  const [originalUrl, setOriginalUrl] = useState('');
  const [shortenedUrl, setShortenedUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [uniqueIPs, setUniqueIPs] = useState(0);
  const [urlsShortened, setUrlsShortened] = useState(0);

  const handleShortenUrl = async () => {
    if (!originalUrl) {
      setError('Please enter a URL to shorten.');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      /**
       * @todo Make request to backend
       */
      const response = await fetch('', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ originalUrl }),
      });

      if (response.ok) {
        const data = await response.json();
        setShortenedUrl(data.shortenedUrl);
        setUrlsShortened(urlsShortened + 1);
      } else {
        console.error('Failed to shorten URL');
        setError('Failed to shorten URL. Please try again.');
      }
    } catch (error) {
      console.error('Error occurred while shortening URL:', error);
      setError('An error occurred while shortening the URL. Please try again.');
    }

    setIsLoading(false);
  };

  const handleCopyUrl = () => {
    navigator.clipboard.writeText(shortenedUrl);
    alert('Shortened URL copied to clipboard!');
  };

  const trackUniqueIPs = () => {
    setUniqueIPs(uniqueIPs + 1);
  };

  return (
    <div className="max-w-lg mx-auto p-6 animate-gradient rounded-md shadow-md">
      <h1 className="text-3xl font-semibold mb-4 text-white">URL Shortener</h1>
      <div className="flex items-center mb-4">
        <input
          type="text"
          placeholder="Enter URL to shorten"
          value={originalUrl}
          onChange={(e) => setOriginalUrl(e.target.value)}
          className="flex-grow border rounded-md px-16 py-3 mr-2 w-full max-w-lg text-lg"
        />
        <button
          onClick={() => { handleShortenUrl(); trackUniqueIPs(); }}
          disabled={isLoading}
          className="bg-blue-500 text-white font-semibold py-3 px-4 rounded-md hover:bg-blue-600"
        >
          {isLoading ? 'Shortening...' : 'Shorten'}
        </button>
      </div>
      {error && <p className="text-red-500">{error}</p>}
      {shortenedUrl && (
        <div className="flex items-center mb-4">
          <input
            type="text"
            value={shortenedUrl}
            readOnly
            className="flex-grow border rounded-md px-4 py-3 mr-2 w-full max-w-lg text-lg"
          />
          <button
            onClick={handleCopyUrl}
            className="bg-blue-500 text-white font-semibold py-3 px-4 rounded-md hover:bg-blue-600"
          >
            Copy
          </button>
        </div>
      )}
      <div className="text-white">
        <p>Unique IPs served: {uniqueIPs}</p>
        <p>Number of URLs shortened: {urlsShortened}</p>
      </div>
    </div>
  );
};

export default UrlShortener;
