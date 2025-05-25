import type { Component } from 'solid-js';
import { useParams } from '@solidjs/router';
import '../../../app/styles.css';

export const StatusPage: Component = () => {
  const params = useParams();
  console.log(params.type, params.id);
  return (
    <div class='page'>
      <div class='header'>
        <a
          class='header__wrapper'
          href='/home'
          title='На главную'
        >
          <img
            class='header__logo'
            src='/logo.svg'
            style='cursor: pointer'
          />
          <div class='header__title'>Сервис поиска потерянных вещей</div>
        </a>
      </div>
    </div>
  );
};
