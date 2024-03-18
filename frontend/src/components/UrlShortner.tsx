import React, { useState } from 'react';
import '../animation.css';
const UrlShortener: React.FC = () => {
  const [url, setUrl] = useState('');
  const [shortenedUrl, setShortenedUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [showCopyButton, setShowCopyButton] = useState(true);
  const [showCopyNotification, setShowCopyNotification] = useState(false);
  const [copyButtonDisabled, setCopyButtonDisabled] = useState(false);
  const [copyFieldDisabled, setCopyFieldDisabled] = useState(false);

  const handleShortenUrl = async () => {
    if (!url) {
      setError('Please enter a URL to shorten.');
      resetError();
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await fetch('http://127.0.0.1:3000/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url }),
      });

      if (response.ok) {
        const data = await response.json();
        setShortenedUrl(data.short_url);
        setIsLoading(false);
        setShowCopyButton(true);
      } else {
        setUrl('');
        setIsLoading(false);
        setError('Failed to shorten URL. Please try again.');
        resetError();
      }
    } catch (error) {
      setError('An error occurred while shortening the URL. Please try again.');
      resetError();
    }
  };

  const handleCopyUrl = () => {
    navigator.clipboard.writeText(shortenedUrl);
    setShowCopyNotification(true);
    

    setTimeout(() => {
      setShowCopyButton(false);
      setShortenedUrl('');
      setShowCopyNotification(false);
      setCopyFieldDisabled(false);
      setCopyButtonDisabled(false);
      setUrl('');    
    }, 5000);
  };

  const resetError = () => {
    setTimeout(() => {
      setError('');
    }, 5000);
  };

  return (
    <div className="bg-gray-900 min-h-screen flex animate-gradient justify-center items-center">
      <div className="max-w-lg mx-auto p-6 animate-gradient rounded-md shadow-md">
        <h1 className="text-2xl font-semibold mb-4 text-white">URL Shortener</h1>
        <div className="flex items-center mb-4">
          <input
            type="text"
            placeholder="Enter URL to shorten"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            className="flex-grow border rounded-md px-6 py-2 mr-2 w-full max-w-lg text-lg"
          />
          <button
            onClick={handleShortenUrl}
            disabled={isLoading}
            className="bg-blue-500 text-white font-semibold py-2 px-4 rounded-md hover:bg-blue-600"
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
              className="flex-grow border rounded-md px-6 py-2 mr-2 w-full max-w-lg text-lg"
              disabled={copyFieldDisabled}
            />
            {showCopyButton && (
              <button
                onClick={handleCopyUrl}
                disabled={copyButtonDisabled}
                className="bg-blue-500 text-white font-semibold py-2 px-4 rounded-md hover:bg-blue-600"
              >
                Copy
              </button>
            )}
          </div>
        )}
        {showCopyNotification && (
          <div className="bg-green-500 text-white px-2 py-1 rounded-md">
            URL copied to clipboard!
          </div>
        )}
      </div>
    </div>
  );
};

export default UrlShortener;
