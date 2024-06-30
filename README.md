# GophKeeper

ЯндексПрактикум финальный проект. Когорта 25+


## Описание проекта

Сервис для хранения паролей и другой чувствительной информации пользователя.

В проект входят:
- `client` - хранит данные пользователей. Сохраняет и предоставляет по запросу от клиента.
- `server` - предоставляет пользователю возможность сохранить и посмотреть свои пароли и другие приватные данные.


Действия пользователя:

- Пользователь скачивает клиент для своей ОС.
- Регистрируется на сервере.
- Аутентифицируется на сервер и получает от сервера JWT ключ для работы. Время жизни ключа можно настроить с помощью ENV.
- После логина пользователь попадает в private-menu, где может просмотреть имеющиеся данные либо сохранить новые.
- При отправке данных на сервер (или получении) программа запросит ввод мастер-ключа, которым будут зашифрованы данные.



### Технологии

- `MySQL` чтобы хранить данные пользователей на сервере.
- `encoding/gob` для кодирования структур с данными.
- `github.com/grpc/grpc-go` для обмена сообщениями между сервером и клиентом.
- `github.com/awnumar/memguard` для хранения мастер-ключа пользователя в памяти.
- `github.com/charmbracelet/bubbletea` красивый TUI на клиенте. 


### Безопасность

Чтобы данные не попали к третьим лицам используются следующие подходы.

- На транспортном уровне используется шифрование с помощью `TSL`. 
    В папке x509 лежат сертификаты и ключи для сервера. Ключи были взяты из [примеров работы с gRPC](https://github.com/grpc/grpc-go/tree/master/examples/data/x509).

- Данные пользователя шифруются мастер-ключом с использованием алгоритма AES. 
    Данный подход защищает от прочтения данных на сервере.
    По идее мастер-ключ должен быть они для всех, но есть ~~баг~~ фича что можно задавать разные мастер-пароли.

- Мастер-ключ хранится в памяти в защищённом хранилище `memguard`. Это может защитить мастер-ключ от дампа памяти.

- В любое время можно нажать `CTRL+C` и закрыть программу. 
    Данная инновация сбережёт ваше ментальное здоровье и не даст злоумышленнику подглядеть пароль из-за спины.


### Что есть

- Регистрация, Аутентификация пользователя.
- Сохранение приватных данных: логин-пароль, тексе, двоичные, банковская карта.
- Синхронизация данных между несколькими авторизованными клиентами.
- Передача приватных данных владельцу по запросу.
- Просмотр сохранённых данных.
- Настройка сервера и клиента через ENV переменные.
- Поддержка терминального интерфейса (TUI — terminal user interface).

### Чего нет (появится в платной версии)

- Удаление данных.
- Изменение уже сохранённых данных.
- Иные типы данных.
- Красивый клиентский интерфейс как в стримах по bubbletea.
- Красивая обработка фатальных ошибок.
- Красиво организованный код.
- Покрытие тестами на 80%.
- Поддержка данных типа OTP (one time password);
- Функциональных и интеграционных тестов.
- Описание протокола взаимодействия клиента и сервера в формате Swagger.


## dev

Обновление proto.
```shell
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/messages.proto
```

### docker

Используется MySQL в качестве хранилища.
```shell
docker compose up -d
```

