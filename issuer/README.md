# Testa Gava

Production-ready Go + HTMX web application for issuing W3C Verifiable Credentials with full A1 Portable Document support.

## Architecural Overview

- **Modular Structure**: Separated handlers, models, and main server
- **Dynamic Form**: Full web form for all credential fields
- **Package Organization**: Clean Model-View-Controller (MVC)-like pattern
- **Complete A1 Support**: All sections of Portable Document A1 with Selective Disclosure support

### Key Features

- **Dynamic Data Entry**: Web form captures all credential subject fields
- **Selective Disclosure**: Complete SD-JWT configuration for all sections
- **Production Ready**: Proper separation of concerns
- **Easy to Extend**: Add new handlers or models easily

### Project Structure

```text
issuer/
├── main.go                    # Server configuration
├── handlers/
│   └── handler.go            # All HTTP handlers
├── models/
│   └── credential.go         # Data structures
├── templates/
│   └── index.html            # Dynamic form
├── static/
│   └── styles.css            # Enhanced CSS with form styles
├── go.mod                     # Go module
├── Dockerfile                 # Container config
├── docker-compose.yml         # Easy deployment
├── Makefile                   # Dev shortcuts
└── README.md                  # This file
```

## Quick Start

```bash
# 1. Navigate to project
cd issuer

# 2. Run the server
go run main.go

# 3. Visit http://localhost:8082
```

### Form Sections

The web form captures all A1 Portable Document sections:

#### Section 1: Personal Information

- Personal ID, Sex, Name
- Date of Birth, Nationalities
- Residence and Stay Addresses

#### Section 2: Legislation Information

- Member State
- Start/End Dates
- Certificate Duration flags

#### Section 3: Activity Type

- Employment status checkboxes
- Exception handling

#### Section 4: Employer/Business Information

- Business name and address
- Employment type

#### Section 5: Work Location

- Fixed address configuration

#### Section 6: Issuing Institution

- Institution details
- Contact information
- Signature

### Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `WALTID_ISSUER_URL` | `http://droplet_ip:7002/openid4vc/sdjwt/issue` | Walt.id issuer endpoint |
| `PORT` | `8082` | Server port |

### Architecture

#### main.go

Simple server setup with routing

#### handlers/handler.go

All HTTP logic:

- `Home()` - Renders form
- `IssueCredential()` - Processes form and issues credential
- `extractFarmerData()` - Parses form data
- `buildCredentialRequest()` - Creates Walt.id request
- `renderSuccess()` / `renderError()` - HTMX responses

#### models/credential.go

Complete data structures:

- `FarmerCredential` - Form data
- `CredentialRequest` - Walt.id request format
- All Section models (Section1-6)
- Selective Disclosure configuration

### Deployment

#### Docker

```bash
docker-compose up -d
```

#### Manual

```bash
go build -o testa-gava
./testa-gava
```

### Customization

#### Add New Fields

1. Add to `models/credential.go` (FarmerCredential struct)
2. Add to `templates/index.html` (form fields)
3. Add to `handlers/handler.go` (extractFarmerData function)

#### Change Styling

Edit `static/styles.css` - form styles are clearly marked

#### Add New Handlers

1. Create function in `handlers/handler.go`
2. Add route in `main.go`

### Testing

```bash
# Start server
go run main.go

# Fill form at http://localhost:8082
# Click "Issue Digital ID Credential"
# Copy the generated credential link
```

#### API Request Format

The form generates a complete Walt.id request with:

- Issuer keys (JWK)
- Credential configuration ID
- Full credential data (all 6 sections)
- Mapping configuration
- Selective disclosure for all fields
- Issuer DID

### Security

- Environment-based configuration
- Form validation
- Error handling
- HTTPS ready

### Tips for DPI Developers

1. **Pre-fill Common Data**: Update default values in `templates/index.html`
2. **Validate Input**: Add validation in `extractFarmerData()`
3. **Custom Sections**: Extend Section models as needed
4. **Database Integration**: Add persistence layer in handlers

## Troubleshooting

### "Package not found"

```bash
go mod tidy
```

### "Templates not found"

Ensure `templates/` directory exists with `index.html`

### "Form not submitting"

Check browser console - HTMX should be loaded from CDN

---

Built with **Go** and **HTMX**, and for **Server Side Rendered** digital identity
