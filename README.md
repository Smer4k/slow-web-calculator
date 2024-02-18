# slow-web-calculator
- Версия golang `1.22.0`
- Символы которые поддерживает оркестр `+` `-` `/` `*`
## Способы запуска
- go run `main.go`
- Запуск `.exe`

## Запуск оркестра:

### go run `main.go`
Перед запуском оркестра этим способом нужно установить `gcc`
#### Установка `gcc`:
1. Скачайте установщик по ссылке ["Нажми сюда"](https://github.com/msys2/msys2-installer/releases/download/2024-01-13/msys2-x86_64-20240113.exe), после установите и следуйте инструкции установщика.
2. После установки не снимайте галочку с `Run MSYS2 now`.
3. Когда откроется командная строка используйте эту команду: `pacman -S --needed base-devel mingw-w64-ucrt-x86_64-toolchain`
4. После ввода команды, нажмите `Enter` и подтвердите в консоли установку введя `Y`
5. Добавьте путь к папке `bin` MinGW-w64 в переменную среды Windows PATH, выполнив следующие действия:
- Откройте поиск и введите `Изменение системных переменных среды` или `Edit environment variables for your account`
- Зайдите в `Переменные среды...`
- В переменнах `пользователя` выберите `Path` и нажмите `изменить`
- Нажмите `Добавить` и введите путь до `bin` gcc, если вы не менял путь установки, то по умолчанию это `C:\msys64\ucrt64\bin`
- Нажмите `OK` чтобы сохранить изменения, после откройте заново все `cmd` и `VS code`

Примечание: если вы не поняли что делать или у вас не получилось, то перейдите по ссылки оригинального источника ["Нажми сюда"](https://code.visualstudio.com/docs/cpp/config-mingw#_installing-the-mingww64-toolchain)

#### Запуск через `.bat`:
1. Зайдите в папку куда вы сохранили файлы проекта
2. Запустите `RunOrchestrator.bat`

#### Запуск через `cmd` или `PowerShell`:
1. Перейдите на диск где вы сохранили файлы проекта введя название диска, пример:
```cmd
S:
```
2. Перейдите в папку `slow-web-calculator\cmd\orchestrator`
```cmd
cd {путь до папки}\slow-web-calculator\cmd\orchestrator
```
3. Включите переменную `CGO_ENABLED`
```cmd
set CGO_ENABLED=1
```
3. Запустите сервер
```cmd
go run main.go
```

### Запуск `.exe`
1. Зайдите в папку куда вы сохранили файлы проекта
2. Зайдите в папку `cmd` после в папку `orchestrator`
3. Запустите файл `orchestrator.exe`

Примечание: перемещять файл `orchestrator.exe` нельзя.

После запуска оркестра откройте браузер и перейдите по ссылки `http://localhost:8080/`
## Запуск агента:
Каждый запущенный агент это отдельный сервер.
### go run `main.go`

1. Откройте `cmd` или `PowerShell`
2. Перейдите на диск где вы сохранили файлы проекта введя название диска, пример:
```cmd
S:
```
3. Зайдите в папку `slow-web-calculator\cmd\agent`
```cmd
cd {путь до папки}\slow-web-calculator\cmd\agent
```
4. Запустить сервер
```cmd
go run main.go
```
Примечание: если нужно больше одного агента, то запустите сервер с доп.параметром:
```cmd
go run main.go -port={ваш порт}
```

### Запуск `.exe`
1. Зайдите в папку куда вы сохранили файлы проекта
2. Зайдите в папку `cmd` после в папку `agent`
3. Запустите файл `agent.exe`

Примечание: больше одного агента таким способом нельзя запустить
