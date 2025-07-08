import React, { useState } from 'react';
import './TableWithImages.css';

// Импортируем локальные изображения
import imageMinus1 from './assets/image-2.png';
import image0 from './assets/image-0.png';
import image1 from './assets/image-1.png';

const TableWithImages = () => {
  const [matrix, setMatrix] = useState([
    [0, 0, 0],
    [0, 0, 0],
    [0, 0, 0],
  ]);

  const [selectedValue, setSelectedValue] = useState(1); // Состояние для RadioGroup

  const imageMap = {
    '-1': imageMinus1,
    '0': image0,
    '1': image1,
  };

  // --- Новая функция для отправки POST-запроса на получение новой матрицы ---
  const fetchNewMatrix = async () => {
    const requestBody = {
      value: 1, // Отправляем int = 1, как договорились
    };

    const url = '/api/get-new-matrix'; // Новый эндпоинт для этого запроса

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
        // Ожидаем, что ответ будет содержать поля 'matrix' и 'message'
        if (data.matrix && Array.isArray(data.matrix) && data.message) {
          setMatrix(data.matrix); // Обновляем матрицу полученной с сервера
          alert(`Матрица обновлена! Сообщение от сервера: ${data.message}`);
          console.log('Новая матрица получена:', data.matrix);
        } else {
          console.error('Некорректный формат ответа от сервера:', data);
          alert('Получен некорректный формат данных от сервера!');
        }
      } else {
        const errorText = await response.text();
        console.error('Ошибка при получении новой матрицы:', response.status, response.statusText, errorText);
        alert(`Ошибка при получении новой матрицы: ${response.statusText}`);
      }
    } catch (error) {
      console.error('Произошла ошибка сети при запросе новой матрицы:', error);
      alert('Произошла ошибка сети при запросе новой матрицы!');
    }
  };
  // --- Конец новой функции ---


  const sendRequest = async (currentMatrix, clickedValue) => {
    const requestBody = {
      matrix: currentMatrix,
      value: clickedValue,
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
    if (matrix[rowIndex][colIndex] !== 0) {
      console.log(`Ячейка [${rowIndex}][${colIndex}] со значением ${matrix[rowIndex][colIndex]} некликабельна.`);
      alert('Клик разрешен только по ячейкам со значением 0!');
      return;
    }

    const newMatrix = matrix.map(row => [...row]);
    newMatrix[rowIndex][colIndex] = selectedValue;

    setMatrix(newMatrix);

    const valueToSend = 0;
    sendRequest(newMatrix, valueToSend);
  };

  const handleRadioChange = (event) => {
    setSelectedValue(parseInt(event.target.value, 10));
  };

  return (
    <div className="table-container">
      <h2>Таблица с картинками</h2>

      {/* --- Новая кнопка здесь --- */}
      <button onClick={fetchNewMatrix} className="fetch-matrix-button">
        Получить новую матрицу
      </button>
      {/* --- Конец новой кнопки --- */}

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
                  className={cellValue === 0 ? 'clickable' : ''}
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
      <p>Выбранное значение для установки: {selectedValue}</p>
    </div>
  );
};

export default TableWithImages;