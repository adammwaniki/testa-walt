# Testa SACCO - Credential Verifier

Production-ready Go + HTMX web application for verifying W3C Verifiable Credentials using OpenID4VP protocol.

## Overview

Testa SACCO provides secure credential verification services for SACCO members and partners. Built with Go and HTMX for optimal performance and server-side rendering.

## Key Features

- **OpenID4VP Verification**: Standards-compliant credential verification
- **Dynamic Configuration**: Web form for verification policies
- **Selective Policies**: Choose which verification checks to perform
- **Server-Side Rendering**: Fast, SEO-friendly pages
- **Production Ready**: Modular architecture with proper separation of concerns

## Project Structure

```text
testa-sacco/
├── main.go                    # Server configuration
├── handlers/
│   └── handler.go            # All HTTP handlers
├── models/
│   └── verification.go       # Data structures
├── templates/
│   └── index.html            # Verification form
├── static/
│   └── styles.css            # Green and Blue themed CSS (blue-collar jobs in agriculture making money)
├── go.mod                     # Go module
├── Dockerfile                 # Container config
├── docker-compose.yml         # Easy deployment
├── Makefile                   # Dev shortcuts
└── README.md                  # This file
```

## Quick Start

```bash
# 1. Navigate to project
cd testa-sacco

# 2. Run the server
go run main.go

# 3. Visit http://localhost:8081
```

## Verification Options

The web form allows you to configure:

### Credential Type

- Verifiable Portable Document A1 (Default)
- Verifiable Attestation
- University Degree Credential
- Permanent Resident Card
- Open Badge Credential

### Verification Policies

- **Verify Signature** - Confirms credential hasn't been tampered with
- **Check Expiration** - Verifies credential is still valid
- **Check Not-Before** - Ensures credential is currently active
- **Check Revocation Status** - Confirms credential hasn't been revoked

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `WALTID_VERIFIER_URL` | `http://139.59.15.151:7003/openid4vc/verify` | Walt.id verifier endpoint |
| `PORT` | `8081` | Server port |

## Architecture

### main.go

Simple server setup with routing - just 30 lines!

### handlers/handler.go

All HTTP logic:

- `Home()` - Renders verification form
- `VerifyCredential()` - Processes form and creates verification request
- `extractVerificationOptions()` - Parses form data
- `buildVerificationRequest()` - Creates Walt.id request
- `renderSuccess()` / `renderError()` - HTMX responses

### models/verification.go

Complete data structures:

- `VerificationRequest` - Walt.id request format
- `RequestCredential` - Credential constraints
- `InputDescriptor` - Verification criteria
- `VerificationOptions` - User selections

## Deployment

### Docker

```bash
docker compose up -d
```

### Manual

```bash
go build -o testa-sacco
./testa-sacco
```

### DigitalOcean

Deploy alongside Testa Gava (issuer) on same droplet or separate server.

## How It Works

### Verification Flow

1. **User Configures Verification**
   - Selects credential type
   - Chooses verification policies
   - Submits form

2. **Generate Verification Request**
   - Server builds OpenID4VP request
   - Sends to Walt.id verifier
   - Receives verification link

3. **Share with Credential Holder**
   - Copy verification link
   - Send to holder via email/SMS/QR code
   - Holder opens in wallet app

4. **Holder Presents Credential**
   - Wallet displays verification request
   - User consents to share
   - Credential presented to verifier

5. **Verification Complete**
   - Policies checked automatically
   - Results returned to verifier
   - Decision made based on results

## API Request Format

The form generates a Walt.id verification request:

```json
{
  "vc_policies": [
    "signature",
    "expired",
    "not-before",
    "revoked-status-list"
  ],
  "request_credentials": [
    {
      "format": "jwt_vc",
      "input_descriptor": {
        "id": "e3d700aa-0988-4eb6-b9c9-e00f4b27f1d8",
        "constraints": {
          "fields": [
            {
              "path": ["$.vc.type"],
              "filter": {
                "contains": {
                  "const": "VerifiablePortableDocumentA1"
                },
                "type": "array"
              }
            }
          ]
        }
      }
    }
  ]
}
```

Headers sent:

- `Content-Type: application/json`
- `Accept: text/plain`
- `authorizeBaseUrl: openid4vp://authorize`
- `responseMode: direct_post`

## Customization

### Add New Credential Types

Edit `templates/index.html` around line 50:

```html
<option value="YourCustomType">Your Custom Credential Type</option>
```

### Change Styling

Edit `static/styles.css` - uses blue theme for verifier

### Add New Handlers

1. Create function in `handlers/handler.go`
2. Add route in `main.go`

## Testing

```bash
# Start server
go run main.go

# Configure verification at http://localhost:8081
# 1. Select "Verifiable Portable Document A1"
# 2. Keep all policies checked
# 3. Click "Generate Verification Request"
# 4. Copy the generated verification link
```

## Security

- Environment-based configuration
- Policy-based verification
- OpenID4VP protocol compliance
- HTTPS ready
- Non-root Docker container

## Tips for SACCO Operators

1. **Credential Types**: Use specific types to prevent incorrect credentials
2. **Link Sharing**: Send verification links via secure channels
3. **Result Handling**: Implement webhook handlers for verification results
4. **Logging**: Monitor verification requests for compliance

## Troubleshooting

### "Package not found"

```bash
go mod tidy
```

### "Templates not found"

Ensure `templates/` directory exists with `index.html`

### "Connection refused"

- Verify Walt.id verifier is running on port 7003
- Check `WALTID_VERIFIER_URL` configuration
- Ensure network connectivity

### "Form not submitting"

Check browser console - HTMX should be loaded from CDN

## 🔗 Integration with Testa Gava

Testa SACCO (verifier) works alongside Testa Gava (issuer):

```text
Testa Gava (Port 8080)           Testa SACCO (Port 8081)
Issue credentials         ←→      Verify credentials
↓                                 ↓
Walt.id Issuer (7002)            Walt.id Verifier (7003)
```

**Typical Flow**:

1. Farmer gets credential from Testa Gava
2. Farmer stores in wallet (e.g., Walt.id, Inji, etc.)
3. SACCO requests verification via Testa SACCO
4. Farmer presents credential from wallet
5. SACCO receives verified information

## Protocol Details

### OpenID4VP

- **Protocol**: OpenID for Verifiable Presentations
- **Standard**: W3C Verifiable Credentials
- **Flow**: Request → Present → Verify
- **Privacy**: Holder controls what to share

### Selective Disclosure

While not configured in the UI, Walt.id supports:

- Requesting specific fields only
- Holder choosing what to reveal
- Zero-knowledge proofs

---

Built with **Go** and **HTMX**, and for **Server Side Rendered** secure credential verification
