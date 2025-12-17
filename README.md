# Тестовое задание: Кошелек (Wallet Service)

Это реализация простого кошелька с REST API, написанная на Go. Приложение позволяет создавать кошельки, пополнять их и снимать средства, а также проверять текущий баланс.

Особое внимание в проекте уделено безопасной обработке конкурентных запросов на списание с одного кошелька с использованием пессимистичных блокировок на уровне базы данных.

**Стек технологий:**
*   **Go**
*   **PostgreSQL** (в качестве базы данных)
*   **Docker** и **Docker Compose** (для развертывания)
*   **Gin** (для роутинга)
*   **pgx** (драйвер для PostgreSQL)
*   **golang-migrate** (для управления миграциями БД)
*   **zap** (для логирования)

---

## API Endpoints

### 1. Создание нового кошелька

Создает новый кошелек с нулевым балансом и возвращает его данные.

- **URL:** `/api/v1/wallets`
- **Method:** `POST`
- **Success Response (201 Created):**
  ```json
  {
      "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "balance": "0"
  }
  ```

### 2. Выполнение операции (пополнение/списание)

Выполняет операцию пополнения (`DEPOSIT`) или списания (`WITHDRAW`) для указанного кошелька.

- **URL:** `/api/v1/wallet`
- **Method:** `POST`
- **Request Body:**
  ```json
  {
      "walletId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "operationType": "DEPOSIT",
      "amount": "1000.50"
  }
  ```
- **Success Response (204 No Content):** Пустое тело ответа.
- **Error Responses:**
    - `404 Not Found`: если кошелек не найден.
    - `422 Unprocessable Entity`: если недостаточно средств для списания.

### 3. Получение баланса

Возвращает текущий баланс кошелька.

- **URL:** `/api/v1/wallets/{WALLET_UUID}`
- **Method:** `GET`
- **Success Response (200 OK):**
  ```json
  {
      "balance": "950.50"
  }
  ```
- **Error Responses:**
    - `404 Not Found`: если кошелек не найден.

---

## Инструкция по запуску

Для запуска проекта на локальной машине должны быть установлены **Docker** и **Docker Compose**.

1.  **Клонируйте репозиторий:**
    ```bash
    git clone https://github.com/DblMOKRQ/test_task.git
    cd test_task
    ```

2.  **Создайте файл с переменными окружения.**
    Создайте файл `config.env` в корне проекта. Вы можете скопировать его из примера:
    ```bash
    cp config.env.example config.env
    ```
    *При необходимости измените значения в `config.env` (например, пароль или порт).*

3.  **Запустите проект с помощью Docker Compose:**
    Эта команда соберет образ Go-приложения, поднимет контейнер с PostgreSQL и применит миграции базы данных.
    ```bash
    docker-compose up --build
    ```
    *Флаг `--build` нужен для пересборки образа приложения при изменениях в коде.*

4.  **Готово!**
    Приложение будет доступно по адресу `http://localhost:8080`. Вы можете использовать `curl` или Postman для отправки запросов.

    **Пример с `curl`:**
    ```bash
    # 1. Создаем кошелек
    curl -X POST http://localhost:8080/api/v1/wallets
    # {"id":"...","balance":"0"} -> скопируйте полученный ID

    # 2. Пополняем его (замените UUID на ваш)
    curl -X POST http://localhost:8080/api/v1/wallet -H "Content-Type: application/json" -d \
    '{"walletId": "PASTE_YOUR_UUID_HERE", "operationType": "DEPOSIT", "amount": "150"}'

    # 3. Проверяем баланс
    curl http://localhost:8080/api/v1/wallets/PASTE_YOUR_UUID_HERE
    # {"balance":"150"}
    ```

5.  **Остановка проекта:**
    Для остановки всех контейнеров нажмите `Ctrl+C` в терминале, где запущен `docker-compose`, или выполните команду в другой вкладке:
    ```bash
    docker-compose down
    ```
