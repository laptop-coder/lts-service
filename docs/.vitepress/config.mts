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
              text: 'Быстрый старт',
              link: '/01-quick-start',
            },
          ],
        },
        {
          text: 'Администраторам',
          base: '/documentation/02-for-admins',
          items: [
            {
              text: 'Развёртывание на сервере',
              link: '/01-deploy',
            },
          ],
        },
        {
          text: 'Разработчикам',
          base: '/documentation/03-for-developers',
          items: [
            {
              text: 'Подготовка к работе',
              link: '/01-preparation-for-work',
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
