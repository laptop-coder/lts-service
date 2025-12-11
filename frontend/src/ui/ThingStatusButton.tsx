import { JSX } from 'solid-js';

import { ASSETS_ROUTE, USER__THING_STATUS__ROUTE } from '../utils/consts';
import ThingContainerButton from './ThingContainerButton/ThingContainerButton';

const ThingStatusButton = (props: { thing: { id: string } }): JSX.Element => (
  <ThingContainerButton
    onclick={() => {
      window.location.href = `${USER__THING_STATUS__ROUTE}?thing_id=${props.thing.id}`;
    }}
    title='Статус объявления'
    name='thing_status_button'
    pathToImage={`${ASSETS_ROUTE}/info.svg`}
    border
  />
);

export default ThingStatusButton;
