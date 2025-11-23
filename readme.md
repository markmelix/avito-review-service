## Avito PR Review Service

Единый микросервис, который автоматически назначает ревьюеров на Pull Request’ы (PR), а также позволяет управлять командами и участниками.

[Подробное описание задания](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2025/Backend-trainee-assignment-autumn-2025.md)

### Запуск

Запуск сервисов проекта осуществляется командой:

```bash
docker compose up
```

Взаимодействие с сервисом происходит по адресу http://localhost:8080

### Нагрузочное тестирование

#### Запуск

Для запуска нагрузочного тестирования остановите production-сервисы и запустите контейнеры соответствующего окружения:

```bash
docker compose down
docker compose -f compose.loadtest.yml down -v
docker compose -f compose.loadtest.yml up
```

После получения репорта в stdout остановите Docker Compose. Результаты тестирования будут лежать в директории `./result`. Далее их можно сконвертировать в график или CLI-представление командами:

```bash
cat result/vegeta_results.bin | vegeta report
vegeta plot result/vegeta_results.bin > results.html
xdg-open results.html
```
