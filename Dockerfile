# Используем официальный образ Go для сборки приложения
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем модули и зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Сборка исполняемого файла
RUN go build -o app .

# Используем минимальный образ для запуска собранного приложения
FROM debian

# Устанавливаем рабочую директорию в конечном образе
WORKDIR /app

# Копируем исполняемый файл из стадии сборки
COPY --from=builder /app/app .

# Открываем порт для взаимодействия
EXPOSE 8080

# Устанавливаем команду запуска
CMD ["./app"]