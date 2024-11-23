# Contact Manager gRPC Service

## Описание

**Contact Manager** — gRPC-сервис для управления контактами. Он предоставляет возможность создания, получения, поиска и удаления контактов.
Proto контракт - https://github.com/tendze/gRPC_ContactManager_Protos. Также в сервис интегрирован сервис авторизации через interceptors - https://github.com/tendze/gRPC_AuthService_SSO

### Функционал:
- **Создание контактов:** Сохраняйте ключевую информацию о пользователях, включая имя, email и номер телефона.
- **Поиск по параметрам:** Мгновенно находите контакты по имени, email или номеру телефона.
- **Удаление контактов:** Удобное управление данными с возможностью безопасного удаления.

---

## Структура API

### Сервис `ContactManager`

#### Методы:
1. **CreateContact**  
   Создает новый контакт.  
   **Вход:**
    - `name` (string) — имя контакта.
    - `email` (string) — email контакта.
    - `phone` (string) — номер телефона.  
      **Выход:**
    - `id` (int64) — уникальный идентификатор контакта.
    - `success` (bool) — статус операции.

2. **GetContactByName**  
   Ищет контакт по имени.  
   **Вход:**
    - `name` (string) — имя контакта.  
      **Выход:**
    - `id` (int64), `name` (string), `email` (string), `phone` (string) — информация о контакте.

3. **GetContactByEmail**  
   Ищет контакт по email.  
   **Вход:**
    - `email` (string) — email контакта.  
      **Выход:**
    - `id` (int64), `name` (string), `email` (string), `phone` (string) — информация о контакте.

4. **GetContactByPhone**  
   Ищет контакт по номеру телефона.  
   **Вход:**
    - `phone` (string) — номер телефона контакта.  
      **Выход:**
    - `id` (int64), `name` (string), `email` (string), `phone` (string) — информация о контакте.

5. **DeleteContact**  
   Удаляет контакт по идентификатору.  
   **Вход:**
    - `id` (int64) — идентификатор контакта.  
      **Выход:**
    - `success` (bool) — статус операции.

---

### Технологии:
- **gRPC:** Быстрый и эффективный протокол коммуникации.
- **ProtoBuf (Protocol Buffers):** Используется для определения контракта API и сериализации данных.
- **Docker:** Сервис может быть упакован в контейнер для легкого деплоя.
- **Интеграция с AuthService:** Для защиты данных используется валидация JWT токенов через внешний сервис авторизации.

---

## Как запустить?

### Требования:
- Go 1.21+
- gRPC & Protocol Buffers
- Docker (опционально)

### Запуск локально:
1. Установите зависимости:
   ```bash
   go mod tidy
2. Примените файлы миграции через команду:
   ```bash
   task migrate-up
3. Затем запустите сам сервис
    ```bash
   task cm