# proxyServer
Прокси сервер на две реплики микросервиса на базе MYSQL.

Практическая работа 31: 
написать микросервис 
написать proxy(который обслуживает две реплики данного приложения)
информацию сохранять в любой базе данных

Запуск
по умолчанию установлено
host1="http://:8080"
прокси  addr="http://:8800" 

доступ к MYSQL
-dsn="web:1@/proxyserver?parseTime=true"

Пример команды запуска 1-й реплики
go run ./cmd/web        (слушать 127.0.0.1:8080)

Изменяем порт, например
go run ./cmd/web -addr="127.0.0.1:8085"

Изменяем порт и доступ к MYSQL
go run ./cmd/web -addr="127.0.0.1:8085" -dsn="web:1@/proxyserver?parseTime=true"

Пример команды запуска прокси
go run ./proxy -host1="http://:8080" -host2="http://:8085"

Работа:
В терминале отображается адрес по которому прокси отправляет запросы
Взаимодействие с сервисом по 127.0.0.1/:8800

Через "меню" можно смотреть
 - список пользователей
 - друзей пользователя
 - cоздать нового пользователя
 - пригласить друга
 - удалить пользователя
 - изменить возраст пользователя
 
Завершение работы
В терминале "CTRL+C" (для Windows)

Техническая информация
При создании сервиса использовался фреймворк GIN и база данных MYSQL

Блокнот базы данных "proxyserver" из 3-х таблиц.
Таблицы: 
"Persons" - хранилище пользователей (4 столбца)
CREATE TABLE Persons (PersonID int NOT NULL PRIMARY KEY AUTO_INCREMENT, FirstName varchar(20) NOT NULL, LastName varchar(20) NOT NULL, Age int NOT NULL);

"Friendship" - хранилище друзей (2 столбца)
CREATE TABLE Friendship  (SourceID int  NOT NULL, TargetID int NOT NULL, FOREIGN KEY(SourceID) REFERENCES Persons(PersonID), FOREIGN KEY (TargetID) REFERENCES Users (UserID) ON UPDATE CASCADE ON DELETE CASCADE);

"Users" - (4 столбца) дублирует "Persons" и необходима для получения "представления" (виртуальной таблицы) при получении списка друзей пользователя
CREATE TABLE Users (UserID int NOT NULL PRIMARY KEY AUTO_INCREMENT, FirstN varchar(20) NOT NULL, LastN varchar(20) NOT NULL, AgeU int NOT NULL);

"Persons_Friendship_Summary" - представляет выборку всех друзей пользователя по его ID
CREATE VIEW Persons_Friendship_Summary AS SELECT PersonID AS pfs_ID, max(FirstName) AS pfs_FirstName, group_concat(LastN ORDER BY LastN SEPARATOR ',') AS pfs_Friend_array FROM Persons INNER JOIN Friendship ON Persons.PersonID = Friendship.SourceID INNER JOIN Users ON Friendship.TargetID = Users.UserID GROUP BY Persons.PersonID;

