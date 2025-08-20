import { JSX, createSignal } from 'solid-js';

import { A } from '@solidjs/router';

import AuthButton from './AuthButton/AuthButton';
import AuthForm from './AuthForm/AuthForm';
import AuthFormAnotherAction from './AuthFormAnotherAction/AuthFormAnotherAction';
import AuthInput from './AuthInput/AuthInput';
import { MODERATOR_LOGIN_ROUTE } from '../utils/consts';

const ModeratorRegisterForm = (): JSX.Element => {
  const [username, setUsername] = createSignal('');
  const [password, setPassword] = createSignal('');
  const [passwordRepetition, setPasswordRepetition] = createSignal('');
  const register = () => {};
  return (
    <AuthForm
      onsubmit={register}
      title='Создание аккаунта модератора'
    >
      <AuthInput
        placeholder='Имя пользователя'
        id='username'
        value={username()}
        onchange={(event) => setUsername(event.target.value)}
      />
      <AuthInput
        placeholder='Пароль'
        id='password'
        type='password'
        value={password()}
        onchange={(event) => setPassword(event.target.value)}
      />
      <AuthInput
        placeholder='Повторите пароль'
        id='password_repetition'
        type='password'
        value={passwordRepetition()}
        onchange={(event) => setPasswordRepetition(event.target.value)}
      />
      <AuthButton>Зарегистрироваться</AuthButton>
      <AuthFormAnotherAction>
        Уже есть аккаунт? <A href={MODERATOR_LOGIN_ROUTE}>Вход</A>
      </AuthFormAnotherAction>
    </AuthForm>
  );
};

export default ModeratorRegisterForm;
