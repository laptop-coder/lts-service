import type { Component, JSX } from 'solid-js';
import type { HeaderProps } from '../model/index';
import { children } from 'solid-js';

export const Header: Component<HeaderProps> = (props) => {
  return (
    <header>
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
      {<>{children(() => props.children)}</>}
    </header>
  );
};
