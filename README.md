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
| **debug**                           | bool   | `false`          |       | Sets logging level to Debug                          |
| **port**                            | int    | `6066`           |       | Bind http server to port                             |
| **name**                            | str    | `Mailer`         |       | Set service name                                     |
| **rest**                            | bool   | `false`          |       | Enable exposed REST API to interact with the service |
| **template_dir**                    | string | `/templates`     |       | Define templates folder path                         |
| **kafka**                           | bool   | `false`          |       | Set kafka as broker backend                          |
| **kafka_address**                   | arr    | `localhost:9092` |       | Set kafka addresses                                  |
| **kafka_group**                     | str    | `my-group`       |       | Set kafka group name                                 |
| **kafka_topic**                     | str    | `emails`         |       | Set kafka topic name                                 |
| **redis_address**                   | string | `localhost:6379` |       | Set host address for redis backend                   |
| **redis_password**                  | string | ``               |       | Set password address for redis backend               |
| **redis_db**                        | int    | `2`              |       | Set redis database number                            |
| **concurrency**                     | int    | `10`             |       | Set number of concurrent workers for redis backend   |
| **otel_exporter_jaeger_enable**     | bool   | `false`          |       | Enable OpenTelemetry based jager tracing             |
| **otel_exporter_jaeger_agent_host** | str    | `jaeger`         |       | Override Jaeger agent hostname                       |
| **otel_exporter_jaeger_agent_port** | int    | `14268`          |       | Override Jaeger agent port                           |
| **otel_exporter_prometheus_enable** | bool   | `false`          |       | Enable OpenTelemetry based prometheus metrics        |
| **otel_exporter_prometheus_port**   | int    | `9464`           |       | Override Prometheus exposed port                     |

## TODO

[] unit testing
