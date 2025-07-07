import React from 'react';
import './App.css'; // Общие стили для App, если есть
import TableWithImages from './TableWithImages'; // Импортируем компонент

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Мое React-приложение</h1>
      </header>
      <TableWithImages /> {/* Вставляем компонент */}
    </div>
  );
}

export default App;