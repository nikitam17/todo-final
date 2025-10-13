Итоговое задание практикума "Go-разработчик с нуля".
Разработан веб-сервер, который реализует функциональность планировщика задач — это будет аналог TODO-листа

# Файлы для итогового задания
В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.
Директория `web` содержит файлы фронтенда.
Директория `api` содержит обработчики и функции реализующие требования к проекту.
Директория `db` содержит функции работы с БД.

Список выполненных заданий со звёздочкой. Если их нет, напишите, что задания повышенной трудности не выполнялись.
1. Реализована возможность определять извне порт при запуске сервера. Если существует переменная окружения TODO_PORT, сервер при старте должен слушать порт со значением этой переменной.
2. Реализована возможность определять путь к файлу базы данных через переменную окружения. Для этого сервер должен получать значение переменной окружения TODO_DBFILE и использовать его в качестве пути к базе данных, если это не пустая строка.
3. Реализована возможность обработки правил для недель и месяцев.
4. Реализована возможность выбрать задачи через поле для поиска.
Инструкция по запуску кода локально: дополнительные флаги, примеры .env и так далее. Напишите, какой адрес следует указывать в браузере.
TODO_DBFILE=scheduler.db
TODO_PORT=7540
TODO_PASSWORD=123456
Инструкция по запуску тестов.
1. go test -run ^TestApp$ ./tests
2. go test -run ^TestDB$ ./tests
3. go test -run ^TestNextDate$ ./tests
4. go test -run ^TestAddTask$ ./tests
5. go test -run ^TestTasks$ ./tests
6. go test -run ^TestEditTask$ ./tests
7.1. go test -run ^TestDelTask$ ./tests
7.2. go test -run ^TestDone$ ./tests
В tests/settings.go следует использовать.
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGVja3N1bSI6IjhkOTY5ZWVmNmVjYWQzYzI5YTNhNjI5MjgwZTY4NmNmMGMzZjVkNWE4NmFmZjNjYTEyMDIwYzkyM2FkYzZjOTIifQ.3NnvnuzGQiP6jAKEXDtc_pJstvrWCbd0K8pBitF9R-c"`
Инструкция по сборке и запуску проекта через докер:
1. сборка: docker build --tag todo-final-app:v1 .
2. запуск: docker run -d -p 7540:7540 todo-final-app:v1
3. просмотр работающих образов: docker ps
4. останов образа: docker stop <id>
5. удаление образа: docker rmi -f todo-final-app:v1