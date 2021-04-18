# Email Service

This service is a configurable and reusable email service based on Postmark.

The goal of this service is to provide an easy plug and play service to send emails for your business. The service is by default configured to rely on a broker as Kafka, in order to treat sending emails as events. The service also allow to start a REST API to manually interact with the service.

## Support

- Kafka
- Swagger documentation
- OpenTelemetry

## Environment variables

| Name      | Type | Default | Range | Description         |
| --------- | ---- | ------- | ----- | ------------------- |
| **DEBUG** | bool | `false` |       | Set debug log level |
