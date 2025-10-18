# Tenant Repository: {{.TenantName}}

This repository contains the data and configuration for tenant: {{.TenantName}}

## Structure

- `data/` - Tenant data files
- `config/` - Configuration files
- `assets/` - Static assets

## Usage

This repository is automatically managed by the transport-app system.

## Configuration

The tenant configuration is stored in `config/tenant.json` and follows this structure:

```json
{
  "tenant": {
    "name": "{{.TenantName}}",
    "description": "",
    "created_at": "{{.CreatedAt}}",
    "status": "active"
  },
  "site": {
    "title": "{{.TenantName}} - Transport App",
    "description": "Transportation services for {{.TenantName}}",
    "theme": "default"
  }
}
```

## Data Files

Place your tenant-specific data files in the `data/` directory. Supported formats:
- JSON files for structured data
- CSV files for tabular data
- Markdown files for documentation
- Images and other assets in `assets/`

## ONA Integration

This repository is designed to work with ONA (Open Network Architecture) for automatic website generation. The structure follows ONA's expected format for seamless integration.
