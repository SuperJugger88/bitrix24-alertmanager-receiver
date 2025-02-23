# Receiver Bitrix24

## Описание
Receiver Bitrix24 — это сервис, предназначенный для облегчения интеграции Bitrix24 с Alertmanager. Он принимает входящие вебхуки и обрабатывает их соответствующим образом, обеспечивая бесшовную связь между системами.

## Как это работает
Receiver работает, получая уведомления вебхуков от Alertmanager. Он обрабатывает эти уведомления и пересылает их в соответствующие каналы Bitrix24, обеспечивая эффективную коммуникацию предупреждений.

## Docker-образ
Receiver Bitrix24 упакован в виде Docker-образа, собрать можно следующей командой:

```bash
docker buildx build --build-arg=BITRIX_WEBHOOK_URL=<your_webhook_url> --build-arg=APP_PORT=4000 --build-arg=MESSAGE_TEMPLATE_PATH=/etc/bitrix24.message.tmpl  --platform linux/amd64 --push  -t <your_docker_registry>/<helpers_path>/bitrix24/receiver:latest .
```

## Развертывание
Сервис развернуть на сервере `<your_metrics_server>` в следующем каталоге:

```
/opt/<container_manifests_path>/bitrix24_receiver/compose.yml
```

Со следующей конфигурацией

```yaml
version: '3.9'

services:
  bitrix24-receiver:
    image: '<your_docker_registry>/<helpers_path>/bitrix24/receiver:latest'
    restart: always
    volumes:
      - /etc/alertmanager/templates/bitrix24.message.tmpl:/etc/bitrix24.message.tmpl
    ports:
      - 127.0.0.1:4000:4000
```

## Интеграция с Alertmanager
Receiver Bitrix24 работает вместе с Alertmanager через конфигурацию вебхука. Он использует конечную точку вебхука `bitrix24.message.tmp` для получения предупреждений и уведомлений от Alertmanager, что позволяет управлять предупреждениями и коммуникацией в реальном времени.

```yaml
# Your configuraion
- name: 'portal'
  webhook_configs:
    - url: 'http://localhost:4000/webhook?dialog_id=chatXXXXX'
      send_resolved: true
```

---

Для получения дополнительной информации или поддержки, пожалуйста, обратитесь к документации проекта или свяжитесь с командой разработчиков.
