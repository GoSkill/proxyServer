# proxyServer
Прокси сервер на две реплики микросервиса на базе MYSQL.
Proxy, который обслуживает две реплики данного приложения, информацию сохраняет в MYSQL

Запуск
по умолчанию установлено
host1="http://:8080"
прокси  addr="http://:8800" 

доступ к MYSQL
-dsn="web:1@/proxyserver?parseTime=true"

Пример команды запуска 1-й реплики

go run ./cmd/web        (слушать 127.0.0.1:8080)

Запуск 2 реплики

Изменяем порт, например
go run ./cmd/web -addr="127.0.0.1:8085"

Изменяем порт и доступ к MYSQL
go run ./cmd/web -addr="127.0.0.1:8085" -dsn="web:1@/proxyserver?parseTime=true"

Пример команды запуска прокси
go run ./proxy -host1="http://:8080" -host2="http://:8085"

Работа:

В терминале отображается адрес, по которому прокси отправляет запросы
Взаимодействие с сервисом по 127.0.0.1/:8800

Через "меню" можно смотреть
 - список пользователей
 - друзей пользователя
 - cоздать нового пользователя
 - пригласить друга
 - удалить пользователя
 - изменить возраст пользователя

![Экран 1](https://github.com/CHvvmu/proxySever/assets/96997574/d7d866b4-fd91-4521-aabd-8b19d0893022)

![Экран 2](https://github.com/CHvvmu/proxySever/assets/96997574/609331fc-9df8-4b13-a786-57d89b178347)

![Экран 3](https://github.com/CHvvmu/proxySever/assets/96997574/75ad83be-b3d0-42dc-b470-320642f683f8)

![Экран 4](https://github.com/CHvvmu/proxySever/assets/96997574/29f93b1b-1852-4627-a837-046366156a97)

Завершение работы
В терминале "CTRL+C" (для Windows)

Техническая информация
При создании сервиса использовался фреймворк GIN и база данных MYSQL

Блокнот базы данных MYSQL "proxyserver" из 3-х таблиц.
Таблицы: 
"Persons" - хранилище пользователей (4 столбца)

Создать 1 таблицу:

CREATE TABLE Persons 

PersonID int NOT NULL PRIMARY KEY AUTO_INCREMENT,

FirstName varchar(20) NOT NULL, 

LastName varchar(20) NOT NULL, 

Age int NOT NULL);

Создать 2 таблицу:

"Friendship" - хранилище друзей (2 столбца)
CREATE TABLE Friendship  
SourceID int  NOT NULL, 
TargetID int NOT NULL, 
FOREIGN KEY(SourceID) REFERENCES Persons(PersonID), 
FOREIGN KEY (TargetID) REFERENCES Users (UserID) ON UPDATE CASCADE ON DELETE CASCADE);

Создать 3 таблицу:

"Users" - (4 столбца) дублирует "Persons" и необходима для получения "представления" (виртуальной таблицы) при получении списка друзей пользователя
CREATE TABLE Users 
UserID int NOT NULL PRIMARY KEY AUTO_INCREMENT, 
FirstN varchar(20) NOT NULL, 
LastN varchar(20) NOT NULL, 
AgeU int NOT NULL);

Создать 4 таблицу:

"Persons_Friendship_Summary" - представляет выборку всех друзей пользователя по его ID
CREATE VIEW Persons_Friendship_Summary AS SELECT PersonID AS pfs_ID, 
max(FirstName) AS pfs_FirstName, 
group_concat(LastN ORDER BY LastN SEPARATOR ',') 
AS pfs_Friend_array FROM Persons INNER JOIN Friendship ON Persons.PersonID = Friendship.SourceID INNER JOIN Users ON Friendship.TargetID = Users.UserID GROUP BY Persons.PersonID;

