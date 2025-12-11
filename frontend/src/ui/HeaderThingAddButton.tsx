import { JSX } from 'solid-js';

import {
  USER__THING_ADD__ROUTE,
  ASSETS_ROUTE,
  ThingType,
} from '../utils/consts';
import HeaderButton from './HeaderButton/HeaderButton';

const HeaderThingAddButton = (props: {
  defaultThingType: ThingType;
}): JSX.Element => (
  <HeaderButton
    name='header_thing_add_button'
    title='Создать объявление'
    onclick={() =>
      (window.location.href =
        USER__THING_ADD__ROUTE +
        `?default_thing_type=${props.defaultThingType}`)
    }
    pathToImage={`${ASSETS_ROUTE}/add.svg`}
  />
);

export default HeaderThingAddButton;
