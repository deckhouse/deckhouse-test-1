---
title: "Каналы обновлений"
permalink: ru/virtualization-platform/documentation/admin/release-channels.html
lang: ru
---

Deckhouse Virtualization Platform использует пять каналов обновлений, предназначенных для использования в разных окружениях, к которым с точки зрения надежности применяются разные требования:

| Канал обновлений | Описание                                                                                                                                                                                                                                                                                          |
| ---------------- |---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Alpha            | Наименее стабильный канал обновлений с наиболее частым появлением новых версий. Ориентирован на кластеры разработки, используемые небольшим количеством разработчиков.                                                                                                                            |
| Beta             | Ориентирован на кластеры разработки, как и канал обновлений Alpha. Получает версии, предварительно опробованные на канале обновлений Alpha.                                                                                                                                                       |
| Early Access     | Рекомендуемый канал обновлений, если вы не уверены в выборе. Подойдет для кластеров, в которых идет активная работа (запускаются, дорабатываются новые приложения и т. п.). Обновления функционала до этого канала обновлений доходят не ранее чем через одну неделю после их появления в релизе. |
| Stable           | Стабильный канал обновлений для кластеров, в которых закончена активная работа и преимущественно осуществляется эксплуатация. Обновления функционала до этого канала обновлений доходят не ранее чем через две недели после их появления в релизе.                                                |
| Rock Solid       | Наиболее стабильный канал обновлений. Подойдет для кластеров, которым необходимо обеспечить повышенный уровень стабильности. Обновления функционала до этого канала доходят не ранее чем через месяц после их появления в релизе.                                                                 |

Компоненты Deckhouse Virtualization Platform могут обновляться автоматически, либо с ручным подтверждением по мере выхода обновлений в каналах обновления.

Информацию по версиям, доступным на каналах обновления, можно получить на сайте [releases.deckhouse.ru](https://releases.deckhouse.ru/).

Подробнее о настройке каналов обновлений читайте в  [Обновлении платформы](./update/update.html).