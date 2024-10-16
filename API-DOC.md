# API Documentation

## Table of Contents
- [API Documentation](#api-documentation)
  - [Table of Contents](#table-of-contents)
  - [Authentication](#authentication)
    - [Signup](#signup)
    - [Complete Signup](#complete-signup)
    - [Login](#login)
    - [Complete Login](#complete-login)
  - [User API](#user-api)
    - [See Profile](#see-profile)
    - [Create Product](#create-product)
  - [Product API](#product-api)
    - [Get All Products](#get-all-products)
    - [Get All Mushroom Products](#get-all-mushroom-products)
    - [Get All Jari Products](#get-all-jari-products)
    - [Get All Products of a Farmer](#get-all-products-of-a-farmer)
    - [Get A Product By ID](#get-a-product-by-id)
    - [Update Product Unavailability](#update-product-unavailability)
    - [Delete Product](#delete-product)
  - [Order API](#order-api)
    - [Create Order](#create-order)

## Authentication

### Signup

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/auth/signup`
- Body:
```json
{
  "first_name": "Rohan",
  "last_name": "Sharma",
  "aadhar_number": "123412341219",
  "email": "rohan.sharma@example.com",
  "phone_number": "6200059008",
  "is_farmer": true,
  "address": "123 Green Farm Lane",
  "city": "Jhabua",
  "state": "Madhya Pradesh",
  "pin_code": "456001",
  "farm_size": "3"
}
```

### Complete Signup

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/auth/complete-signup`
- Body:
```json
{
  "user": {
    "first_name": "Rohan",
    "last_name": "Sharma",
    "aadhar_number": "582228381168",
    "email": "rohan.sharma@example.com",
    "phone_number": "6200059008",
    "is_farmer": true,
    "address": "123 Green Farm Lane",
    "city": "Jhabua",
    "state": "Madhya Pradesh",
    "pin_code": "456001",
    "farm_size": "3"
  },
  "verification_code": "355193"
}
```

### Login

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/auth/login`
- Body:
```json
{
  "phone_number": "6200059008",
  "aadhar_number": "582228381168",
  "verification_code": ""
}
```

### Complete Login

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/auth/complete-login`
- Body:
```json
{
  "phone_number": "6200059008",
  "aadhar_number": "582228381168",
  "verification_code": "935279"
}
```

## User API

### See Profile

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/user/1/profile`

**Response:**
```json
{
  "id": "1",
  "first_name": "Rohan",
  "last_name": "Sharma",
  "aadhar_number": "582228381168",
  "email": "rohan.sharma@example.com",
  "phone_number": "6200059008",
  "is_farmer": true,
  "address": "123 Green Farm Lane",
  "city": "Jhabua",
  "state": "Madhya Pradesh",
  "pin_code": "456001",
  "farm_size": "3",
  "password": "",
  "is_verified_by_admin": false,
  "created_at": "2024-10-16T17:22:24.208101Z",
  "updated_at": "2024-10-16T17:22:24.208101Z",
  "last_login_at": "2024-10-16T12:17:19.975054Z",
  "img": "",
  "user_type": "farmer"
}
```

### Create Product

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/v1/user/1/newproduct`
- Body:
```json
{
  "img": "https://example.com/images/product_101.jpg",
  "farmer_id": 1,
  "name": "jari",
  "type": "organic",
  "quantity_in_kg": 2500,
  "rate_per_kg": 200.36,
  "jari_size": "small",
  "expected_delivery": null,
  "created_at": "2024-09-15T14:30:00Z",
  "updated_at": "2024-09-25T10:00:00Z",
  "farmer_phone_number": "9876543219"
}
```

**Response:**
```json
{
  "message": "product created successfully!"
}
```

## Product API

### Get All Products

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/product`

**Response:**
```json
[
  {
    "id": 1,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "mushroom",
    "type": "organic",
    "quantity": 1200,
    "rate_per_kg": 280.65,
    "created_at": "2024-10-16T12:50:03.476421Z",
    "updated_at": "2024-10-16T12:50:03.476421Z",
    "farmer_phone_number": 6200059008,
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  }
]
```

### Get All Mushroom Products

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/product/mushroom`

**Response:**
```json
[
  {
    "id": 4,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "mushroom",
    "type": "non-organic",
    "quantity_in_kg": 5000,
    "rate_per_kg": 865.36,
    "created_at": "2024-10-16T14:08:43.942844Z",
    "updated_at": "2024-10-16T14:08:43.942844Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  },
  {
    "id": 3,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "mushroom",
    "type": "organic",
    "quantity_in_kg": 1000,
    "rate_per_kg": 865.36,
    "created_at": "2024-10-16T14:08:20.909822Z",
    "updated_at": "2024-10-16T14:08:20.909822Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  }
]
```

### Get All Jari Products

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/product/jari`

**Response:**
```json
[
  {
    "id": 2,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "jari",
    "type": "organic",
    "quantity_in_kg": 2000,
    "rate_per_kg": 265.36,
    "jari_size": "medium",
    "created_at": "2024-10-16T14:07:48.748941Z",
    "updated_at": "2024-10-16T14:07:48.748941Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  },
  {
    "id": 1,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "jari",
    "type": "organic",
    "quantity_in_kg": 2500,
    "rate_per_kg": 200.36,
    "jari_size": "small",
    "created_at": "2024-10-16T14:07:27.694646Z",
    "updated_at": "2024-10-16T14:07:27.694646Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  }
]
```

### Get All Products of a Farmer

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/product/farmer/1`

**Response:**
```json
[
  {
    "id": 4,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "mushroom",
    "type": "non-organic",
    "quantity_in_kg": 5000,
    "rate_per_kg": 865.36,
    "created_at": "2024-10-16T14:08:43.942844Z",
    "updated_at": "2024-10-16T14:08:43.942844Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  },
  {
    "id": 3,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "mushroom",
    "type": "organic",
    "quantity_in_kg": 1000,
    "rate_per_kg": 865.36,
    "created_at": "2024-10-16T14:08:20.909822Z",
    "updated_at": "2024-10-16T14:08:20.909822Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  },
  {
    "id": 2,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "jari",
    "type": "organic",
    "quantity_in_kg": 2000,
    "rate_per_kg": 265.36,
    "jari_size": "medium",
    "created_at": "2024-10-16T14:07:48.748941Z",
    "updated_at": "2024-10-16T14:07:48.748941Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  },
  {
    "id": 1,
    "img": "https://example.com/images/product_101.jpg",
    "farmer_id": 1,
    "name": "jari",
    "type": "organic",
    "quantity_in_kg": 2500,
    "rate_per_kg": 200.36,
    "jari_size": "small",
    "created_at": "2024-10-16T14:07:27.694646Z",
    "updated_at": "2024-10-16T14:07:27.694646Z",
    "farmer_phone_number": "9876543219",
    "farmers_first_name": "Rohan",
    "farmers_last_name": "Sharma",
    "is_available": false,
    "is_verified_by_admin": false
  }
]
```

### Get A Product By ID

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/product/4`

**Response:**
```json
{
  "id": 4,
  "img": "https://example.com/images/product_101.jpg",
  "farmer_id": 1,
  "name": "mushroom",
  "type": "non-organic",
  "quantity_in_kg": 5000,
  "rate_per_kg": 865.36,
  "created_at": "2024-10-16T14:08:43.942844Z",
  "updated_at": "2024-10-16T14:08:43.942844Z",
  "farmer_phone_number": "9876543219",
  "farmers_first_name": "Rohan",
  "farmers_last_name": "Sharma",
  "is_available": false,
  "is_verified_by_admin": false
}
```

### Update Product Unavailability

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/product/4/mark-unavailable`

**Response:**
```json
{
  "message": "Product availability updated successfully!"
}
```

### Delete Product

**Request:**
- Method: `DELETE`
- URL: `http://localhost:8080/api/v1/product/4`

**Response:**
```json
{
  "message": "product deleted successfully!"
}
```

## Order API

### Create Order

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/v1/product/3/order`
- Body:
```json
{
  "quantity_in_kg": 200,
  "delivery_address": "123 Maple St",
  "delivery_city": "Springfield",
  "delivery_address_zip": 62701,
  "mode_of_delivery": "Standard Shipping",