---
title: Управление control plane
permalink: ru/stronghold/documentation/admin/platform-management/control-plane-settings/about.html
lang: ru
---

Управление компонентами control plane кластера осуществляется с помощью модуля `control-plane-manager`, который запускается на всех master-узлах кластера (узлы с меткой `node-role.kubernetes.io/control-plane: ""`).

Функционал управления control plane:

- **Управление сертификатами**, необходимыми для работы control plane, в том числе продление, выпуск при изменении конфигурации и т. п. Позволяет автоматически поддерживать безопасную конфигурацию control plane и быстро добавлять дополнительные SAN для организации защищенного доступа к API Kubernetes.
- **Настройка компонентов**. Автоматически создает необходимые конфигурации и манифесты компонентов control plane.
- **Upgrade или downgrade компонентов**. Поддерживает в кластере одинаковые версии компонентов.
- **Управление конфигурацией etcd-кластера** и его узлов. Масштабирует master-узлы, выполняет миграцию из `single-master` в `multi-master` и обратно.
- **Настройка kubeconfig**. Обеспечивает всегда актуальную конфигурацию для работы kubectl. Генерирует, продлевает, обновляет kubeconfig с правами cluster-admin и создает симлинк пользователю root, чтобы kubeconfig использовался по умолчанию.
- **Расширение работы планировщика**, за счет подключения внешних плагинов через вебхуки. Управляется ресурсом [KubeSchedulerWebhookConfiguration](../../../reference/cr.html#kubeschedulerwebhookconfiguration). Позволяет использовать более сложную логику при решении задач планирования нагрузки в кластере. Например:
  - размещение подов приложений организации хранилища данных ближе к самим данным;
  - приоритизация узлов в зависимости от их состояния (сетевой нагрузки, состояния подсистемы хранения и т. д.);
  - разделение узлов на зоны.

## Управление сертификатами

Управление SSL-сертификатами компонентов control plane:

- Серверными сертификатами для `kube-apiserver` и `etcd`. Они хранятся в секрете `d8-pki` пространства имен `kube-system`:
  - корневой CA Kubernetes (`ca.crt` и `ca.key`);
  - корневой CA etcd (`etcd/ca.crt` и `etcd/ca.key`);
  - RSA-сертификат и ключ для подписи Service Account (`sa.pub` и `sa.key`);
  - корневой CA для extension API-серверов (`front-proxy-ca.key` и `front-proxy-ca.crt`).
- Клиентскими сертификатами для подключения компонентов `control-plane` друг к другу. Выписывает, продлевает и перевыписывает, если что-то изменилось (например, список SAN). Следующие сертификаты хранятся только на узлах:
  - серверный сертификат API-сервера (`apiserver.crt` и `apiserver.key`);
  - клиентский сертификат для подключения `kube-apiserver` к `kubelet` (`apiserver-kubelet-client.crt` и `apiserver-kubelet-client.key`);
  - клиентский сертификат для подключения `kube-apiserver` к `etcd` (`apiserver-etcd-client.crt` и `apiserver-etcd-client.key`);
  - клиентский сертификат для подключения `kube-apiserver` к extension API-серверам (`front-proxy-client.crt` и `front-proxy-client.key`);
  - серверный сертификат `etcd` (`etcd/server.crt` и `etcd/server.key`);
  - клиентский сертификат для подключения `etcd` к другим членам кластера (`etcd/peer.crt` и `etcd/peer.key`);
  - клиентский сертификат для подключения `kubelet` к `etcd` для helthcheck'ов (`etcd/healthcheck-client.crt` и `etcd/healthcheck-client.key`).

Также позволяет добавить дополнительные SAN в сертификаты, это дает возможность быстро и просто добавлять дополнительные «точки входа» в API Kubernetes.

При изменении сертификатов также автоматически обновляется соответствующая конфигурация kubeconfig.

## Масштабирование

Поддерживается работа control plane в конфигурации как `single-master`, так и `multi-master`.

В конфигурации `single-master`:

- `kube-apiserver` использует только тот экземпляр `etcd`, который размещен с ним на одном узле;
- На узле настраивается прокси-сервер, отвечающий на localhost,`kube-apiserver` отвечает на IP-адрес master-узла.

В конфигурации `multi-master` компоненты control plane автоматически разворачиваются в отказоустойчивом режиме:

- `kube-apiserver` настраивается для работы со всеми экземплярами `etcd`.
- На каждом master-узле настраивается дополнительный прокси-сервер, отвечающий на localhost. Прокси-сервер по умолчанию обращается к локальному экземпляру `kube-apiserver`, но в случае его недоступности последовательно опрашивает остальные экземпляры `kube-apiserver`.

### Масштабирование master-узлов

Масштабирование узлов `control-plane` осуществляется автоматически, с помощью метки `node-role.kubernetes.io/control-plane=""`:

- Установка метки `node-role.kubernetes.io/control-plane=""` на узле приводит к развертыванию на нем компонентов `control-plane`, подключению нового узла `etcd` в etcd-кластер, а также перегенерации необходимых сертификатов и конфигурационных файлов.
- Удаление метки `node-role.kubernetes.io/control-plane=""` с узла приводит к удалению всех компонентов `control-plane`, перегенерации необходимых конфигурационных файлов и сертификатов, а также корректному исключению узла из etcd-кластера.

> **Внимание.** При масштабировании узлов с 2 до 1 требуются [ручные действия](./faq.html#что-делать-если-кластер-etcd-развалился) с `etcd`. В остальных случаях все необходимые действия происходят автоматически. Обратите внимание, что при масштабировании с любого количества master-узлов до 1 рано или поздно на последнем шаге возникнет ситуация масштабирования узлов с 2 до 1.

## Управление версиями

Обновление **patch-версии** компонентов control plane (то есть в рамках минорной версии, например с `1.27.3` на `1.27.5`) происходит автоматически вместе с обновлением версии Deckhouse. Управлять обновлением patch-версий нельзя.

Обновлением **минорной-версии** компонентов control plane (например, с `1.26.*` на `1.28.*`) можно управлять с помощью параметра [kubernetesVersion](../../installing/configuration.html#clusterconfiguration-kubernetesversion), в котором можно выбрать автоматический режим обновления (значение `Automatic`) или указать желаемую минорную версию control plane. Версию control plane, которая используется по умолчанию (при `kubernetesVersion: Automatic`), а также список поддерживаемых версий Kubernetes можно найти в [документации](../../supported_versions.html#kubernetes).

Обновление control plane выполняется безопасно и для кластера с одним master-узлом, и для мультикластерных. Во время обновления может быть кратковременная недоступность API-сервера. На работу приложений в кластере обновление не влияет и может выполняться без выделения окна для регламентных работ.

Если указанная для обновления версия (параметр [kubernetesVersion](../../installing/configuration.html#clusterconfiguration-kubernetesversion)) не соответствует текущей версии control plane в кластере, запускается умная стратегия изменения версий компонентов:

- Общие замечания:
  - Обновление в разных NodeGroup выполняется параллельно. Внутри каждой NodeGroup узлы обновляются последовательно, по одному.
- При upgrade:
  - Обновление происходит **последовательными этапами**, по одной минорной версии: 1.26 -> 1.27, 1.27 -> 1.28, 1.28 -> 1.29.
  - На каждом этапе сначала обновляется версия control plane, затем происходит обновление kubelet на узлах кластера.
- При downgrade:
  - Успешный downgrade гарантируется только на одну версию вниз от максимальной минорной версии control plane, когда-либо использовавшейся в кластере.
  - Сначала происходит downgrade kubelet'a на узлах кластера, затем — downgrade компонентов control plane.

## Аудит

Если требуется журналировать операции с API или отдебажить неожиданное поведение, для этого в Kubernetes предусмотрен [Auditing](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/). Его можно настроить путем создания правил [Audit Policy](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/#audit-policy), а результатом работы аудита будет лог-файл `/var/log/kube-audit/audit.log` со всеми интересующими операциями.

В установках Deckhouse по умолчанию созданы базовые политики, которые отвечают за логирование событий:

- связанных с операциями создания, удаления и изменения ресурсов;
- совершаемых от имен сервисных аккаунтов из системных пространств имен `kube-system`, `d8-*`;
- совершаемых с ресурсами в системных пространствах имен `kube-system`, `d8-*`.

Для выключения базовых политик установите флаг [basicAuditPolicyEnabled](configuration.html#parameters-apiserver-basicauditpolicyenabled) в `false`.

Настройка политик аудита подробно рассмотрена в [одноименной секции FAQ](faq.html#как-настроить-дополнительные-политики-аудита).
