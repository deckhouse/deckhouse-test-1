---
title: "Механизм секретов PKI"
permalink: ru/stronghold/documentation/user/secrets-engines/pki.html
lang: ru
---

Когда этот модуль секретов включен, экземпляр Stronghold может работать в качестве сервера сертификационного центра (CA), генерируя динамические корневые и промежуточные сертификаты, а также подписывая сертификаты по запросу.

Модуль секретов инфраструктуры открытых ключей (Public Key Infrastructure, сокращённо PKI) генерирует динамические сертификаты X.509. С этим механизмом секретов сервисы могут получать сертификаты без прохождения обычного ручного процесса создания закрытого ключа и CSR, подачи заявки в CA и ожидания завершения процесса верификации и подписания. Встроенные механизмы аутентификации и авторизации Stronghold обеспечивают функцию верификации.

Если поддерживать TTL относительно коротким, необходимость отзыва сертификатов становится менее вероятной, что позволяет держать CRL (certificate revocation list / лист отзыва сертификатов) коротким и помогает механизму секретов выдерживать большие рабочие нагрузки. Это, в свою очередь, позволяет каждому экземпляру запущенного приложения иметь уникальный сертификат, устраняя необходимость их совместного использования и связанные с этим проблемы отзыва и ротации.

Кроме того, благодаря возможности обойтись без отзыва, этот механизм секретов позволяет использовать эфемерные сертификаты. Сертификаты могут извлекаться и храниться в памяти при запуске приложения и удаляться при завершении работы, никогда не записываясь на диск.

## Настройка

Большинство механизмов секретов должно быть настроено заранее, чтобы они могли выполнять свои функции. Эти шаги обычно выполняются оператором или инструментом управления конфигурацией.

1. Включите механизм секретов PKI:

   ```shell
   $ d8 stronghold secrets enable pki
   Success! Enabled the pki secrets engine at: pki/
   ```

   По умолчанию механизм секретов будет установлен с именем движка. Чтобы включить механизм секретов по другому пути, используйте аргумент `-path`.

1. Увеличьте TTL, настроив механизм секретов. Значение по умолчанию в 30 дней может быть слишком коротким, поэтому увеличьте его до 1 года:

   ```shell
   $ d8 stronghold secrets tune -max-lease-ttl=8760h pki
   Success! Tuned the secrets engine at: pki/
   ```

   Обратите внимание, что отдельные роли могут ограничивать это значение до более короткого на основе каждого сертификата. Это лишь настраивает глобальное максимальное значение для этого механизма секретов.

1. Настройте сертификат CA и приватный ключ. Stronghold может использовать уже существующую пару ключей или сгенерировать собственный самоподписанный корневой сертификат. В общем случае, мы рекомендуем поддерживать ваш корневой CA вне Stronghold и предоставлять Stronghold подписанный промежуточный CA.

   ```shell
   $ d8 stronghold write pki/root/generate/internal \
        common_name=my-website.ru \
        ttl=8760h
   
   Key              Value
   ---              -----
   certificate      -----BEGIN CERTIFICATE-----...
   expiration       1756317679
   issuing_ca       -----BEGIN CERTIFICATE-----...
   serial_number    fc:f1:fb:2c:6d:4d:99:1e:82:1b:08:0a:81:ed:61:3e:1d:fa:f5:29
   ```

   Возвращаемый сертификат является чисто информативным. Закрытый ключ безопасно хранится внутри Stronghold.

1. Обновите местоположение CRL и выпускающие сертификаты. Эти значения могут быть обновлены в будущем.

   ```shell
   $ d8 stronghold write pki/config/urls \
        issuing_certificates="http://127.0.0.1:8200/v1/pki/ca" \
        crl_distribution_points="http://127.0.0.1:8200/v1/pki/crl"
   Success! Data written to: pki/config/urls
   ```

1. Настройте роль, которая сопоставляет имя в Stronghold с процедурой генерации сертификата. Когда пользователи или машины генерируют учетные данные, они генерируются для этой роли:

   ```shell
   $ d8 stronghold write pki/roles/example-dot-ru \
        allowed_domains=my-website.ru \
        allow_subdomains=true \
        max_ttl=72h
   Success! Data written to: pki/roles/example-dot-ru
   ```

## Использование

После того как механизм секретов настроен и у пользователя/машины есть Stronghold-токен с соответствующими правами, можно генерировать учетные данные.

1. Сгенерируйте новые учетные данные, записав их в путь `/issue` с именем роли:

   ```shell
   $ d8 stronghold write pki/issue/example-dot-ru \
        common_name=www.my-website.ru
   
   Key                 Value
   ---                 -----
   certificate         -----BEGIN CERTIFICATE-----...
   issuing_ca          -----BEGIN CERTIFICATE-----...
   private_key         -----BEGIN RSA PRIVATE KEY-----...
   private_key_type    rsa
   serial_number       1d:2e:c6:06:45:18:60:0e:23:d6:c5:17:43:c0:fe:46:ed:d1:50:be
   ```

   Вывод будет включать динамически сгенерированный закрытый ключ и сертификат, который соответствует данной роли и истекает через 72 часа (как указано в нашем определении роли). Также возвращаются выпускающий CA и цепочка доверия для упрощения автоматизации.
