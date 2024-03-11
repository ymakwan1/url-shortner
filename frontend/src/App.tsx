// frontend/src/App.tsx

import React from 'react';
import UrlShortener from './components/UrlShortner';

const App: React.FC = () => {
  return (
    <div className="bg-gray-100 min-h-screen flex justify-center items-center">
      <UrlShortener />
    </div>
  );
};

export default App;
