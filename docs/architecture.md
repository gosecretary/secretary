# Secretary Project Architecture

## System Architecture Diagram

This diagram illustrates the high-level architecture of the Secretary system, including the API gateway, session management, ephemeral credential generation, access request workflow, and connections to target resources.

```mermaid
graph TD
    User["User"] -- "HTTP/S" --> API
    Admin["Administrator"] -- "HTTP/S" --> API
    
    subgraph "Secretary Gateway"
        API["API Gateway"] --> Auth["Authentication"]
        API --> RBAC["Role-Based Access Control"]
        API --> Sessions["Session Management"]
        API --> Audit["Audit Logging"]
        API --> EphCreds["Ephemeral Credential\nGenerator"]
        API --> ReqWF["Access Request\nWorkflow"]
    end
    
    Auth --> DB[("Database")]
    RBAC --> DB
    Sessions --> DB
    EphCreds --> DB
    ReqWF --> DB
    Audit --> AuditLogs[("Audit Logs")]
    
    subgraph "Target Resources"
        EphCreds -. "temporary\ncredentials" .-> SSH["SSH Servers"]
        EphCreds -. "temporary\ncredentials" .-> DB_Servers["Database Servers"]
        EphCreds -. "temporary\ncredentials" .-> Other["Other Resources"]
    end
    
    Sessions -. "monitors\n& records" .-> SSH
    Sessions -. "monitors\n& records" .-> DB_Servers
    Sessions -. "monitors\n& records" .-> Other
    
    Admin -- "approves/denies" --> ReqWF
```

## Domain Model Class Diagram

This diagram shows the main domain models and their relationships in the Secretary system.

```mermaid
classDiagram
    class User {
        +String ID
        +String Username
        +String Email
        +String Password
        +String Name
        +String Role
        +DateTime CreatedAt
        +DateTime UpdatedAt
    }
    
    class Resource {
        +String ID
        +String Name
        +String Description
        +String Type
        +DateTime CreatedAt
        +DateTime UpdatedAt
    }
    
    class Credential {
        +String ID
        +String ResourceID
        +String Username
        +String Password
        +DateTime CreatedAt
        +DateTime UpdatedAt
    }
    
    class Permission {
        +String ID
        +String UserID
        +String ResourceID
        +String Role
        +String Action
        +DateTime CreatedAt
        +DateTime UpdatedAt
    }
    
    class Session {
        +String ID
        +String UserID
        +String ResourceID
        +DateTime StartTime
        +DateTime EndTime
        +String Status
        +String ClientIP
        +String ClientMetadata
        +String AuditPath
        +DateTime CreatedAt
        +DateTime UpdatedAt
    }
    
    class AccessRequest {
        +String ID
        +String UserID
        +String ResourceID
        +String Reason
        +String Status
        +String ReviewerID
        +String ReviewNotes
        +DateTime RequestedAt
        +DateTime ReviewedAt
        +DateTime ExpiresAt
        +DateTime CreatedAt
        +DateTime UpdatedAt
    }
    
    class EphemeralCredential {
        +String ID
        +String UserID
        +String ResourceID
        +String Username
        +String Password
        +String Token
        +DateTime ExpiresAt
        +DateTime CreatedAt
        +DateTime UsedAt
    }
    
    User "1" -- "*" Permission : has
    User "1" -- "*" Session : creates
    User "1" -- "*" AccessRequest : submits
    User "1" -- "*" AccessRequest : reviews
    User "1" -- "*" EphemeralCredential : receives
    
    Resource "1" -- "*" Credential : has
    Resource "1" -- "*" Permission : grants
    Resource "1" -- "*" Session : hosts
    Resource "1" -- "*" AccessRequest : target of
    Resource "1" -- "*" EphemeralCredential : provides access to
``` 