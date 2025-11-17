# Testa-Walt: Digital Identity Infrastructure for Agricultural Finance

## Executive Summary

**The Challenge**: Small-scale farmers in Kenya struggle to access formal financial services due to lack of verifiable identity credentials. Traditional KYC processes are expensive, time-consuming, and exclude those without formal identification documents.

**The Solution**: A complete Digital Public Infrastructure (DPI) implementation using W3C Verifiable Credentials, demonstrating how governments and SACCOs can issue and verify farmer identity credentials at scale.

**Impact Potential**:

- Reduce onboarding costs
- Enable instant identity verification (minutes vs weeks)
- Privacy-preserving with selective disclosure

## What This Project Demonstrates

This is a **production-ready reference implementation** show casing:

1. **Government Issuer (Testa Gava)** - Issues portable social security documents (A1 certificates)
2. **SACCO Verifier (Testa SACCO)** - Verifies credentials for account opening
3. **Complete DPI Stack** - Walt.id services deployed on cloud infrastructure

**Technical Stack**: Walt.id (SSI), OpenID4VC, W3C Verifiable Credentials, Go + HTMX, Docker, DigitalOcean

---

## Table of Contents

- [Testa-Walt: Digital Identity Infrastructure for Agricultural Finance](#testa-walt-digital-identity-infrastructure-for-agricultural-finance)
  - [Executive Summary](#executive-summary)
  - [What This Project Demonstrates](#what-this-project-demonstrates)
  - [Table of Contents](#table-of-contents)
  - [Business Context](#business-context)
    - [The Agricultural Finance Gap in Kenya](#the-agricultural-finance-gap-in-kenya)
    - [Why Verifiable Credentials?](#why-verifiable-credentials)
    - [Alignment with CDPI Principles](#alignment-with-cdpi-principles)
  - [Architecture Overview](#architecture-overview)
    - [System Components](#system-components)
    - [Technology Decisions](#technology-decisions)
    - [Data Flow](#data-flow)
  - [Key Terms](#key-terms)
  - [Server Deployment](#server-deployment)
    - [Prerequisites - Server](#prerequisites---server)
    - [Server Setup](#server-setup)
    - [Installation](#installation)
    - [Services Configuration](#services-configuration)
  - [Issuing W3C Verifiable Credentials](#issuing-w3c-verifiable-credentials)
    - [Issuing Context](#issuing-context)
    - [Issuing Workflow](#issuing-workflow)
    - [Issuing Step-by-Step: Behind The Scenes](#issuing-step-by-step-behind-the-scenes)
      - [1. Generate a Credential Offer](#1-generate-a-credential-offer)
      - [2. Citizen Logs Into Wallet](#2-citizen-logs-into-wallet)
      - [3. List Available Wallets](#3-list-available-wallets)
      - [4. Create a DID](#4-create-a-did)
      - [5. Accept the Credential Offer](#5-accept-the-credential-offer)
  - [Verifying W3C Verifiable Credentials](#verifying-w3c-verifiable-credentials)
    - [Verification Context](#verification-context)
    - [Verification Workflow](#verification-workflow)
    - [Verification Step-by-Step](#verification-step-by-step)
      - [1. SACCO Generates Presentation Request](#1-sacco-generates-presentation-request)
      - [2. Citizen Resolves the Presentation Request](#2-citizen-resolves-the-presentation-request)
      - [3. Citizen Shares the Credential](#3-citizen-shares-the-credential)
      - [4. SACCO Verifies the Credentials](#4-sacco-verifies-the-credentials)
  - [API Reference](#api-reference)
    - [Issuer API Endpoints](#issuer-api-endpoints)
    - [Wallet API Endpoints](#wallet-api-endpoints)
    - [Verifier API Endpoints](#verifier-api-endpoints)
  - [Troubleshooting](#troubleshooting)
  - [Testa Web App Deployment](#testa-web-app-deployment)
    - [Prerequisites - Web App](#prerequisites---web-app)
      - [Issuer](#issuer)
      - [Verifier](#verifier)
    - [Go + HTMX Stats](#go--htmx-stats)
  - [Lessons Learned \& Best Practices](#lessons-learned--best-practices)
    - [What Worked Well](#what-worked-well)
    - [Challenges \& Solutions](#challenges--solutions)
      - [Challenge 1: CORS Error](#challenge-1-cors-error)
      - [Challenge 2: Selective Disclosure Complexity](#challenge-2-selective-disclosure-complexity)
      - [Challenge 3: Mobile Form UX](#challenge-3-mobile-form-ux)

---

## Business Context

### The Agricultural Finance Gap in Kenya

Kenya has:

- **7.5 million+ smallholder farmers** producing over 75% of agricultural output per the farm to market alliance
- **5,000+ SACCOs** serving rural communities (178 Deposit taking and 177 Non deposit taking SASRA regulated)
- **Minimal financial inclusion** among rural farmers
- **High KYC costs** preventing SACCO expansion

### Why Verifiable Credentials?

Traditional identity verification requires:

- Physical document presentation
- Manual verification
- Paper records
- Multiple visits

**With Verifiable Credentials**:

```text
Farmer → Gets A1 cert from gov't → Stores in wallet → 
Presents to SACCO → Instant verification → Account opened
```

**Result**: Same-day account opening, lower costs, higher trust.

### Alignment with CDPI Principles

This implementation follows Centre for Digital Public Infrastructure best practices:

1. **Open Standards**: W3C VC, OpenID4VC, DID methods
2. **Interoperability**: Works with any W3C-compliant wallet
3. **Privacy by Design**: Selective disclosure, minimal data sharing
4. **Inclusive Design**: USSD support, mobile-first UI
5. **Modularity**: Built for local context while maintaining extensibility

---

## Architecture Overview

### System Components

```text
┌─────────────────────────────────────────────────────────────┐
│                    Digital Public Infrastructure            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────┐   │
│  │ Testa Gava   │      │   Wallet     │      │  Testa   │   │
│  │  (Issuer)    │──────│  (Holder)    │──────│  SACCO   │   │
│  │  Port 8082   │      │  Port 7101   │      │(Verifier)│   │
│  │              │      │              │      │Port 8081 │   │
│  └──────┬───────┘      └──────────────┘      └────┬─────┘   │
│         │                                         │         │
│         │         Walt.id Services                │         │
│  ┌──────▼─────────────────────────────────────────▼────┐    │
│  │  Issuer API (7002)  │  Verifier API (7003)          │    │
│  │  Wallet API (7001)  │  VC Repository (7103)         │    │
│  │  Portal API (7102)  │                               │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Technology Decisions

| Component | Technology | Why? |
|-----------|-----------|------|
| **SSI Framework** | Walt.id | Open-source, production-ready, OpenID4VC support |
| **Application Layer** | Go + HTMX | Fast (TTFB <1ms), server-side rendering, low bandwidth |
| **Deployment** | Docker + DigitalOcean | Cost-effective ($28/mo), scalable data centers |
| **Credential Format** | SD-JWT | Selective disclosure, mobile-friendly, JSON-based |
| **DID Method** | did:key | Simple, no blockchain required, instant resolution |

### Data Flow

**Issuance Flow**:

1. Farmer visits Testa Gava web app
2. Submits KYC information (A1 form)
3. App calls Walt.id Issuer API
4. Returns credential offer link (openid-credential://...)
5. Farmer scans QR code or clicks link
6. Wallet imports credential

**Verification Flow**:

1. SACCO generates verification request
2. Farmer receives request (QR/link)
3. Wallet prompts consent
4. Farmer selects what to share (selective disclosure)
5. Wallet presents credential
6. SACCO verifies automatically
7. Account opened

## Key Terms

Walt.id provides a comprehensive set of tools to implement verifiable credentials according to the W3C standards. This system consists of three main components:

1. **Issuer** - Creates and signs verifiable credentials
2. **Wallet** - Stores and manages credentials for the credential holder
3. **Verifier** - Requests and verifies credentials from holders

This documentation covers the deployment of walt.id services and a complete implementation of both the issuance and verification workflows.

## Server Deployment

### Prerequisites - Server

- DigitalOcean account
- Basic familiarity with Linux commands and Docker
- SSH key pair for secure access

### Server Setup

1. **Create a DigitalOcean Droplet**
   - Log in to your DigitalOcean account
   - Click on "Create" and select "Droplets"
   - Choose a datacenter region closest to your users
   - Select Ubuntu 24.04 LTS (Noble) or the most recent LTS
   - Choose droplet size:
     - Recommended: 4 GB RAM / 2 AMD CPUs / 80 GB NVMe SSD / 4 TB transfer ($28/mo)
   - Add your SSH key
   - Choose a hostname (e.g., `walt-id-test-server`)
   - Click "Create Droplet"

2. **Connect to Your Droplet**

   ```bash
   ssh root@your_droplet_ip
   ```

### Installation

1. **Install Required Dependencies**

   ```bash
   apt update && apt upgrade -y
   apt install git docker-compose -y
   ```

2. **Clone the Repository**

   ```bash
   git clone https://github.com/walt-id/waltid-identity.git && cd waltid-identity
   ```

3. **Launch Services**

   ```bash
   cd docker-compose && docker compose up -d
   ```

   Depending on your docker-compose version you may use `docker-compose up -d`

   You can ommit the `-d` flag if you want to watch the logs

   **Note:** If you encounter a "pull_policy: missing" error, update the docker-compose.yaml file:

   ```bash
   sed -i 's/pull_policy: missing/pull_policy: if_not_present/g' docker-compose.yaml
   ```

### Services Configuration

The following services are exposed on these ports:

**APIs:**

- Wallet API: `http://your_server_ip:7001`
- Issuer API: `http://your_server_ip:7002`
- Verifier API: `http://your_server_ip:7003`

**Note**: To directly test the API you can run the [postman collection](/postmanCollection.json)

**Web Applications:**

- Demo Web Wallet: `http://your_server_ip:7101`
- Dev Web Wallet: `http://your_server_ip:7104`
- Web Portal: `http://your_server_ip:7102`
- Credential Repo: `http://your_server_ip:7103`

## Issuing W3C Verifiable Credentials

### Issuing Context

In this implementation, we simulate a government agency issuing identity credentials to citizens:

1. A citizen visits a government website
2. The citizen logs into their account
3. The government offers digital credentials (ID, passport, etc.)
4. The citizen must accept the credential offer to receive it in their wallet
5. The credential is securely stored in their digital wallet for future use

### Issuing Workflow

1. Government generates a credential offer
2. Citizen authenticates to their wallet
3. Citizen creates a DID (Decentralized Identifier)
4. Citizen scans the credential offer (via QR code)
5. Citizen accepts the credential
6. Credential is stored in the citizen's wallet

### Issuing Step-by-Step: Behind The Scenes

#### 1. Generate a Credential Offer

The government service generates a credential offer URL by calling the `/openid4vc/jwt/issue` or `/openid4vc/sdjwt/issue` endpoint for ordinary or selective diclosure issuance respectively:

```bash
curl -X 'POST' \
'http://your_server_ip:7002/openid4vc/jwt/issue' \
-H 'accept: text/plain' \
-H 'Content-Type: application/json' \
-d '{
  "issuerKey": {
    "type": "jwk",
    "jwk": {
      "kty": "OKP",
      "d": "mDhpwaH6JYSrD2Bq7Cs-pzmsjlLj4EOhxyI-9DM1mFI",
      "crv": "Ed25519",
      "kid": "Vzx7l5fh56F3Pf9aR3DECU5BwfrY6ZJe05aiWYWzan8",
      "x": "T3T4-u1Xz3vAV2JwPNxWfs4pik_JLiArz_WTCvrCFUM"
    }
  },
  "issuerDid": "did:key:z6MkjoRhq1jSNJdLiruSXrFFxagqrztZaXHqHGUTKJbcNywp",
  "credentialConfigurationId": "OpenBadgeCredential_jwt_vc_json",
  "credentialData": {
    "@context": [
      "https://www.w3.org/2018/credentials/v1",
      "https://purl.imsglobal.org/spec/ob/v3p0/context.json"
    ],
    "id": "urn:uuid:THIS WILL BE REPLACED WITH DYNAMIC DATA FUNCTION",
    "type": [
      "VerifiableCredential",
      "OpenBadgeCredential"
    ],
    "name": "JFF x vc-edu PlugFest 3 Interoperability",
    "issuer": {
      "type": [
        "Profile"
      ],
      "id": "did:key:THIS WILL BE REPLACED WITH DYNAMIC DATA FUNCTION FROM CONTEXT",
      "name": "Jobs for the Future (JFF)",
      "url": "https://www.jff.org/",
      "image": "https://w3c-ccg.github.io/vc-ed/plugfest-1-2022/images/JFF_LogoLockup.png"
    },
    "credentialSubject": {
      "id": "did:key:123 (THIS WILL BE REPLACED BY DYNAMIC DATA FUNCTION)",
      "type": [
        "AchievementSubject"
      ],
      "achievement": {
        "id": "urn:uuid:ac254bd5-8fad-4bb1-9d29-efd938536926",
        "type": [
          "Achievement"
        ],
        "name": "JFF x vc-edu PlugFest 3 Interoperability",
        "description": "This wallet supports the use of W3C Verifiable Credentials and has demonstrated interoperability during the presentation request workflow during JFF x VC-EDU PlugFest 3.",
        "criteria": {
          "type": "Criteria",
          "narrative": "Wallet solutions providers earned this badge by demonstrating interoperability during the presentation request workflow. This includes successfully receiving a presentation request, allowing the holder to select at least two types of verifiable credentials to create a verifiable presentation, returning the presentation to the requestor, and passing verification of the presentation and the included credentials."
        },
        "image": {
          "id": "https://w3c-ccg.github.io/vc-ed/plugfest-3-2023/images/JFF-VC-EDU-PLUGFEST3-badge-image.png",
          "type": "Image"
        }
      }
    }
  },
  "mapping": {
    "id": "<uuid>",
    "issuer": {
      "id": "<issuerDid>"
    },
    "credentialSubject": {
      "id": "<subjectDid>"
    },
    "issuanceDate": "<timestamp>",
    "expirationDate": "<timestamp-in:365d>"
  }
}'
```

or with sdjwt:

```bash
curl -X 'POST' \
  'http://your_server_ip:7002/openid4vc/sdjwt/issue' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "issuerKey": {
    "type": "jwk",
    "jwk": {
      "kty": "OKP",
      "d": "mDhpwaH6JYSrD2Bq7Cs-pzmsjlLj4EOhxyI-9DM1mFI",
      "crv": "Ed25519",
      "kid": "Vzx7l5fh56F3Pf9aR3DECU5BwfrY6ZJe05aiWYWzan8",
      "x": "T3T4-u1Xz3vAV2JwPNxWfs4pik_JLiArz_WTCvrCFUM"
    }
  },
  "credentialConfigurationId": "OpenBadgeCredential_jwt_vc_json",
  "credentialData": {
    "@context": [
      "https://www.w3.org/2018/credentials/v1",
      "https://purl.imsglobal.org/spec/ob/v3p0/context.json"
    ],
    "id": "urn:uuid:THIS WILL BE REPLACED WITH DYNAMIC DATA FUNCTION (see below)",
    "type": [
      "VerifiableCredential",
      "OpenBadgeCredential"
    ],
    "name": "JFF x vc-edu PlugFest 3 Interoperability",
    "issuer": {
      "type": [
        "Profile"
      ],
      "name": "Jobs for the Future (JFF)",
      "url": "https://www.jff.org/",
      "image": "https://w3c-ccg.github.io/vc-ed/plugfest-1-2022/images/JFF_LogoLockup.png"
    },
    "credentialSubject": {
      "type": [
        "AchievementSubject"
      ],
      "achievement": {
        "id": "urn:uuid:ac254bd5-8fad-4bb1-9d29-efd938536926",
        "type": [
          "Achievement"
        ],
        "name": "JFF x vc-edu PlugFest 3 Interoperability",
        "description": "This wallet supports the use of W3C Verifiable Credentials and has demonstrated interoperability during the presentation request workflow during JFF x VC-EDU PlugFest 3.",
        "criteria": {
          "type": "Criteria",
          "narrative": "Wallet solutions providers earned this badge by demonstrating interoperability during the presentation request workflow. This includes successfully receiving a presentation request, allowing the holder to select at least two types of verifiable credentials to create a verifiable presentation, returning the presentation to the requestor, and passing verification of the presentation and the included credentials."
        },
        "image": {
          "id": "https://w3c-ccg.github.io/vc-ed/plugfest-3-2023/images/JFF-VC-EDU-PLUGFEST3-badge-image.png",
          "type": "Image"
        }
      }
    }
  },
  "mapping": {
    "id": "<uuid>",
    "issuer": {
      "id": "<issuerDid>"
    },
    "credentialSubject": {
      "id": "<subjectDid>"
    },
    "issuanceDate": "<timestamp>",
    "expirationDate": "<timestamp-in:365d>"
  },
  "selectiveDisclosure": {
    "fields": {
      "name": {
        "sd": true
      }
    }
  },
  "issuerDid": "did:key:z6MkjoRhq1jSNJdLiruSXrFFxagqrztZaXHqHGUTKJbcNywp"
}'
```

The result is a credential offer URL, for example with the `/openid4vc/jwt/issue` endpoint:

```text
openid-credential-offer://issuer.demo.walt.id/draft13/?credential_offer_uri=https%3A%2F%2Fissuer.demo.walt.id%2Fdraft13%2FcredentialOffer%3Fid%3D8a2c0ce8-a2e3-43e4-9225-164432bcd76e
```

#### 2. Citizen Logs Into Wallet

The citizen logs into their wallet:

```bash
curl -X 'POST' \
  'http://your_server_ip/wallet-api/auth/login' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "type": "email",
  "email": "user@email.com",
  "password": "password"
}'
```

Response:

```json
{
  "id": "8a26f0a5-e51a-40eb-8da7-2bfb999b6080",
  "token": "eyJhbGciOiJIUzI1NiJ9.eyJuYmY...",
  "username": "user@email.com"
}
```

Save the token for future API calls:

```bash
TOKEN=eyJhbGciOiJIUzI1NiJ9.eyJuYmY...
```

#### 3. List Available Wallets

The citizen checks which wallets are available:

```bash
curl -X 'GET' \
  'http://your_server_ip/wallet-api/wallet/accounts/wallets' \
  -H 'accept: application/json' \
  -H "authorization: Bearer $TOKEN"
```

Response:

```json
{
  "account":"8a26f0a5-e51a-40eb-8da7-2bfb999b6080",
  "wallets":[
    {
      "id":"3cf48671-4b91-4fc3-9d79-496c1c0ba91b",
      "name":"Wallet of Max Mustermann",
      "createdOn":"2025-03-03T12:22:22.811Z",
      "addedOn":"2025-03-03T12:22:22.811Z",
      "permission":"ADMINISTRATE"
    }
  ]
}
```

#### 4. Create a DID

The citizen creates a decentralized identifier (DID):

```bash
curl -X 'POST' \
  'http://your_server_ip/wallet-api/wallet/3cf48671-4b91-4fc3-9d79-496c1c0ba91b/dids/create/key' \
  -H 'accept: */*' \
  -H "authorization: Bearer $TOKEN" \
  -d ''
```

Response:

```text
did:key:zDnaeiXE5gkKfmaUzN7e1MANvMj8NJg4W7KtWfxLXnA1i4Zh6
```

#### 5. Accept the Credential Offer

The citizen accepts the credential offer:

```bash
curl -X 'POST' \
  'http://your_server_ip/wallet-api/wallet/3cf48671-4b91-4fc3-9d79-496c1c0ba91b/exchange/useOfferRequest?did=did%3Akey%3AzDnaeiXE5gkKfmaUzN7e1MANvMj8NJg4W7KtWfxLXnA1i4Zh6' \
  -H 'accept: application/json' \
  -H 'Content-Type: text/plain' \
  -H "authorization: Bearer $TOKEN" \
  -d 'openid-credential-offer://issuer.demo.walt.id/draft13/?credential_offer_uri=https%3A%2F%2Fissuer.demo.walt.id%2Fdraft13%2FcredentialOffer%3Fid%3D8a2c0ce8-a2e3-43e4-9225-164432bcd76e'
```

The response contains the issued credential, including an ID to reference it later:

```json
[{
  "wallet": "3cf48671-4b91-4fc3-9d79-496c1c0ba91b",
  "id": "urn:uuid:3d712efa-1a30-4625-80ce-e481fcb907e7",
  "document": "eyJraWQiOiJkaWQ6a2V5On...",
  ...etc.
}]
```

## Verifying W3C Verifiable Credentials

### Verification Context

In this implementation, we simulate a financial institution (SACCO) verifying a citizen's identity:

1. The citizen, a farmer, wants to open an account at a SACCO
2. The SACCO following regulations needs to verify the farmer's identity
3. The citizen presents their government-issued credential
4. The SACCO verifies the credential and opens the account

### Verification Workflow

1. SACCO generates a presentation request
2. Farmer receives the request (via QR code or link)
3. Farmer resolves the presentation request
4. Farmer selects and shares the requested credential
5. SACCO verifies the presented credential
6. SACCO checks verification results and policies

### Verification Step-by-Step

#### 1. SACCO Generates Presentation Request

The SACCO generates a presentation request:

```bash
curl -X 'POST' \
  'http://your_server_ip:7003/openid4vc/verify' \
  -H 'accept: */*' \
  -H 'authorizeBaseUrl: openid4vp://authorize' \
  -H 'responseMode: direct_post' \
  -H 'Content-Type: application/json' \
  -d '{
  "request_credentials": [
    {
      "type": "OpenBadgeCredential",
      "format": "jwt_vc_json"
    }
  ]
}'
```

Response (presentation request URL):

```text
openid4vp://authorize?response_type=vp_token&client_id=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify&response_mode=direct_post&state=ke8eFZteF7RU&presentation_definition_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fpd%2Fke8eFZteF7RU&client_id_scheme=redirect_uri&client_metadata=%7B%22authorization_encrypted_response_alg%22%3A%22ECDH-ES%22%2C%22authorization_encrypted_response_enc%22%3A%22A256GCM%22%7D&nonce=8e615c90-d6b7-48d5-bee9-7c1ba8208b4c&response_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify%2Fke8eFZteF7RU
```

**Important**: Extract the `state` parameter value (e.g., `ke8eFZteF7RU`) for later use.

#### 2. Citizen Resolves the Presentation Request

The citizen resolves the presentation request URL:

```bash
curl -X 'POST' \
  'http://your_server_ip:7003/wallet-api/wallet/3cf48671-4b91-4fc3-9d79-496c1c0ba91b/exchange/resolvePresentationRequest' \
  -H 'accept: text/plain' \
  -H 'Content-Type: text/plain' \
  -H "authorization: Bearer $TOKEN" \
  -d 'openid4vp://authorize?response_type=vp_token&client_id=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify&response_mode=direct_post&state=ke8eFZteF7RU&presentation_definition_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fpd%2Fke8eFZteF7RU&client_id_scheme=redirect_uri&client_metadata=%7B%22authorization_encrypted_response_alg%22%3A%22ECDH-ES%22%2C%22authorization_encrypted_response_enc%22%3A%22A256GCM%22%7D&nonce=8e615c90-d6b7-48d5-bee9-7c1ba8208b4c&response_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify%2Fke8eFZteF7RU'
```

Response (resolved URL):

```text
openid4vp://authorize?response_type=vp_token&client_id=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify&response_mode=direct_post&state=ke8eFZteF7RU&presentation_definition=%7B%22id%22%3A%22eMGijqmoddAK%22%2C%22input_descriptors%22%3A%5B%7B%22id%22%3A%22OpenBadgeCredential%22%2C%22format%22%3A%7B%22jwt_vc_json%22%3A%7B%22alg%22%3A%5B%22EdDSA%22%5D%7D%7D%2C%22constraints%22%3A%7B%22fields%22%3A%5B%7B%22path%22%3A%5B%22%24.vc.type%22%5D%2C%22filter%22%3A%7B%22type%22%3A%22string%22%2C%22pattern%22%3A%22OpenBadgeCredential%22%7D%7D%5D%7D%7D%5D%7D&presentation_definition_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fpd%2Fke8eFZteF7RU&client_id_scheme=redirect_uri&client_metadata=%7B%22authorization_encrypted_response_alg%22%3A%22ECDH-ES%22%2C%22authorization_encrypted_response_enc%22%3A%22A256GCM%22%7D&nonce=8e615c90-d6b7-48d5-bee9-7c1ba8208b4c&response_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify%2Fke8eFZteF7RU
```

#### 3. Citizen Shares the Credential

The citizen selects and shares the credential:

```bash
curl -X 'POST' \
  'http://your_server_ip:7003/wallet-api/wallet/3cf48671-4b91-4fc3-9d79-496c1c0ba91b/exchange/usePresentationRequest' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -H "authorization: Bearer $TOKEN" \
  -d '{
  "did": "did:key:zDnaeiXE5gkKfmaUzN7e1MANvMj8NJg4W7KtWfxLXnA1i4Zh6",
  "presentationRequest": "openid4vp://authorize?response_type=vp_token&client_id=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify&response_mode=direct_post&state=ke8eFZteF7RU&presentation_definition=%7B%22id%22%3A%22eMGijqmoddAK%22%2C%22input_descriptors%22%3A%5B%7B%22id%22%3A%22OpenBadgeCredential%22%2C%22format%22%3A%7B%22jwt_vc_json%22%3A%7B%22alg%22%3A%5B%22EdDSA%22%5D%7D%7D%2C%22constraints%22%3A%7B%22fields%22%3A%5B%7B%22path%22%3A%5B%22%24.vc.type%22%5D%2C%22filter%22%3A%7B%22type%22%3A%22string%22%2C%22pattern%22%3A%22OpenBadgeCredential%22%7D%7D%5D%7D%7D%5D%7D&presentation_definition_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fpd%2Fke8eFZteF7RU&client_id_scheme=redirect_uri&client_metadata=%7B%22authorization_encrypted_response_alg%22%3A%22ECDH-ES%22%2C%22authorization_encrypted_response_enc%22%3A%22A256GCM%22%7D&nonce=8e615c90-d6b7-48d5-bee9-7c1ba8208b4c&response_uri=https%3A%2F%2Fverifier.demo.walt.id%2Fopenid4vc%2Fverify%2Fke8eFZteF7RU",
  "selectedCredentials": [
    "urn:uuid:3d712efa-1a30-4625-80ce-e481fcb907e7"
  ]
}'
```

Response (if successful):

```json
{"redirectUri": null}
```

A 200 status code with this response indicates the presentation was accepted.

#### 4. SACCO Verifies the Credentials

The SACCO checks verification results using the session ID from step 1:

```bash
curl -X 'GET' \
  'http://your_server_ip:7003/openid4vc/session/ke8eFZteF7RU' \
  -H 'accept: */*'
```

The response contains the detailed verification results, including:

- Credential validity
- Digital signature verification
- Policy compliance
- Credential content

A `"verificationResult": true` in the response indicates successful verification.

## API Reference

### Issuer API Endpoints

- **POST /openid4vc/jwt/issue** - Issue a new credential and get an offer URL
- **POST /openid4vc/sdjwt/issue** - Issue a new credential with selective disclosure and get an offer URL
- **GET /draft13/credentialOffer** - Retrieve a credential offer by ID

### Wallet API Endpoints

- **POST /wallet-api/auth/login** - Authenticate a user
- **GET /wallet-api/wallet/accounts/wallets** - List wallets for an account
- **POST /wallet-api/wallet/{id}/dids/create/key** - Create a new DID
- **POST /wallet-api/wallet/{id}/exchange/useOfferRequest** - Accept a credential offer
- **POST /wallet-api/wallet/{id}/exchange/resolvePresentationRequest** - Resolve a presentation request
- **POST /wallet-api/wallet/{id}/exchange/usePresentationRequest** - Share credentials in response to a request

### Verifier API Endpoints

- **POST /openid4vc/verify** - Create a presentation request
- **GET /openid4vc/session/{id}** - Check verification results

## Troubleshooting

- **Session Expiration**: OpenID4VC sessions expire quickly. Complete the verification flow within 1-2 minutes of generating the presentation request.
- **Authentication Issues**: Ensure the authorization token is correctly included in all requests.
- **Docker Compose Issues**: If you encounter pull policy errors, update the docker-compose.yaml file as noted in the [installation](#installation) section.
- **Formatting Errors**: When using cURL, ensure proper escaping of special characters in URLs and JSON payloads.

For more information, refer to the [walt.id documentation](https://docs.walt.id/community-stack/guides/issue-verify-w3c-credential/).

## Testa Web App Deployment

### Prerequisites - Web App

Basic familiarity with Linux commands, Docker, Go 1.24+ and HTMX

#### Issuer

Read the full documentation [in the README](/issuer/README.md)

#### Verifier

Read the full documentation [in the README](/verifier/README.md)

### Go + HTMX Stats

Time to First Byte (TTFB):

```bash
curl -w "@-" -o /dev/null -s http://139.59.15.151/:8080 << 'EOF'

time_total: %{time_total}s

time_starttransfer: %{time_starttransfer}s

EOF
```

Response:

```bash
time_total: 0.000197s

time_starttransfer: 0.000000s
```

## Lessons Learned & Best Practices

### What Worked Well

1. **Go + HTMX Stack**:

   - Fast development (built both apps in 1 week)
   - Extremely low bandwidth (<30KB per page load)
   - Server-side rendering = SEO friendly
   - Perfect for Kenya's mobile-first, high-latency environment

2. **Docker Deployment**:

   - Easy to replicate across environments
   - Consistent behavior dev → staging → production
   - Simple updates with docker-compose pull

3. **Walt.id Choice**:

   - Well-documented APIs
   - Active community support
   - OpenID4VC compliance out of box
   - No blockchain dependency = lower costs

### Challenges & Solutions

#### Challenge 1: CORS Error

**Problem**: Portal tried calling localhost:7103 instead of the server IP on the same port
**Solution**: Set `SERVICE_HOST=server_ip` in docker-compose .env file
**Lesson**: Always configure public-facing URLs for production

#### Challenge 2: Selective Disclosure Complexity

**Problem**: 40+ fields in A1 certificate need selective disclosure config
**Solution**: Structured selective disclosure configurations in `models/credential.go`
**Lesson**: Model complex credentials carefully upfront

#### Challenge 3: Mobile Form UX

**Problem**: Long A1 form overwhelming on mobile
**Solution**:

- Collapsible sections (planned)
- Progress indicator (planned)
- Save as draft (planned)
**Lesson**: Design for 3G connections and small screens
