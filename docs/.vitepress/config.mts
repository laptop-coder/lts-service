import { defineConfig } from 'vitepress';

export default defineConfig({
  title: 'LTS-сервис',
  description: 'Поиск потерянных вещей',
  srcDir: 'src',
  themeConfig: {
    nav: [
      {
        text: 'Документация',
        link: '/documentation/detailed/01-intro',
      },
      { text: 'Примеры', link: '/examples/01-intro' },
    ],

    sidebar: {
      '/documentation': [
        {
          text: 'Начало работы',
          base: '/documentation',
          link: '/getting-started',
        },
        {
          text: 'Подробная документация',
          base: '/documentation/detailed',
          items: [
            {
              text: 'Введение',
              link: '/01-intro',
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
