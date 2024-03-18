import React, { useState } from 'react';

const UrlShortener: React.FC = () => {
  const [url, setUrl] = useState('');
  const [shortenedUrl, setShortenedUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [showCopyButton, setShowCopyButton] = useState(true);
  const [showCopyNotification, setShowCopyNotification] = useState(false);

  const handleShortenUrl = async () => {
    if (!url) {
      setError('Please enter a URL to shorten.');
      resetError();
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await fetch('http://localhost:3000/', {
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

        // Hide copy button and shortened URL field after 20 seconds
        setTimeout(() => {
          setShortenedUrl('');
          setShowCopyButton(false);
        }, 20000);
      } else {
        setUrl('');
        setIsLoading(false);
        console.error('Failed to shorten URL');
        setError('Failed to shorten URL. Please try again.');
        resetError();
      }
    } catch (error) {
      console.error('Error occurred while shortening URL:', error);
      setError('An error occurred while shortening the URL. Please try again.');
      resetError();
    }
  };

  const handleCopyUrl = () => {
    
    navigator.clipboard.writeText(shortenedUrl);
    setShowCopyNotification(true);
    
    
    setTimeout(() => {
      setUrl('');
      setShortenedUrl('');
      setShowCopyButton(false);
      setShowCopyNotification(false);
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
        <h1 className="text-3xl font-semibold mb-4 text-white">URL Shortener</h1>
        <div className="flex items-center mb-4">
          <input
            type="text"
            placeholder="Enter URL to shorten"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            className="flex-grow border rounded-md px-4 py-3 mr-2 w-full max-w-lg text-lg"
          />
          <button
            onClick={handleShortenUrl}
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
            {showCopyButton && (
              <button
                onClick={handleCopyUrl}
                className="bg-blue-500 text-white font-semibold py-3 px-4 rounded-md hover:bg-blue-600"
              >
                Copy
              </button>
            )}
          </div>
        )}
        {showCopyNotification && (
          <div className="bg-green-500 text-white px-4 py-2 rounded-md">
            URL copied to clipboard!
          </div>
        )}
      </div>
    </div>
  );
};

export default UrlShortener;
