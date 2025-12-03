import { JSX, createSignal } from 'solid-js';

import { A } from '@solidjs/router';

import Form from '../ui/Form/Form';
import AuthFormSubmitButton from '../ui/AuthFormSubmitButton';
import Input from '../ui/Input/Input';
import {
  REGISTER_MODERATOR__ROUTE,
  BACKEND__LOGIN_MODERATOR__ROUTE,
  MODERATOR__HOME__ROUTE,
  PASSWORD_MIN_LEN,
  PASSWORD_MAX_LEN,
  USERNAME_MIN_LEN,
  USERNAME_MAX_LEN,
} from '../utils/consts';
import AuthFormOtherChoice from '../ui/AuthFormOtherChoice/AuthFormOtherChoice';
import FormTitle from '../ui/FormTitle/FormTitle';
import axiosInstance from '../utils/axiosInstance';
import { usernameRegExpStr, passwordRegExpStr } from '../utils/regExps';
import BackButton from '../ui/BackButton/BackButton';

const LoginModeratorForm = (): JSX.Element => {
  const [username, setUsername] = createSignal('');
  const [password, setPassword] = createSignal('');

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    await axiosInstance
      //TODO: refactor(move to a separate file, in utils module)
      .post(
        BACKEND__LOGIN_MODERATOR__ROUTE,
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
          window.location.replace(MODERATOR__HOME__ROUTE);
        }
      })
      .catch((error) => {
        alert(
          'Ошибка отправки. Возможно, имя пользователя или пароль неверны. Попробуйте ещё раз',
        );
        console.log(error);
      });
  };

  return (
    <Form onsubmit={handleSubmit}>
      <BackButton />
      <FormTitle>Вход в аккаунт модератора</FormTitle>
      <Input
        placeholder='Имя пользователя'
        name='login_moderator_form_username'
        value={username()}
        oninput={(event) => setUsername(event.target.value)}
        required
        pattern={usernameRegExpStr}
        minlength={USERNAME_MIN_LEN}
        maxlength={USERNAME_MAX_LEN}
      />
      <Input
        type='password'
        placeholder='Пароль'
        name='login_moderator_form_password'
        value={password()}
        oninput={(event) => setPassword(event.target.value)}
        required
        pattern={passwordRegExpStr}
        minlength={PASSWORD_MIN_LEN}
        maxlength={PASSWORD_MAX_LEN}
      />
      <AuthFormSubmitButton
        title='Войти в аккаунт'
        name='login_moderator_form_submit_button'
      />
      <AuthFormOtherChoice>
        Нет аккаунта? <A href={REGISTER_MODERATOR__ROUTE}>Регистрация</A>
      </AuthFormOtherChoice>
    </Form>
  );
};

export default LoginModeratorForm;
