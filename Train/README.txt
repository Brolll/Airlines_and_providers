Инструкция для запуска на Windows

// Пункты 1 и 2 необходимо проводить в powershell от имени администратора, остальные - нет
1. Ввести в Powershell: Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
2. Ввести в Powershell: Y
//

3. Ввести в Powershell: irm get.scoop.sh | iex
4. Ввести в Powershell: scoop install make
5. Ввести в Powershell: scoop install migrate
6. Запустить Docker

//Опциональный пункт, если при открытии файла main.go в Goland написано GOROOT is not defined, выполните следующий пункт//
7. Установить GOROOT, если там написано <No SDK>, необходимо нажать на +, затем на download, выбрать go1.21.0, нажать OK
//

8. Открыть Powershell в папке с проектом
9. Ввести в Powershell: make postgres
10. Ввести в Powershell: make createdb
11. Ввести в Powershell: make migrateup
12. Ввести в Powershell: .\train.exe
13. Открыть в браузере localhost:8080


Инструкция для запуска на Linux
1. Запустить Docker
2. Открыть терминал в папке с проектом
3. Ввести в терминал: sudo apt install make
4. Ввести в терминал: sudo apt install migrate
5. Ввести в терминал: make postgres
6. Ввести в терминал: make createdb
7. Ввести в терминал: make migrateup
8. Ввести в терминал: ./train.exe
9. Открыть в браузере localhost:8080

Если хотите проверить работу миграции, введите в терминале Goland: "make migratedown" а затем "make migrateup"
