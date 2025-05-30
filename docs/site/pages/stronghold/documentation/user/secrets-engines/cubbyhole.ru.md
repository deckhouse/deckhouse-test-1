---
title: "Механизм секретов Cubbyhole"
permalink: ru/stronghold/documentation/user/secrets-engines/cubbyhole.html
lang: ru
---

Механизм секретов `cubbyhole` используется для хранения произвольных секретов в хранилище, привязанном к токену. В `cubbyhole`пути распределены по токенам. Ни один токен не может получить доступ к cubbyhole другого токена. Когда срок действия токена истекает, его `cubbyhole` уничтожается.

Также в отличие от движка секретов `kv`, поскольку время жизни `cubbyhole` связано со временем жизни токена аутентификации, не существует понятия TTL или интервала обновления для значений, содержащихся в `cubbyhole` токене.

Запись в ключ в механизме секретов `cubbyhole` полностью заменит старое значение.

## Настройка

Большинство секретных механизмов необходимо предварительно настроить, прежде чем они смогут выполнять свои функции. Эти шаги обычно выполняются оператором или инструментом управления конфигурацией.

Механизм секретов `cubbyhole` включен по умолчанию. Его нельзя отключить, переместить или включить несколько раз.

## Использование

После того как механизм секретов настроен и у пользователя есть токен с соответствующими правами, который может генерировать учетные данные. Механизм секретов `cubbyhole` позволяет записывать ключи с произвольными значениями

1. Запись произвольных данных:

   ```console
      $ d8 stronghold write cubbyhole/my-secret my-value=s3cr3t
   
      Success! Data written to: cubbyhole/my-secret
   ```

2. Чтение произвольных данных:

   ```console
      $ d8 stronghold read cubbyhole/my-secret
   
      Key         Value
   
      ---         -----
   
      my-value    s3cr3t
   ```
