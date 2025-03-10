# UOFDBot - User of the day [bot]

<p align="center">
  <img width="100px" src="https://user-images.githubusercontent.com/2866780/72239870-6a010c00-35f3-11ea-9d8f-9d499762e1bb.png"></img>
</p>

Веселый телеграм бот. Поможет определиться, кто сегодня 'пидор', а кто 'герой'. Каждая игра запускается отдельно. При
необходимости статистику можно сбросить.

### Запуск

#### Docker Compose

```shell
docker compose up -d --build
```

##### Переменные окружения

| Переменная            | Описание                                  | Пример                  | По умолчанию                                                                                                |
|-----------------------|-------------------------------------------|-------------------------|-------------------------------------------------------------------------------------------------------------|
| UOFD_DB_FILE_PATH     | Путь до файла БД SQLite                   | /opt/UOFDBot/uofd.db    | `/opt/UOFDBot/default/uofd.db`                                                                              |
| UOFD_CONFIG_FILE_PATH | Путь до файла конфигурации                | /opt/UOFDBot/config.yml | Берётся файл из репозитория `configs/config.yml` и копируется в контейнер `/opt/UOFDBot/default/config.yml` |
| UOFD_LANG_DIR_PATH    | Путь до директории с ресурсами локализаци | /opt/UOFDBot/lang       | Файлы из репозитория по пути `lang` копируется в контейнер `/opt/UOFDBot/default/lang`                      |
