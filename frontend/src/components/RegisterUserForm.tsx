import { JSX, createSignal } from 'solid-js';

import { A } from '@solidjs/router';

import Form from '../ui/Form/Form';
import AuthFormSubmitButton from '../ui/AuthFormSubmitButton';
import Input from '../ui/Input/Input';
import {
  LOGIN_USER__ROUTE,
  BACKEND__REGISTER_USER__ROUTE,
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

const RegisterUserForm = (): JSX.Element => {
  const [username, setUsername] = createSignal('');
  const [email, setEmail] = createSignal('');
  const [password, setPassword] = createSignal('');
  const [passwordRepeat, setPasswordRepeat] = createSignal('');

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    if (password() == passwordRepeat()) {
      await axiosInstance
        .post(
          BACKEND__REGISTER_USER__ROUTE,
          {
            username: username(),
            email: email(),
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
            window.location.replace(LOGIN_USER__ROUTE);
          }
        })
        .catch((error) => {
          alert('Ошибка отправки. Попробуйте ещё раз');
          console.log(error);
        });
    } else {
      alert('Пароли не совпадают');
    }
  };

  return (
    <Form onsubmit={handleSubmit}>
      <BackButton />
      <FormTitle>Создание аккаунта</FormTitle>
      <Input
        placeholder='Имя пользователя'
        name='register_user_form_username'
        value={username()}
        oninput={(event) => setUsername(event.target.value)}
        required
        pattern={usernameRegExpStr}
        minlength={USERNAME_MIN_LEN}
        maxlength={USERNAME_MAX_LEN}
      />
      <Input
        // TODO: refactor inputs (length, pattern)
        type='email'
        placeholder='Email'
        name='register_user_form_email'
        value={email()}
        oninput={(event) => setEmail(event.target.value)}
        required
      />
      <Input
        type='password'
        placeholder='Пароль'
        name='register_user_form_password'
        value={password()}
        oninput={(event) => setPassword(event.target.value)}
        required
        pattern={passwordRegExpStr}
        minlength={PASSWORD_MIN_LEN}
        maxlength={PASSWORD_MAX_LEN}
      />
      <Input
        type='password'
        placeholder='Повторите пароль'
        name='register_user_form_password_repeat'
        value={passwordRepeat()}
        oninput={(event) => setPasswordRepeat(event.target.value)}
        required
        pattern={passwordRegExpStr}
        minlength={PASSWORD_MIN_LEN}
        maxlength={PASSWORD_MAX_LEN}
      />
      <AuthFormSubmitButton
        title='Зарегистрироваться'
        name='register_user_form_submit_button'
      />
      <AuthFormOtherChoice>
        Уже есть аккаунт? <A href={LOGIN_USER__ROUTE}>Вход</A>
      </AuthFormOtherChoice>
    </Form>
  );
};

export default RegisterUserForm;
