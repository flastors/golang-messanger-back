Инструкции по Swagger (swag)

1) Установить инструмент swag (если ещё не установлен):

```bash
# go 1.20+ recommended
go install github.com/swaggo/swag/cmd/swag@latest
```

2) Перейти в корень проекта и сгенерировать документацию (будет создана папка `docs`):

```bash
swag init -g internal/controllers/restapi/docs.go -o docs
```

3) Запустить приложение (локально или в Docker). Swagger UI будет доступен по:

http://localhost:8080/swagger/index.html

Примечания:
- Я добавил минимальный placeholder `docs/docs.go` чтобы приложение собиралось без `swag init`. После запуска `swag init` вы получите полный json/yaml в `docs`.
- Если вы используете Docker, пробросьте порт 8080 и убедитесь, что приложение слушает на 0.0.0.0 (переменная окружения `CHAT_SERVICE_HOST` уже добавлена в compose).

