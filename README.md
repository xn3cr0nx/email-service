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

| Name                                | Type   | Default          | Range | Description                                          |
| ----------------------------------- | ------ | ---------------- | ----- | ---------------------------------------------------- |
| **DEBUG**                           | bool   | `false`          |       | Sets logging level to Debug                          |
| **PORT**                            | int    | `6066`           |       | Bind http server to port                             |
| **NAME**                            | str    | `Mailer`         |       | Set service name                                     |
| **REST**                            | bool   | `false`          |       | Enable exposed REST API to interact with the service |
| **TEMPLATE_DIR**                    | string | `/templates`     |       | Define templates folder path                         |
| **KAFKA**                           | bool   | `false`          |       | Set kafka as broker backend                          |
| **KAFKA_ADDRESS**                   | arr    | `localhost:9092` |       | Set kafka addresses                                  |
| **KAFKA_GROUP**                     | str    | `my-group`       |       | Set kafka group name                                 |
| **KAFKA_TOPIC**                     | str    | `emails`         |       | Set kafka topic name                                 |
| **REDIS_ADDRESS**                   | string | `localhost:6379` |       | Set host address for redis backend                   |
| **REDIS_PASSWORD**                  | string | ``               |       | Set password address for redis backend               |
| **REDIS_DB**                        | int    | `2`              |       | Set redis database number                            |
| **CONCURRENCY**                     | int    | `10`             |       | Set number of concurrent workers for redis backend   |
| **OTEL_EXPORTER_JAEGER_ENABLE**     | bool   | `false`          |       | Enable OpenTelemetry based jager tracing             |
| **OTEL_EXPORTER_JAEGER_AGENT_HOST** | str    | `jaeger`         |       | Override Jaeger agent hostname                       |
| **OTEL_EXPORTER_JAEGER_AGENT_PORT** | int    | `14268`          |       | Override Jaeger agent port                           |
| **OTEL_EXPORTER_PROMETHEUS_ENABLE** | bool   | `false`          |       | Enable OpenTelemetry based prometheus metrics        |
| **OTEL_EXPORTER_PROMETHEUS_PORT**   | int    | `9464`           |       | Override Prometheus exposed port                     |

## TODO

[ ] unit testing
