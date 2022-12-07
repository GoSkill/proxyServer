# proxySever
Практическая работа 31
написать микросервис 
написать proxy(который обслуживает две реплики данного приложения)
информацию сохранять в любой базе данных

Запуск:
1 реплика /cmd/web/go run main.go (Addr = http://localhost:8080)
2 реплика /cmd/web/go run main.go (изменить Addr = http://localhost:8081)
прокси /proxy/go run main.go (Addr = http://localhost:8800)

Работа:
В терминале отображается адрес по которому прокси отправляет запросы
Взаимодействие с сервисом по http://localhost:8800
Просмотр через "меню" - список пользователей, "кнопки ID" - ФИО пользователя и его друзей
Создание пользователя через "форму"
Другие функции через терминал. Команды "curl" на главной странице

Завершение работы:
В терминале "CTRL+C" (для Windows)

Техническая информация:
При создании сервиса использовался фреймворк GIN и база данных MYSQL

Блокнот базы данных "proxyserver"
Таблицы: 
"Persons" - хранилище пользователей (4 столбца)
CREATE TABLE Persons (PersonID int NOT NULL PRIMARY KEY AUTO_INCREMENT, FirstName varchar(20) NOT NULL, LastName varchar(20) NOT NULL, Age int NOT NULL);

"Friendship" - хранилище друзей (2 столбца)
CREATE TABLE Friendship  (SourceID int  NOT NULL, TargetID int NOT NULL, FOREIGN KEY(SourceID) REFERENCES Persons(PersonID), FOREIGN KEY (TargetID) REFERENCES Users (UserID) ON UPDATE CASCADE ON DELETE CASCADE);

"Users" - (4 столбца) дублирует "Persons" и необходима для получения "представления" (виртуальной таблицы) при получении списка друзей пользователя
CREATE TABLE Users (UserID int NOT NULL PRIMARY KEY AUTO_INCREMENT, FirstN varchar(20) NOT NULL, LastN varchar(20) NOT NULL, AgeU int NOT NULL);

"Persons_Friendship_Summary" - представляет выборку всех друзей пользователя по его ID
CREATE VIEW Persons_Friendship_Summary AS SELECT PersonID AS pfs_ID, max(FirstName) AS pfs_FirstName, cast(concat('[', group_concat(json_quote(LastN) ORDER BY LastN SEPARATOR ','), ']') as json) AS pfs_Friend_array FROM Persons INNER JOIN Friendship ON Persons.PersonID = Friendship.SourceID INNER JOIN Users ON Friendship.TargetID = Users.UserID GROUP BY Persons.PersonID
