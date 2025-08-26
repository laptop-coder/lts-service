import { JSX, createSignal } from 'solid-js';

import { A } from '@solidjs/router';

import AuthButton from './AuthButton/AuthButton';
import AuthForm from './AuthForm/AuthForm';
import AuthFormAnotherAction from './AuthFormAnotherAction/AuthFormAnotherAction';
import AuthInput from './AuthInput/AuthInput';
import { MODERATOR_REGISTER_ROUTE } from '../utils/consts';
import axiosInstanceUnauthorized from '../utils/axiosInstanceUnauthorized';

const ModeratorLoginForm = (): JSX.Element => {
  const fieldsAreNotFilledMessage = () => alert('Не все поля заполнены');
  const sendingErrorMessage = () =>
    alert('Ошибка отправки. Попробуйте ещё раз');

  const [username, setUsername] = createSignal('');
  const [password, setPassword] = createSignal('');
  const login = async (event: SubmitEvent) => {
    event.preventDefault();
    if (username() != ''&& password() != '') {
      await axiosInstanceUnauthorized
        .post(
          '/moderator/login',
          {
            username: username(),
            password: password(),
          },
          {
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded',
            },
          },
        )
        .then((response) => {
          if (response.status == 200) {
            window.location.replace('/moderator');
          }
        })
        .catch((error) => {
          sendingErrorMessage();
          console.log(error);
        });
    } else {
      fieldsAreNotFilledMessage();
      return;
    }
  };
  return (
    <AuthForm
      onsubmit={login}
      title='Вход в аккаунт модератора'
    >
      <AuthInput
        placeholder='Имя пользователя'
        id='username'
        value={username()}
        onChange={(event) => setUsername(event.target.value)}
      />
      <AuthInput
        placeholder='Пароль'
        id='password'
        type='password'
        value={password()}
        onChange={(event) => setPassword(event.target.value)}
      />
      <AuthButton>Войти</AuthButton>
      <AuthFormAnotherAction>
        Нет аккаунта? <A href={MODERATOR_REGISTER_ROUTE}>Регистрация</A>
      </AuthFormAnotherAction>
    </AuthForm>
  );
};

export default ModeratorLoginForm;
