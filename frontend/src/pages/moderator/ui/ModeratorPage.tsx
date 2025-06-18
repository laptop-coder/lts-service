import type { Component } from 'solid-js';
import '../../../app/styles.css';

export const ModeratorPage: Component = () => {
  document.title = 'Модератор: ' + document.title;
  return (
    <div class='page'>
      <div class='header'>
        <div class='header__wrapper'>
          <a
            class='header__wrapper'
            href='/'
            title='На главную'
          >
            <img
              class='header__logo'
              src='/logo-512.png'
              style='cursor: pointer'
            />
            <div class='header__title'>Сервис поиска потерянных вещей</div>
          </a>
        </div>
        <div class='header__title'>Модератор</div>
      </div>
      <div class='box'></div>
    </div>
  );
};
