import React from 'react';
import UrlShortener from './components/UrlShortner';
import './animation.css'; 

const App: React.FC = () => {
  return (
    <div className="bg-gray-900 min-h-screen flex justify-center items-center animate-gradient">
      <UrlShortener />
    </div>
  );
};

export default App;
