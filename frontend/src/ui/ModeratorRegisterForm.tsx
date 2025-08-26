import { JSX, createSignal } from 'solid-js';

import { A } from '@solidjs/router';

import AuthButton from './AuthButton/AuthButton';
import AuthForm from './AuthForm/AuthForm';
import AuthFormAnotherAction from './AuthFormAnotherAction/AuthFormAnotherAction';
import AuthInput from './AuthInput/AuthInput';
import { MODERATOR_LOGIN_ROUTE } from '../utils/consts';
import axiosInstanceUnauthorized from '../utils/axiosInstanceUnauthorized';

const ModeratorRegisterForm = (): JSX.Element => {
  const fieldsAreNotFilledMessage = () => alert('Не все поля заполнены');
  const sendingErrorMessage = () =>
    alert('Ошибка отправки. Попробуйте ещё раз');

  const [username, setUsername] = createSignal('');
  const [email, setEmail] = createSignal('');
  const [password, setPassword] = createSignal('');
  const [passwordRepetition, setPasswordRepetition] = createSignal('');
  const register = async (event: SubmitEvent) => {
    event.preventDefault();
    if (username() != '' && email() != '' && password() != '') {
      await axiosInstanceUnauthorized
        .post(
          '/moderator/register',
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
            window.location.replace('/moderator/login');
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
