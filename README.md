# 🎟️ Ticketing API - Route Documentation

This API powers a simple ticketing system with user authentication, event management, ticket purchases, and reports. Below is the list of available endpoints grouped by functionality and access role.

## Base URL


---

## 🔐 Auth Routes (Public Access)

| Method | Endpoint      | Description           |
|--------|---------------|-----------------------|
| POST   | `/register`   | Register a new user   |
| POST   | `/login`      | Login and get JWT     |

---

## 👤 User Routes (Admin Only)

All routes under `/api/users` require **admin** role.

| Method | Endpoint      | Description         |
|--------|---------------|---------------------|
| GET    | `/users/`     | Get all users       |
| GET    | `/users/:id`  | Get user by ID      |

---

## 📅 Event Routes

### Public Access

| Method | Endpoint        | Description        |
|--------|-----------------|--------------------|
| GET    | `/events`       | Get all events     |
| GET    | `/events/:id`   | Get event by ID    |

### Admin Only (Authenticated)

| Method | Endpoint        | Description        |
|--------|-----------------|--------------------|
| POST   | `/events`       | Create an event    |
| PUT    | `/events/:id`   | Update an event    |
| DELETE | `/events/:id`   | Delete an event    |

---

## 🎫 Ticket Routes (User Only)

All routes under `/api/tickets` require **user** role.

| Method | Endpoint                   | Description                        |
|--------|----------------------------|------------------------------------|
| POST   | `/tickets`                 | Purchase a ticket                  |
| GET    | `/tickets`                 | Get user's tickets                 |
| GET    | `/tickets/:id`             | Get ticket details by ID          |
| PATCH  | `/tickets/:id`             | Cancel a ticket                    |
| PATCH  | `/tickets/:id/payment`     | Confirm ticket payment             |
| PATCH  | `/tickets/:id/cancel-payment` | Cancel ticket payment           |

---

## 📊 Report Routes (Admin Only)

All routes under `/api/reports` require **admin** role.

| Method | Endpoint              | Description                    |
|--------|-----------------------|--------------------------------|
| GET    | `/reports/summary`    | Get summary report             |
| GET    | `/reports/events`     | Get event sales reports        |
| GET    | `/reports/ticket`     | Get all purchased tickets      |

---

## 🔒 Role-based Access

| Role   | Access Level                                      |
|--------|---------------------------------------------------|
| Public | Register, Login, View Events                     |
| User   | Purchase & manage tickets                        |
| Admin  | Manage users, events, tickets, and view reports  |

---

## 🛡️ Middleware Notes

- `AuthMiddleware("admin")`: restricts access to admins.
- `AuthMiddleware("user")`: restricts access to authenticated users.
- JWT is required in `Authorization` header:  


---

## 🏁 Getting Started

1. Clone the repo
2. Run migrations & seeders
3. Start the server:  
 ```bash
 go run main.go

Access API via http://localhost:8080/api
