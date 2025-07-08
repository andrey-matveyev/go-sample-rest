import React, { useState, useEffect } from 'react';
import './TableWithImages.css';

// Импортируем локальные изображения
import imageMinus1 from './assets/image-2.png'; // Укажите правильный путь и расширение
import image0 from './assets/image-0.png';           // Укажите правильный путь и расширение
import image1 from './assets/image-1.png';           // Укажите правильный путь и расширение

const TableWithImages = () => {
  // Инициализируем матрицу значениями, где кликабельны только '0'
  // Пример: можно инициализировать матрицу случайными -1, 0, 1
  const [matrix, setMatrix] = useState([
    [0, 0, 0],
    [0, 0, 0],
    [0, 0, 0],
  ]);

  // Объект для быстрого сопоставления значений и изображений
  const imageMap = {
    '-1': imageMinus1,
    '0': image0,
    '1': image1,
  };

  // Опционально: Загрузка картинок с бэкенда.
  // Если вы хотите грузить URL картинок с бэкенда,
  // вам понадобится состояние для их хранения и useEffect для загрузки.
  /*
  const [remoteImageUrls, setRemoteImageUrls] = useState({});
  useEffect(() => {
    const fetchImageUrls = async () => {
      try {
        const response = await fetch('/api/images'); // Ваш эндпоинт на бэкенде для картинок
        if (response.ok) {
          const data = await response.json();
          setRemoteImageUrls(data); // Ожидаем { "-1": "url1", "0": "url2", "1": "url3" }
        } else {
          console.error('Failed to fetch image URLs from backend');
        }
      } catch (error) {
        console.error('Network error while fetching image URLs:', error);
      }
    };
    fetchImageUrls();
  }, []);
  // Тогда в JSX вы бы использовали imageMap['-1'] = remoteImageUrls['-1']
  // или сделали бы проверку if (remoteImageUrls['0']) { ... }
  */


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
    // Разрешаем клик только по ячейкам со значением 0
    if (matrix[rowIndex][colIndex] !== 0) {
      console.log(`Ячейка [${rowIndex}][${colIndex}] со значением ${matrix[rowIndex][colIndex]} некликабельна.`);
      alert('Клик разрешен только по ячейкам со значением 0!');
      return;
    }

    const newMatrix = matrix.map(row => [...row]);
    // Изменяем значение ячейки с 0 на 1 при клике
    newMatrix[rowIndex][colIndex] = 1;

    setMatrix(newMatrix);

    // Целое число для параметра (например, само значение ячейки до изменения)
    const clickedValue = 0; // Или любое другое число, которое вам нужно передать как параметр

    sendRequest(newMatrix, clickedValue);
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
                  // Добавляем класс 'clickable' только для ячеек со значением 0,
                  // чтобы стилизовать курсор и т.д.
                  className={cellValue === 0 ? 'clickable' : ''}
                  // Отключаем pointer-events для некликабельных ячеек
                  style={{ pointerEvents: cellValue === 0 ? 'auto' : 'none' }}
                >
                  <img
                    // Выбираем изображение в зависимости от значения ячейки
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
    </div>
  );
};

export default TableWithImages;