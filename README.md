# Сервис подсчёта арифметических выражений
Финальная задача Спринта 1 на курсе Golang в Яндекс Лицее

## Описание
Программа основывается на преобразовании арифметического выражения из инфиксной в постфиксную запись при помощи алгоритма сортировочной станции и вычислении его.

## Особенности реализации
В данном задании реализуется работа с HTTP запросами и ответами.

# Тестирование
Для запуска кода нужно сохранить архив, распаковать его, запустить терминал в папке и запустить сервер с помощью команды `go run main.go`

Далее в отдельном терминале нужно писать запросы

## Успешное решение (200 OK)
### Ввод
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(6+9)*3\"}" http://127.0.0.1:8080/api/v1/calculate
### Вывод
{"result":"45.000000"}

## Неверное выражение (422 Unprocessable Entity)
### Ввод
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(6+9)*a\"}" http://127.0.0.1:8080/api/v1/calculate
### Вывод
{"error":"Expression is not valid"}

### Ввод
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"\"}" http://127.0.0.1:8080/api/v1/calculate
### Вывод
{"error":"Expression is not valid"}

## Иная ошибка (500 Internal Server Error)
### Ввод
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(6+9)/0\"}" http://127.0.0.1:8080/api/v1/calculate
### Вывод
{"error":"Internal server error"}
