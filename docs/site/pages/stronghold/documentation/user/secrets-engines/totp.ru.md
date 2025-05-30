---
title: Механизм секретов TOTP
sidebar_label: Механизм секретов TOTP
permalink: ru/stronghold/documentation/user/secrets-engines/totp.html
description: Механизм секретов TOTP для Stronghold генерирует одноразовые пароли на основе текущего времени.
lang: ru
---

Механизм секретов TOTP генерирует временные учетные данные в соответствии со стандартом TOTP.
Этот механизм также может использоваться для создания нового ключа и проверки паролей,
сгенерированных с помощью этого ключа.

Механизм секретов TOTP может работать как в режиме генератора (аналогично Google Authenticator),
так и в режиме провайдера (аналогично сервису входа Google.com).

## Режим генератора

Механизм секретов TOTP может выступать в роли генератора TOTP-кодов.
В этом режиме он может заменить традиционные генераторы TOTP, такие как Google Authenticator.
Это обеспечивает дополнительный уровень безопасности, поскольку возможность генерации кодов
защищена политиками, а весь процесс аудируется.

### Настройка

Большинство механизмов секретов требуют предварительной настройки перед использованием.
Эти шаги обычно выполняются оператором или инструментом управления конфигурацией.

Активируйте механизм секретов TOTP:

```shell-session
$ stronghold secrets enable totp
Success! Enabled the totp secrets engine at: totp/
```

По умолчанию механизм монтируется по имени. Чтобы активировать его по другому пути, используйте аргумент `-path`.

Настройте именованный ключ. Имя этого ключа будет использоваться для его идентификации.

```shell-session
$ stronghold write totp/keys/my-key \
    url="otpauth://totp/Stronghold:test@test.com?secret=Y64VEVMBTSXCYIWRSHRNDZW62MPGVU2G&issuer=Stronghold"
Success! Data written to: totp/keys/my-key
```

Параметр `url` соответствует секретному ключу или значению из штрих-кода, предоставленного сторонним сервисом.

### Использование

После настройки механизма и получения пользователем или системой токена Stronghold с соответствующими
разрешениями можно генерировать учетные данные.

Сгенерируйте новый OTP-код, прочитав данные из `/code` с указанием имени ключа:

```shell-session
$ stronghold read totp/code/my-key
Key     Value
---     -----
code    260610
```

Используя ACL, можно ограничить доступ к механизму TOTP так, чтобы доверенные операторы
могли управлять определениями ключей, а пользователи и приложения были ограничены в доступе к
определенным учетным данным.

## Режим провайдера

Механизм секретов TOTP также может выступать в роли провайдера TOTP.
В этом режиме он может использоваться для генерации новых ключей и проверки паролей, созданных с помощью этих ключей.

### Настройка

Активируйте механизм секретов TOTP:

```shell-session
$ stronghold secrets enable totp
Success! Enabled the totp secrets engine at: totp/
```

По умолчанию механизм монтируется по имени. Чтобы активировать его по другому пути, используйте аргумент `-path`.

Создайте именованный ключ, используя параметр `generate`. Это указывает Stronghold работать в режиме провайдера:

```shell-session
$ stronghold write totp/keys/my-user \
    generate=true \
    issuer=Stronghold \
    account_name=user@test.com

Key        Value
---        -----
barcode    iVBORw0KGgoAAAANSUhEUgAAAMgAAADIEAAAAADYoy0BA...
url        otpauth://totp/Stronghold:user@test.com?algorithm=SHA1&digits=6&issuer=Stronghold&period=30&secret=V7MBSK324I7KF6KVW34NDFH2GYHIF6JY
```

Ответ включает base64-кодированный штрих-код и OTP-ссылку. Оба варианта эквивалентны.
Передайте их пользователю, который проходит аутентификацию с помощью TOTP.

### Использование

Пользователь может проверить TOTP-код, сгенерированный сторонним приложением:

```shell-session
$ stronghold write totp/code/my-user code=886531
Key      Value
---      -----
valid    true
```
