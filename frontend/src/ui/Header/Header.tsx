import { JSX, For, Switch, Match } from 'solid-js';

import { A } from '@solidjs/router';

import styles from './Header.module.css';
import Logo from '../Logo/Logo';
import {
  Role,
  HOME__ROUTE,
  MODERATOR__HOME__ROUTE,
  ThingType,
  HeaderButton,
} from '../../utils/consts';
import HeaderProfileButton from '../HeaderProfileButton/HeaderProfileButton';
import ThingAddButton from '../ThingAddButton/ThingAddButton';
import HeaderLogoutButton from '../HeaderLogoutButton/HeaderLogoutButton';
import HeaderLoginButton from '../HeaderLoginButton/HeaderLoginButton';

const Header = (props: {
  role: Role;
  addThingDefaultThingType?: ThingType;
  buttons: HeaderButton[];
}): JSX.Element => (
  <header class={styles.header}>
    <div class={styles.header_wrapper}>
      <Logo role={props.role} />
      <A
        href={
          props.role === Role.moderator ? MODERATOR__HOME__ROUTE : HOME__ROUTE
        }
        title='На главную'
      >
        <h1>
          <span style={{ color: 'var(--gray)' }}>LTS</span>
          <span style={{ color: 'var(--white)' }}>-сервис</span>
        </h1>
      </A>
    </div>
    <div class={styles.header_buttons}>
      <For each={props.buttons as HeaderButton[]}>
        {(item: HeaderButton) => (
          <Switch>
            <Match when={item === HeaderButton.add_thing}>
              <ThingAddButton
                defaultThingType={
                  props.addThingDefaultThingType !== undefined
                    ? props.addThingDefaultThingType
                    : ThingType.lost
                }
              />
            </Match>
            <Match when={item === HeaderButton.profile}>
              <HeaderProfileButton role={props.role} />
            </Match>
            <Match when={item === HeaderButton.login}>
              <HeaderLoginButton role={props.role} />
            </Match>
            <Match when={item === HeaderButton.logout}>
              <HeaderLogoutButton role={props.role} />
            </Match>
          </Switch>
        )}
      </For>
    </div>
  </header>
);

export default Header;
