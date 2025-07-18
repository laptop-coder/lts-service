import '@/app/styles.css';
import type { Component } from 'solid-js';
import { Header } from '@/shared/ui/index';

export const ModeratorPage: Component = () => {
  document.title = 'Модератор: ' + document.title;
  return (
    <div class='page'>
      <Header>
        <div class='header__title'>Модератор</div>
      </Header>
      <div class='box'></div>
    </div>
  );
};
