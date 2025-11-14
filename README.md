# **E-Commerce API**
Це API написанний на **GO** з використанням **PostgreSQL** та принцепів **Clean Architecture**.
Проєкт підтримує управління товарами, користувачами, кошиками та замовленнями

# Технології
Go 1.24+

Docker & Docker Compose 3.9+

# Інструкція по запуску
## 1. Clone Repository
```git
git clone https://github.com/Xiancel/ecommerce.git
cd ecommerce
```

## 2. Налаштування
створіть **.env** файл
Скопіюйте **.env.example** в **.env** та налаштуйте змінні:
```txt
DB_USER=<your_db_user>
DB_PASSWORD=<your_db_password>
DB_NAME=<your_db_name>

JWT_SECRET=<your_jwt_secret>

ADMIN_EMAIL=<admin_email>
ADMIN_PASSWORD=<admin_password>
```

## 3. Запуск Docker
```bash
docker-compose --profile dev up --build
```
# Схема проєкту
```txt
/cmd
  /main.go            # Точка входу додатку
              
/internal                    
  /domain             # Домені моделі

  /service            # Реалізація бізнес-логіки
    /auth             # Автентифікація
    /cart             # Логіка кошика
    /order            # Обробка замовлень
    /product          # Управління товарами
    /user             # Управління користувачами
  
  /repository         # Інтерфейси Репозиторіїв
  /handler            # HTTP Layer
    /http             # REST handlers
    /middleware       # Auth, CORS, Logging
           
/migrations           # SQL міграції
/docs                 # Документація   
/scripts              # Допоміжні скрипти
```

# API Endpoints

## Автентифікація
```txt
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
```

## Товари (публічні)
```txt
GET  /api/v1/products
GET  /api/v1/products/:id
GET  /api/v1/products/search
GET  /api/v1/categories
```

## Кошик (тільки для авторизованних користувачів)
```txt
GET    /api/v1/cart
POST   /api/v1/cart/items
PUT    /api/v1/cart/items/:id
DELETE /api/v1/cart/items/:id
DELETE /api/v1/cart
```

## Замовлення (тільки для авторизованних користувачів)
```txt
POST   /api/v1/orders
GET    /api/v1/orders
GET    /api/v1/orders/:id
PUT    /api/v1/orders/:id/cancel
```

## Користувачі (тільки для авторизованних користувачів)
```txt
GET    /api/v1/users
PUT    /api/v1/users
```

## Адмін (тільки для ролі **admin**)
```txt
POST   /api/v1/admin/products
PUT    /api/v1/admin/products/:id
GET    /api/v1/admin/orders
GET    /api/v1/admin/users
GET    /api/v1/admin/statistics
```

# Документація
1. **Swagger** можно переглянути після запуску на endpoint **/swagger/index.html***
2. **Insomnia** у розроботці
