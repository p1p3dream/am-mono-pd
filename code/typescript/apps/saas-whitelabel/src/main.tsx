import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import './styles/globals.css';

// Create a comment explaining the change
// Temporarily disabling StrictMode to test if it affects hot reloading
// If hot reloading works after this change, we can re-enable StrictMode with a different approach
ReactDOM.createRoot(document.getElementById('root')!).render(
  // <React.StrictMode>
  <App />
  // </React.StrictMode>,
);
