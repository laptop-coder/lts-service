import type { Component } from "solid-js";
import "../../../app/styles.css";

export const ModeratorPage: Component = () => {
  document.title = "Модератор: " + document.title;
  return (
    <div class="page">
      <div class="header">
        <div class="header__wrapper">
          <img
            class="header__logo"
            src="/logo.svg"
          />
          <div class="header__title">Сервис поиска потерянных вещей</div>
        </div>
        <div class="header__title">Модератор</div>
      </div>
      <div class="box"></div>
    </div>
  );
};
