# Farm-to-Market Backend Design Document

## 1. Overview

This document outlines the backend design for a farm-to-market application that connects farmers with buyers. The application allows farmers to sell their crops and buyers to place orders. It includes user authentication, admin verification for farmers, and a marketplace for crop listings.

## 2. Architecture

We'll use a microservices architecture to ensure scalability and maintainability. The backend will be built using Go with the Echo framework. Here's a high-level overview of the architecture:

[FlowChart](https://claude.site/artifacts/fe1165ad-ef01-4598-94fa-87e7e809fd4a)


```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   API Gateway   │◄────┤  Load Balancer  │◄────┤   Client Apps   │
└────────┬────────┘     └─────────────────┘     └─────────────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Auth Service   │     │   User Service  │     │  Admin Service  │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                      │                        │
         │                      │                        │
         ▼                      ▼                        ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ Product Service │     │  Order Service  │     │ Notification    │
└─────────────────┘     └─────────────────┘     │     Service     │
                                                └─────────────────┘
         │                      │                        │
         │                      │                        │
         ▼                      ▼                        ▼
┌─────────────────────────────────────────────────────────────────┐
│                           Database                              │
└─────────────────────────────────────────────────────────────────┘
```


```
flowchart TD
    Start([Start]) --> UserType{User Type?}
    UserType -->|Farmer| FarmerSignup[Farmer Signup]
    UserType -->|Buyer| BuyerSignup[Buyer Signup]

    FarmerSignup --> CollectFarmerInfo[Collect Farmer Info]
    CollectFarmerInfo --> UploadDocuments[Upload Documents]
    UploadDocuments --> SendOTP[Send OTP for Verification]
    SendOTP --> VerifyOTP{Verify OTP}
    VerifyOTP -->|Success| CreateFarmerAccount[Create Farmer Account]
    VerifyOTP -->|Failure| RetryOTP[Retry OTP]
    RetryOTP --> SendOTP
    CreateFarmerAccount --> AdminReview{Admin Review}
    AdminReview -->|Approved| NotifyFarmer[Notify Farmer - Account Approved]
    AdminReview -->|Rejected| NotifyFarmerRejection[Notify Farmer - Account Rejected]
    NotifyFarmer --> FarmerLogin[Farmer Login]
    FarmerLogin --> FarmerDashboard[Farmer Dashboard]

    FarmerDashboard --> ListCrop[List Crop for Sale]
    ListCrop --> EnterCropDetails[Enter Crop Details]
    EnterCropDetails --> UploadCropPhotos[Upload Crop Photos]
    UploadCropPhotos --> SubmitForApproval[Submit for Approval]
    SubmitForApproval --> AdminReviewCrop{Admin Review Crop}
    AdminReviewCrop -->|Approved| ListCropOnMarket[List Crop on Market]
    AdminReviewCrop -->|Rejected| NotifyFarmerCropRejected[Notify Farmer - Crop Rejected]

    BuyerSignup --> CollectBuyerInfo[Collect Buyer Info]
    CollectBuyerInfo --> SendBuyerOTP[Send OTP for Verification]
    SendBuyerOTP --> VerifyBuyerOTP{Verify OTP}
    VerifyBuyerOTP -->|Success| CreateBuyerAccount[Create Buyer Account]
    VerifyBuyerOTP -->|Failure| RetryBuyerOTP[Retry OTP]
    RetryBuyerOTP --> SendBuyerOTP
    CreateBuyerAccount --> BuyerLogin[Buyer Login]
    BuyerLogin --> BuyerDashboard[Buyer Dashboard]

    BuyerDashboard --> BrowseCrops[Browse Available Crops]
    BrowseCrops --> SelectCrop[Select Crop]
    SelectCrop --> PlaceOrder[Place Order]
    PlaceOrder --> EnterOrderDetails[Enter Order Details]
    EnterOrderDetails --> ConfirmOrder[Confirm Order]
    ConfirmOrder --> NotifyFarmerOrder[Notify Farmer of New Order]
    NotifyFarmerOrder --> ConnectBuyerFarmer[Connect Buyer and Farmer]

    subgraph Admin Panel
    AdminReview
    AdminReviewCrop
    end

    subgraph Farmer Flow
    FarmerSignup
    FarmerLogin
    FarmerDashboard
    ListCrop
    end

    subgraph Buyer Flow
    BuyerSignup
    BuyerLogin
    BuyerDashboard
    BrowseCrops
    end
```

JWT DIAGRAM:
```
sequenceDiagram
    participant Client
    participant JWTMiddleware
    participant ExtractUserID
    participant IsFarmer
    participant RequestHandler

    Client->>JWTMiddleware: Request with JWT
    JWTMiddleware->>ExtractUserID: Validated token
    ExtractUserID->>ExtractUserID: Extract user_id
    ExtractUserID->>IsFarmer: Request with user_id in context
    IsFarmer->>IsFarmer: Check user_type
    IsFarmer->>RequestHandler: Request (if user is farmer)
    RequestHandler->>RequestHandler: Handle request (user_id available)
```


### Services:

1. **API Gateway**: Routes requests to appropriate microservices.
2. **Auth Service**: Handles user authentication and authorization.
3. **User Service**: Manages user profiles and registration.
4. **Admin Service**: Handles farmer verification and product approval.
5. **Product Service**: Manages crop listings.
6. **Order Service**: Handles order creation and management.
7. **Notification Service**: Sends OTPs and notifications.

## 3. Database Selection

For this application, I recommend using a SQL database, specifically PostgreSQL. Here's why:

1. **Data Relationships**: Your application has clear relationships between entities (users, products, orders), which SQL databases handle well.
2. **ACID Compliance**: SQL databases ensure data integrity, which is crucial for financial transactions and user data.
3. **Scalability**: While NoSQL databases are known for horizontal scalability, modern SQL databases like PostgreSQL can also scale well.
4. **Complex Queries**: SQL databases excel at complex joins and transactions, which your application might require.
5. **Consistency**: SQL databases provide strong consistency, which is important for your use case.

You mentioned using Supabase, which is an excellent choice. Supabase is built on top of PostgreSQL and provides additional features like real-time subscriptions, authentication, and storage.

## 4. Project Structure

Here's a suggested directory structure for your Go project:

```
farm-to-market/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── user/
│   ├── admin/
│   ├── product/
│   ├── order/
│   └── notification/
├── pkg/
│   ├── database/
│   ├── middleware/
│   └── utils/
├── api/
│   └── routes/
├── config/
├── scripts/
├── test/
└── go.mod
```

## 5. Key Components

### 5.1 User Authentication

- Implement JWT-based authentication.
- Use secure password hashing (e.g., bcrypt).
- Implement role-based access control (RBAC) for farmers, buyers, and admins.

### 5.2 User Registration

- Implement separate flows for farmers and buyers.
- Integrate with UIDAI for Aadhaar verification and OTP.
- Store user details securely in the database.

### 5.3 Admin Verification

- Create an admin dashboard for farmer verification.
- Implement document upload and review system.
- Use database transactions to ensure data consistency during approval.

### 5.4 Product Listing

- Allow farmers to create and manage crop listings.
- Implement admin approval workflow for listings.
- Use database triggers or application logic to handle status changes.

### 5.5 Order Management

- Implement order creation and tracking.
- Ensure proper validation of order quantities against available stock.

### 5.6 Notification System

- Integrate with SMS gateway for OTP and notifications.
- Implement email notifications as a fallback.

## 6. API Design

Design RESTful APIs for each service. Example endpoints:

- `/api/auth/register`
- `/api/auth/login`
- `/api/users/{id}`
- `/api/admin/verify-farmer/{id}`
- `/api/products`
- `/api/orders`

Use proper HTTP methods (GET, POST, PUT, DELETE) and status codes.

## 7. Security Considerations

- Implement rate limiting to prevent abuse.
- Use HTTPS for all communications.
- Sanitize and validate all user inputs.
- Implement proper error handling to avoid information leakage.
- Regular security audits and penetration testing.

## 8. Scalability

- Use container orchestration (e.g., Kubernetes) for easy scaling.
- Implement caching (e.g., Redis) for frequently accessed data.
- Use database read replicas for scaling read operations.
- Implement asynchronous processing for non-critical tasks.

## 9. Monitoring and Logging

- Implement comprehensive logging across all services.
- Use a centralized logging system (e.g., ELK stack).
- Set up monitoring and alerting (e.g., Prometheus and Grafana).

## 10. Testing

- Implement unit tests for all critical components.
- Use integration tests to ensure proper interaction between services.
- Perform load testing to ensure the system can handle the expected user load.

## 11. Deployment

- Use CI/CD pipelines for automated testing and deployment.
- Consider using a cloud provider (e.g., AWS, GCP, or Azure) for hosting.
- Implement blue-green or canary deployment strategies for zero-downtime updates.

## 12. Cost Optimization

- Use serverless functions for infrequently used features.
- Implement auto-scaling to handle varying loads efficiently.
- Use managed services where possible to reduce operational overhead.
- Optimize database queries and implement proper indexing.

## Conclusion

This design provides a scalable and maintainable architecture for your farm-to-market application. By using Go with the Echo framework and PostgreSQL (via Supabase), you'll have a solid foundation for building a robust system that can handle 10k+ daily users while keeping costs low.

Remember to iterate on this design as you develop and gather more specific requirements. Regular performance testing and optimization will be key to ensuring the system scales effectively.