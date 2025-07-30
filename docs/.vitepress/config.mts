import { defineConfig } from 'vitepress';

export default defineConfig({
  title: 'LTS-сервис',
  description: 'Поиск потерянных вещей',
  srcDir: 'src',
  themeConfig: {
    nav: [
      {
        text: 'Документация',
        link: '/documentation/01-intro',
      },
      { text: 'Примеры', link: '/examples/01-intro' },
    ],

    sidebar: {
      '/documentation': [
        {
          text: 'Введение',
          base: '/documentation',
          link: '/01-intro',
        },
        {
          text: 'Пользователям',
          base: '/documentation/01-for-users',
          items: [
            {
              text: 'Сценарии использования',
              link: '/01-use-cases',
            },
          ],
        },
        {
          text: 'Администраторам',
          base: '/documentation/02-for-admins',
          items: [
            {
              text: 'Установка зависимостей',
              link: '/01-installing-dependencies',
            },
            {
              text: 'Развёртывание на сервере',
              link: '/02-deploy',
            },
          ],
        },
        {
          text: 'Разработчикам',
          base: '/documentation/03-for-developers',
          items: [
            {
              text: 'Установка зависимостей',
              link: '/01-installing-dependencies',
            },
            {
              text: 'Подготовка к работе',
              link: '/02-preparation-for-work',
            },
            {
              text: 'Где хранятся данные',
              link: '/03-where-the-data-is-stored',
            },
            {
              text: 'Бэкап данных',
              link: '/04-data-backup',
            },
            {
              text: 'CI/CD',
              link: '/05-ci-cd',
            },
          ],
        },
      ],
      '/examples': [
        {
          base: '/examples',
          items: [
            {
              text: 'Введение',
              link: '/01-intro',
            },
            {
              text: 'Скриншоты',
              link: '/02-screenshots',
            },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/laptop-coder/lts-service' },
    ],
  },
});
