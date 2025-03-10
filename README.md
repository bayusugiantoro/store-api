# 📘 Store Redemption API Documentation

## 🚀 Overview

Brand & Voucher API provides various endpoints for managing brands, vouchers and redemption transactions.

---

## 🏗️ Tech

This project was built using the following technology :

- **Golang** — The main programming language for building API.
- **PostgreSQL** — Relational database to store brand, voucher and transaction information.
- **golang-migrate** — Tools for managing database migrations.
- **golang-http-router** — HTTP routers are light and fast to handle routing.
- **Testify** — Library for unit testing in Golang.
- **Golang Validator** — Library for validating data on struct.
- **Postman** — Used to manually test API and document requests/responses.

---

## 📦 Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/bayusugiantoro/store-api.git
```

### 2. Install Dependencies

```bash
go get
```

### 3. Run Database Migration

Using CLI migration tool like `golang-migrate` or `goose`.

```bash
migrate -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" -path migrations up
```

### 4. ERD

![App Screenshot](./assets/RDBMS.JPG)

---

### 5. Run Server

```bash
go run main.go
```

---

## 📡 API Endpoints

### 1. Create Brand

- **Method:** `POST`
- **URL:** `http://localhost:3000/brand`

---

### 2. Get Single Brand

- **Method:** `GET`
- **URL:** `http://localhost:3000/brand?id={brand_id}`

---

### 3. Get All Brand

- **Method:** `GET`
- **URL:** `http://localhost:3000/brand`

---

### 4. Create Voucher

- **Method:** `POST`
- **URL:** `http://localhost:3000/voucher`

---

### 5. Get Single Voucher

- **Method:** `GET`
- **URL:** `http://localhost:3000/voucher?id={voucher_id}`

---

### 6. Get All Vouchers by Brand

- **Method:** `GET`
- **URL:** `http://localhost:3000/brand?id={brand_id}/vouchers`

---

### 7. Create Transaction

- **Method:** `POST`
- **URL:** `http://localhost:3000/transaction/redemption`

---

### 8. Get Transaction Detail

- **Method:** `GET`
- **URL:** `http://localhost:3000/transaction/redemption?transactionId={transactionId}`
