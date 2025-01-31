# AddressInfo Structure

This document provides an overview of the `AddressInfo` structure and its fields, explaining their purpose and usage.

## Field Context

| Field          | Purpose                                                                                         |
|-----------------|------------------------------------------------------------------------------------------------|
| `Contact`      | Information about the contact associated with the address.                                      |
| `State`        | Represents the state, generally used in the U.S. and similar countries.                         |
| `Province`     | Represents the province, as in Canada or countries with specific subdivisions.                  |
| `County`       | A smaller administrative level than the state or province (e.g., counties).                     |
| `District`     | For urban subdivisions like neighborhoods or sectors.                                           |
| `AddressLine1` | Main address line (street and number).                                                          |
| `AddressLine2` | Additional information, such as an apartment or unit.                                           |
| `AddressLine3` | Optional additional details, such as references.                                                |
| `Latitude`     | Geographic coordinate for geolocation.                                                          |
| `Longitude`    | Geographic coordinate for geolocation.                                                          |
| `ZipCode`      | Postal code.                                                                                    |
| `TimeZone`     | Time zone associated with the location. 


Enviar logs a openobserve (el paso de permisos a topico) :

gcloud iam service-accounts create cloud-logs --display-name="Cloud Logs Service Account" --project=einar-404623


gcloud pubsub topics add-iam-policy-binding transport-app-events --member=serviceAccount:cloud-logs@einar-404623.iam.gserviceaccount.com --role=roles/pubsub.publisher --project=einar-404623

gcloud pubsub topics add-iam-policy-binding transport-app-events --member=serviceAccount:cloud-logs@einar-404623.iam.gserviceaccount.com --role=roles/pubsub.publisher  --project=einar-404623

KO_DOCKER_REPO=gcr.io/[TU_PROYECTO_ID]