Инструкция для запуска

// Пункты 1 и 2 необходимо проводить в powershell от имени администратора, остальные - нет
1. Ввести в Powershell: Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
2. Ввести в Powershell: Y
//

3. Ввести в Powershell: irm get.scoop.sh | iex
4. Ввести в Powershell: scoop install make
5. Ввести в Powershell: scoop install migrate
6. Ввести в Powershell: wsl --update
7. Запустить Docker
8. Открыть проект в Goland

//Опциональный пункт, если при открытии файла main.go в Goland написано GOROOT is not defined, выполните следующий пункт//
9. Установить GOROOT, если там написано <No SDK>, необходимо нажать на +, затем на download, выбрать go1.21.0, нажать OK
//

10. Ввести в терминал Goland: make postgres
11. Ввести в терминал Goland: make createdb
12. Ввести в терминал Goland: make migrateup
13. Запустить приложение (Shift + F10)
14. Открыть в браузере localhost:8080

Если хотите проверить работу миграции, введите в терминале Goland: "make migratedown" а затем "make migrateup"