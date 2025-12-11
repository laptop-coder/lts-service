import { JSX } from 'solid-js';

import { ASSETS_ROUTE } from '../utils/consts';
import ThingContainerButton from './ThingContainerButton/ThingContainerButton';

const WriteToUserButton = (props: { email: string }): JSX.Element => (
  <ThingContainerButton
    title='Написать пользователю'
    name='write_to_user_button'
    onclick={() => (window.location.href = 'mailto:' + props.email)}
    pathToImage={`${ASSETS_ROUTE}/message.svg`}
    border
  />
);

export default WriteToUserButton;
