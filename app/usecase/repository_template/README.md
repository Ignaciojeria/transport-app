# Repository Templates

This directory contains embedded templates for creating tenant repositories.

## Files

- `tenant_repo_template.md` - README template for tenant repositories
- `tenant_config_template.json` - Configuration template for tenant data
- `gitignore_template` - Git ignore template for tenant repositories

## Usage

These templates are used by the `CreateRepositoryWorkflow` to generate initial files for new tenant repositories.

## Template Variables

The templates support the following variables:

- `{{.TenantName}}` - The name of the tenant (e.g., "tenant-123e4567-e89b-12d3-a456-426614174000")
- `{{.TenantID}}` - The tenant ID extracted from the name (e.g., "123e4567-e89b-12d3-a456-426614174000")
- `{{.CreatedAt}}` - The creation timestamp in RFC3339 format

## Structure Generated

When a tenant repository is created, it will have this structure:

```
tenant-{uuid}/
├── README.md              # Generated from tenant_repo_template.md
├── .gitignore            # Generated from gitignore_template
├── config/
│   └── tenant.json       # Generated from tenant_config_template.json
├── data/                 # Empty directory for tenant data
└── assets/               # Empty directory for static assets
```

## ONA Integration

These templates are designed to work with ONA (Open Network Architecture) for automatic website generation. The structure follows ONA's expected format for seamless integration.
