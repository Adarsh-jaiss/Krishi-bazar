# Signup :

req : http://localhost:8080/api/auth/signup
```
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

req : http://localhost:8080/api/auth/complete-signup

```
{
  "user": {
    "first_name": "Rahul",
    "last_name": "Sharma",
    "aadhar_number": "123412341212",
    "email": "rahul.sharma@example.com",
    "phone_number": "6200059008",
    "is_farmer": true,
    "address": "123 Green Farm Lane",
    "city": "Jhabua",
    "state": "Madhya Pradesh",
    "pin_code": "456001",
    "farm_size": "3"
  },
  "verification_code": "406632"
}
```

# HOMEPAGE REq hEADER:

```
fetch('https://your-api.com/protected-route', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  },
})
```
# ADMIN APPROVE :

req: http://localhost:8080/api/admin/approve

```
{
 "user_id":"2",
 "is_verified":true
}
```