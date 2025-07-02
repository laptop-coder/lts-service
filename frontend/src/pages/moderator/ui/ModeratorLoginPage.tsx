import '../../../app/styles.css';
import type { Component } from 'solid-js';
import { createSignal } from 'solid-js';

export const ModeratorLoginPage: Component = () => {
  document.title = 'Модератор: вход в аккаунт (LTS-сервис)';
  const [username, setUsername] = createSignal('');
  const [password, setPassword] = createSignal('');
  return (
    <div class='page'>
      <div class='auth-container__wrapper'>
        <div class='auth-container'>
          <div class='auth-container__title'>Вход в аккаунт модератора</div>
          <form class='auth-container__form'>
            <div class='auth-container__input-group'>
              <input
                class='auth-container__input'
                id='username'
                // Placeholder is needed here to check if the input field is not empty.
                // :not(:placeholder-shown)
                placeholder=''
                value={username()}
                onInput={(event) => setUsername(event.target.value)}
              />
              <label
                class='auth-container__label'
                for='username'
              >
                Имя пользователя
              </label>
              <span class='auth-container__underline' />
            </div>
            <div class='auth-container__input-group'>
              <input
                class='auth-container__input'
                id='password'
                placeholder=''
                type='password'
                value={password()}
                onInput={(event) => setPassword(event.target.value)}
              />
              <label
                class='auth-container__label'
                for='password'
              >
                Пароль
              </label>
              <span class='auth-container__underline' />
            </div>
            <button class='auth-container__submit-button'>Войти</button>
            <span class='auth-container__another-action'>
              Нет аккаунта? <a href='/moderator/register'>Регистрация</a>
            </span>
          </form>
        </div>
      </div>
    </div>
  );
};
