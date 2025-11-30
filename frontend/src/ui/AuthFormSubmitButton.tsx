import { JSX } from 'solid-js';

import SubmitButton from './SubmitButton/SubmitButton';

const AuthFormSubmitButton = (props: {
  title: string;
  name: string;
}): JSX.Element => (
  <SubmitButton
    title={props.title}
    name={props.name}
  >
    Отправить
  </SubmitButton>
);

export default AuthFormSubmitButton;
