import React, { useState } from 'react';
import './TableWithImages.css';

// Импортируем локальные изображения
import imageMinus1 from './assets/image-2.png';
import image0 from './assets/image-0.png';
import image1 from './assets/image-1.png';

const TableWithImages = () => {
  // Инициализируем матрицу нулями
  const [matrix, setMatrix] = useState([
    [0, 0, 0],
    [0, 0, 0],
    [0, 0, 0],
  ]);

  // Новое состояние для хранения выбранного значения из RadioGroup
  // Инициализируем его в 1 по умолчанию
  const [selectedValue, setSelectedValue] = useState(1);

  const imageMap = {
    '-1': imageMinus1,
    '0': image0,
    '1': image1,
  };

  const sendRequest = async (currentMatrix, clickedValue) => {
    const requestBody = {
      matrix: currentMatrix,
      value: clickedValue, // Это значение, которое было в кликнутой ячейке до изменения (т.е. 0)
    };

    const url = '/api/table-data';

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
        // Можно обновить матрицу из ответа сервера, если это необходимо
        // setMatrix(data.updatedMatrix); // Пример
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
    // Клик разрешен только по ячейкам со значением 0
    if (matrix[rowIndex][colIndex] !== 0) {
      console.log(`Ячейка [${rowIndex}][${colIndex}] со значением ${matrix[rowIndex][colIndex]} некликабельна.`);
      alert('Клик разрешен только по ячейкам со значением 0!');
      return;
    }

    const newMatrix = matrix.map(row => [...row]);
    // Устанавливаем значение ячейки согласно выбранному в RadioGroup
    newMatrix[rowIndex][colIndex] = selectedValue; // Используем selectedValue

    setMatrix(newMatrix);

    // В качестве параметра передаем старое значение ячейки (0)
    // или другое значение, которое вам нужно
    const valueToSend = 0; // Или selectedValue, в зависимости от вашей логики
    sendRequest(newMatrix, valueToSend);
  };

  // Обработчик изменения значения RadioGroup
  const handleRadioChange = (event) => {
    // Преобразуем значение из строки в число
    setSelectedValue(parseInt(event.target.value, 10));
  };

  return (
    <div className="table-container">
      <h2>Таблица с картинками</h2>

      <div className="radio-group-container">
        <h3>Выберите значение для установки:</h3>
        <label>
          <input
            type="radio"
            name="cellValue"
            value="1"
            checked={selectedValue === 1}
            onChange={handleRadioChange}
          />
          Установить 1
        </label>
        <label>
          <input
            type="radio"
            name="cellValue"
            value="-1"
            checked={selectedValue === -1}
            onChange={handleRadioChange}
          />
          Установить -1
        </label>
      </div>

      <table>
        <tbody>
          {matrix.map((row, rowIndex) => (
            <tr key={rowIndex}>
              {row.map((cellValue, colIndex) => (
                <td
                  key={`${rowIndex}-${colIndex}`}
                  onClick={() => handleCellClick(rowIndex, colIndex)}
                  // Класс 'clickable' только для ячеек со значением 0
                  className={cellValue === 0 ? 'clickable' : ''}
                  // Отключаем pointer-events для некликабельных ячеек
                  style={{ pointerEvents: cellValue === 0 ? 'auto' : 'none' }}
                >
                  <img
                    src={imageMap[cellValue.toString()]}
                    alt={`Cell value: ${cellValue}`}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
      
      <p>Текущая матрица (для отладки): {JSON.stringify(matrix)}</p>
      <p>Выбранное значение для установки: {selectedValue}</p> {/* Для отладки */}
    </div>
  );
};

export default TableWithImages;