import { JSX, createSignal } from 'solid-js';

import { A } from '@solidjs/router';

import AuthButton from './AuthButton/AuthButton';
import AuthForm from './AuthForm/AuthForm';
import AuthFormAnotherAction from './AuthFormAnotherAction/AuthFormAnotherAction';
import AuthInput from './AuthInput/AuthInput';
import { MODERATOR_REGISTER_ROUTE } from '../utils/consts';

const ModeratorLoginForm = (): JSX.Element => {
  const [username, setUsername] = createSignal('');
  const [email, setEmail] = createSignal('');
  const [password, setPassword] = createSignal('');
  const login = () => {};
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
        placeholder='Email'
        id='email'
        type='email'
        value={email()}
        onchange={(event) => setEmail(event.target.value)}
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
