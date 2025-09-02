import { JSX, createSignal } from 'solid-js';

import { ASSETS_ROUTE } from '../../utils/consts';
import EmailLink from '../EmailLink/EmailLink';
import ThingPhoto from '../../ui/ThingPhoto/ThingPhoto';
import FormatUTCDatetime from '../FormatUTCDatetime/FormatUTCDatetime';
import Thing from '../Thing/Thing';
import ThingDescriptionTitle from '../ThingDescriptionTitle/ThingDescriptionTitle';
import ThingDescriptionGroup from '../ThingDescriptionGroup/ThingDescriptionGroup';
import ThingDescriptionItem from '../ThingDescriptionItem/ThingDescriptionItem';
import checkPhotoAvailability from '../../utils/checkPhotoAvailability';
import { STORAGE_ROUTE } from '../../utils/consts';
import type Email from '../../types/email';
import type LostThing from '../../types/LostThing';
import type FoundThing from '../../types/FoundThing';
import type utcDatetime from '../../types/utcDatetime';
import ChangeThingStatusButton from '../../ui/ChangeThingStatusButton/ChangeThingStatusButton';
import styles from './ModeratorUnauthorized.module.css';
import { A } from '@solidjs/router';

const ModeratorUnauthorized = (): JSX.Element => (
  <div class={styles.moderator_unauthorized}>
    Для получения доступа к этой странице{' '}
    <A href='/moderator/login'>войдите в аккаунт</A> модератора или{' '}
    <A href='/moderator/register'>создайте новый</A>
  </div>
);

export default ModeratorUnauthorized;
