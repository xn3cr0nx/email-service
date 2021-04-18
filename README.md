# Email Service

This service is a configurable and reusable email service based on Postmark.

The goal of this service is to provide an easy plug and play service to send emails for your business.
The service can be configured in order to rely on a broker as Kafka, or async queue, backend by redis as fallback, in order to treat sending emails as events. The service also allow to start a REST API to instantly ask for email sending.

## Support

- Kafka
- Asynq
- Redis
- Swagger documentation
- OpenTelemetry

## Environment variables

| Name      | Type | Default | Range | Description         |
| --------- | ---- | ------- | ----- | ------------------- |
| **DEBUG** | bool | `false` |       | Set debug log level |

## TODO

[] Prometheus emails counter
