import React, { useState } from 'react';
import './TableWithImages.css';

// Импортируем локальные изображения
import imageMinus1 from './assets/image-2.png';
import image0 from './assets/image-0.png';
import image1 from './/assets/image-1.png';

const TableWithImages = () => {
  const [matrix, setMatrix] = useState([
    [0, 0, 0],
    [0, 0, 0],
    [0, 0, 0],
  ]);

  const [selectedValue, setSelectedValue] = useState(1);
  // Новое состояние для хранения подсвеченных ячеек
  // Это будет массив объектов {row, col}, или другая структура
  const [highlightedCells, setHighlightedCells] = useState([]);

  const imageMap = {
    '-1': imageMinus1,
    '0': image0,
    '1': image1,
  };

  // --- Новая функция для проверки выигрышных комбинаций ---
  const checkWinningCombinations = (currentMatrix) => {
    const wins = []; // Массив для хранения координат всех ячеек, которые нужно подсветить
    const size = currentMatrix.length;

    // Вспомогательная функция для проверки линии (строки, столбца или диагонали)
    const checkLine = (line) => {
    // Комбинация считается выигрышной, если все 3 элемента линии одинаковы
    // И это значение не равно 0
    if (line[0] !== 0 && line[0] === line[1] && line[1] === line[2]) {
      return line[0]; // Возвращаем выигрышное значение (1 или -1)
    }
    return null; // Нет выигрышной комбинации
  };

    // Проверка строк
    for (let r = 0; r < size; r++) {
      const row = currentMatrix[r];
      if (checkLine(row) !== null) {
        for (let c = 0; c < size; c++) {
          wins.push({ row: r, col: c });
        }
      }
    }
      

    // Проверка столбцов
    for (let c = 0; c < size; c++) {
      const col = [];
      for (let r = 0; r < size; r++) {
        col.push(currentMatrix[r][c]);
      }
      if (checkLine(col) !== null) {
        for (let r = 0; r < size; r++) {
          wins.push({ row: r, col: c });
        }
      }
    }

    // Проверка главной диагонали (сверху-слева донизу-справа)
    const mainDiag = [];
    for (let i = 0; i < size; i++) {
      mainDiag.push(currentMatrix[i][i]);
    }
    if (checkLine(mainDiag) !== null) {
      for (let i = 0; i < size; i++) {
        wins.push({ row: i, col: i });
      }
    }

    // Проверка побочной диагонали (сверху-справа донизу-слева)
    const antiDiag = [];
    for (let i = 0; i < size; i++) {
      antiDiag.push(currentMatrix[i][size - 1 - i]);
    }
    if (checkLine(antiDiag) !== null) {
      for (let i = 0; i < size; i++) {
        wins.push({ row: i, col: size - 1 - i });
      }
    }

    // Убираем дубликаты (если одна ячейка участвует в нескольких комбинациях)
    const uniqueWins = Array.from(new Set(wins.map(JSON.stringify))).map(JSON.parse);
    
    setHighlightedCells(uniqueWins);
    //setHighlightedCells(wins);

    // Можно добавить логику, если игра окончена, например, если количество заполненных ячеек равно 9
    if (uniqueWins.length > 0) {
      console.log('Найдены выигрышные комбинации! Подсвечено:', uniqueWins);
      // alert('Выигрышная комбинация найдена!'); // Опционально
    } else {
      console.log('Выигрышных комбинаций пока нет.');
    }
   
  };
  // --- Конец функции проверки комбинаций ---


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

  const fetchNewMatrix = async () => {
    // При запросе новой матрицы, сбрасываем подсветку
    setHighlightedCells([]);

    const requestBody = {
      value: 1,
    };

    const url = '/api/get-new-matrix';

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
        if (data.matrix && Array.isArray(data.matrix) && data.message) {
          setMatrix(data.matrix); // Обновляем матрицу
          // После получения новой матрицы, сразу проверяем её на комбинации
          checkWinningCombinations(data.matrix);
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


  const handleCellClick = (rowIndex, colIndex) => {
    if (matrix[rowIndex][colIndex] !== 0) {
      console.log(`Ячейка [${rowIndex}][${colIndex}] со значением ${matrix[rowIndex][colIndex]} некликабельна.`);
      alert('Клик разрешен только по ячейкам со значением 0!');
      return;
    }

    const newMatrix = matrix.map(row => [...row]);
    newMatrix[rowIndex][colIndex] = selectedValue;

    setMatrix(newMatrix);
    // После обновления матрицы, проверяем на выигрышные комбинации
    checkWinningCombinations(newMatrix); // Вызываем проверку здесь!

    const valueToSend = 0;
    sendRequest(newMatrix, valueToSend);
  };

  const handleRadioChange = (event) => {
    setSelectedValue(parseInt(event.target.value, 10));
  };

  return (
    <div className="table-container">
      <h2>Таблица с картинками</h2>

      <button onClick={fetchNewMatrix} className="fetch-matrix-button">
        Получить новую матрицу
      </button>

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
              {row.map((cellValue, colIndex) => {
                // Проверяем, должна ли текущая ячейка быть подсвечена
                const isHighlighted = highlightedCells.some(
                  (cell) => cell.row === rowIndex && cell.col === colIndex
                );
                return (
                  <td
                    key={`${rowIndex}-${colIndex}`}
                    onClick={() => handleCellClick(rowIndex, colIndex)}
                    className={`${cellValue === 0 ? 'clickable' : ''} ${isHighlighted ? 'highlighted' : ''}`}
                    style={{ pointerEvents: cellValue === 0 ? 'auto' : 'none' }}
                  >
                    <img
                      src={imageMap[cellValue.toString()]}
                      alt={`Cell value: ${cellValue}`}
                    />
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
      <p>Текущая матрица (для отладки): {JSON.stringify(matrix)}</p>
      <p>Выбранное значение для установки: {selectedValue}</p>
      {/* Для отладки: <p>Подсвеченные ячейки: {JSON.stringify(highlightedCells)}</p> */}
    </div>
  );
};

export default TableWithImages;