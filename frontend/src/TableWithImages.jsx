import React, { useState } from 'react';
import './TableWithImages.css'; // Импортируем стили для этого компонента

const TableWithImages = () => {
  const [matrix, setMatrix] = useState([
    [0, 0, 0],
    [0, 0, 0],
    [0, 0, 0],
  ]);

  const images = [
    'https://via.placeholder.com/100x100?text=Image1',
    'https://via.placeholder.com/100x100?text=Image2',
    'https://via.placeholder.com/100x100?text=Image3',
    'https://via.placeholder.com/100x100?text=Image4',
    'https://via.placeholder.com/100x100?text=Image5',
    'https://via.placeholder.com/100x100?text=Image6',
    'https://via.placeholder.com/100x100?text=Image7',
    'https://via.placeholder.com/100x100?text=Image8',
    'https://via.placeholder.com/100x100?text=Image9',
  ];

  const sendRequest = async (currentMatrix, clickedValue) => {
    const requestBody = {
      matrix: currentMatrix,
      value: clickedValue,
    };

    const url = '/api/table-data'; // Эндпоинт для POST-запроса

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      if (response.ok) {
        const data = await response.json();
        console.log('Данные успешно отправлены. Ответ сервера:', data);
        alert('Запрос успешно отправлен!');
      } else {
        const errorText = await response.text();
        console.error('Ошибка при отправке запроса:', response.status, response.statusText, errorText);
        alert(`Ошибка при отправке запроса: ${response.statusText}`);
      }
    } catch (error) {
      console.error('Произошла ошибка сети:', error);
      alert('Произошла ошибка сети!');
    }
  };

  const handleCellClick = (rowIndex, colIndex) => {
    if (matrix[rowIndex][colIndex] === 1) {
      console.log(`Ячейка [${rowIndex}][${colIndex}] уже выбрана. Невозможно кликнуть повторно.`);
      alert('Эта ячейка уже выбрана!');
      return;
    }

    const newMatrix = matrix.map(row => [...row]);
    newMatrix[rowIndex][colIndex] = 1;

    setMatrix(newMatrix);

    const randomInt = Math.floor(Math.random() * 100);

    sendRequest(newMatrix, randomInt);
  };

  return (
    <div className="table-container">
      <h2>Таблица с картинками</h2>
      <table>
        <tbody>
          {matrix.map((row, rowIndex) => (
            <tr key={rowIndex}>
              {row.map((cellValue, colIndex) => (
                <td
                  key={`${rowIndex}-${colIndex}`}
                  onClick={() => handleCellClick(rowIndex, colIndex)}
                  className={cellValue === 1 ? 'selected' : ''}
                  style={{ pointerEvents: cellValue === 1 ? 'none' : 'auto' }}
                >
                  <img
                    src={images[rowIndex * 3 + colIndex]}
                    alt={`Cell ${rowIndex}-${colIndex}`}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
      <p>Текущая матрица (для отладки): {JSON.stringify(matrix)}</p>
    </div>
  );
};

export default TableWithImages;