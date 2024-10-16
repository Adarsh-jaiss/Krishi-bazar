# Krishi Bazar Backend Design Document

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


## API DOC : [DOC]("API-doc.md")

### Services:

1. **Auth Service**: Handles user authentication and authorization.
2. **User Service**: Manages user profiles and registration.
3. **Admin Service**: Handles farmer verification and product approval.
4. **Product Service**: Manages crop listings.
5. **Order Service**: Handles order creation and management.

## 4. Project Structure

Here's a suggested directory structure for your Go project:

```
farm-to-market/
├── main.go
├── internal/
│   ├── auth/
│   ├── user/
│   ├── admin/
│   ├── product/
│   ├── order/
│   └── notification/
├── database/
├── types/
├── scripts/
├── test/
└── go.mod
```

