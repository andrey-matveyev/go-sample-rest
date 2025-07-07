import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css'; // Глобальные стили
import App from './App'; // Импортируем главный компонент App
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App /> {/* Рендерим компонент App */}
  </React.StrictMode>
);

// Если хотим измерять производительность вашего приложения:
reportWebVitals();