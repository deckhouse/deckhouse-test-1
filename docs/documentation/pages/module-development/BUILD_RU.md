---
title: "Сборка и публикация модуля"
permalink: ru/module-development/build/
lang: ru
---

Deckhouse Kubernetes Platform (DKP) использует container registry для загрузки и обновления модуля. В container registry хранятся артефакты модуля. Артефакты модуля появляются в результате сборки модуля, после чего их можно загрузить (опубликовать) в registry.

## Состав артефактов модуля

В результате сборки модуля создаются три типа артефактов, которые впоследствии загружаются в container registry:
- **Образы контейнеров приложений**. Правила сборки и исходный код таких образов находятся в директории с _именем приложения_ внутри [images](../structure/#images). Собранные образы указываются в шаблонах и запускаются в кластере. Образы тегируются [content-based тегами](https://werf.io/documentation/v1.2/usage/build/process.html#tagging-images), и для их использования в шаблонах необходимо подключить библиотеку [lib-helm](https://github.com/deckhouse/lib-helm).
- **Образ модуля**. Правила сборки модуля находятся в файле `werf.yaml` в директории модуля. В качестве тегов образов используется [семантическое версионирование](https://semver.org/lang/ru/).
- **Релиз**. Артефакт версии модуля. На основе данных релиза DKP принимает решение об обновлении модуля в кластере. У релизов есть два типа тегов: тег [семантического версионирования](https://semver.org/lang/ru/) (как у образа модуля) и тег, соответствующий каналу обновлений (например, `alpha`, `beta` и т. д.). [В шаблоне модуля](https://github.com/deckhouse/modules-template/) есть пример workfklow для GitHub Actions, где релиз создается автоматически при сборке.

## Сборка артефактов модуля и публикация в container registry

Для сборки артефактов модуля и загрузки их в container registry в рамках процесса CI/CD мы предлагаем воспользоваться подготовленными [GitHub Actions](https://github.com/deckhouse/modules-actions).

[В репозитории шаблона модуля](https://github.com/deckhouse/modules-template/) представлен пример модуля, который содержит простой workflow для GitHub Actions и использует [GitHub Packages](https://github.com/features/packages) (ghcr.io) в качестве container registry. В представленном примере workflow используется следующая логика:
- Сборка артефактов модуля при изменениях в рамках PR и при слиянии изменений в ветку main.
- Сборка артефактов модуля из тегов с использованием продуктивного container registry.
- Публикация модуля в container registry GitHub Packages в выбранный [канал обновлений](../versioning/#каналы-обновлений) из тега.  

Артефакты модуля будут загружаться по адресу `ghcr.io/<OWNER>/modules/`, который будет являться [источником модулей](../../cr.html#modulesource).

Выполните следующие настройки в свойствах вашего проекта на GitHub, чтобы workflow модуля работал корректно:
- Откройте страницу _Settings -> Actions -> General_.
- Установите параметр _Read and write permissions_ в разделе _Workflow permissions_.

Вы также можете модифицировать workflow, предусмотреть использование своего container registry и более сложный процесс сборки и публикации (например, с использованием отдельных container registry для разработки и промышленной эксплуатации).

{% alert level="warning" %}
При разработке **нескольких модулей** и их публикации в GitHub Packages необходимо использовать [Personal Access Token (PAT)](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic) аккаунта.

**Не используйте** `GITHUB_TOKEN` в GitHub Workflows, чтобы избежать проблем с правами доступа при загрузке образов. Это связано с тем, что конечные релизные образы сохраняются по адресу `ghcr.io/<OWNER>/modules/`, принадлежащему первому созданному репозиторию.

Пример адаптации [шаблона модуля](https://github.com/deckhouse/modules-template/) для использования PAT:
- На странице _Settings -> Secrets and variables -> Actions_ создайте Secret с названием `TOKEN`, содержащий PAT.
- Замените переменную `GITHUB_TOKEN` на `TOKEN` в `.github/workflows/`:

    ```shell
    cd <REPO>
    sed -i -e 's/GITHUB_TOKEN/TOKEN/g' $(find .github/workflows/ -type f)
    ```

{% endalert %}

{% alert level="info" %}
Артефакты модуля также можно собрать локально с помощью [werf](https://werf.io/) (это может потребоваться, например, [при отладке](../development/)).

Вы также можете самостоятельно сделать сборку артефактов модуля и публикацию для вашей системы CI/CD по аналогии с workflow для GitHub Actions, приведенном в шаблоне модуля, но это может потребовать глубокого понимания процессов сборки и публикации модуля. Обратитесь [к сообществу](/community/) при появлении вопросов и затруднений.
{% endalert %}

Общий сценарий работы с workflow, приведенным [в шаблоне модуля](https://github.com/deckhouse/modules-template/):
1. Опубликуйте изменения в коде модуля в ветке проекта на GitHub. Это запустит сборку артефактов модуля и их публикацию в container registry.
1. Создайте новый релиз модуля или установите тег в формате [семантического версионирования](https://semver.org/lang/ru/) на нужном коммите.
1. Перейдите в раздел _Actions_ репозитория модуля на GitHub и слева, в списке workflow, выберите _Deploy_.
1. В правой части страницы нажмите на выпадающий список _Run workflow_, выберите необходимый канал обновлений и укажите нужный тег в поле ввода тега. Нажмите кнопку _Run workflow_.
1. После успешного выполнения workflow в container registry появится новая версия модуля. Опубликованную версию модуля можно [использовать в кластере](../run/).
