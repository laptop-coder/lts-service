import { JSX, createSignal } from 'solid-js';

import styles from './ThingContainer.module.css';
import type { Thing } from '../../types/thing';
import checkPhotoAvailability from '../../utils/checkPhotoAvailability';
import {
  ASSETS_ROUTE,
  STORAGE_ROUTE,
  Role,
  ThingType,
  VerificationAction,
} from '../../utils/consts';
import ThingPhoto from '../../ui/ThingPhoto/ThingPhoto';
import ThingContainerItem from '../../ui/ThingContainerItem/ThingContainerItem';
import WriteToUserButton from '../../ui/WriteToUserButton/WriteToUserButton';
import getOtherUserEmail from '../../utils/getOtherUserEmail';
import getUsername from '../../utils/getUsername';
import ThingStatusButton from '../../ui/ThingStatusButton/ThingStatusButton';
import ThingEditButton from '../../ui/ThingEditButton/ThingEditButton';
import ThingDeleteButton from '../../ui/ThingDeleteButton/ThingDeleteButton';
import FormButtonsGroup from '../../ui/FormButtonsGroup/FormButtonsGroup';
import formatDate from '../../utils/formatDate';
import ThingChangeVerificationButton from '../../ui/ThingChangeVerificationButton/ThingChangeVerificationButton';

import { Motion } from 'solid-motionone';

const ThingContainer = (props: {
  thing: Thing;
  role: Role;
  status?: boolean;
}): JSX.Element => {
  const pathToPhoto = `${STORAGE_ROUTE}/${props.thing.Id}.jpeg`;
  const [thingPhotoIsAvailable, setThingPhotoIsAvailable] = createSignal(false);
  checkPhotoAvailability({
    pathToPhoto: pathToPhoto,
    success: () => setThingPhotoIsAvailable(true),
  });

  // Username of own user
  const [username, setUsername] = createSignal('');
  // Email of notice owner
  const [noticeOwnerEmail, setNoticeOwnerEmail] = createSignal('');
  // Get username and email
  if (props.role !== Role.moderator) {
    getUsername().then((data) => setUsername(data));
    getOtherUserEmail({ username: props.thing.NoticeOwner }).then((data) =>
      setNoticeOwnerEmail(data),
    );
  }

  return (
    <Motion
      class={
        props.status ? styles.thing_container_status : styles.thing_container
      }
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.5 }}
    >
      <h2 class={styles.thing_container_title}>{props.thing.Name}</h2>
      <div class={styles.thing_container_content}>
        <ThingContainerItem
          pathToImage={`${ASSETS_ROUTE}/profile.svg`}
          title={`${props.thing.NoticeOwner} (автор объявления)`}
        >
          {props.thing.NoticeOwner}
        </ThingContainerItem>
        <ThingContainerItem
          pathToImage={`${ASSETS_ROUTE}/datetime.svg`}
          title={`${props.thing.Name} (дата и время публикации)`}
        >
          {formatDate(props.thing.PublicationDatetime)}
        </ThingContainerItem>
        {props.thing.UserMessage !== '' && (
          <ThingContainerItem
            pathToImage={`${ASSETS_ROUTE}/comment.svg`}
            title={`${props.thing.Name} (сообщение пользователя)`}
          >
            {props.thing.UserMessage}
          </ThingContainerItem>
        )}
        {props.status && (
          <>
            <ThingContainerItem
              pathToImage={`${ASSETS_ROUTE}/category.svg`}
              title={`${props.thing.Name} (тип объявления)`}
            >
              {props.thing.Type === ThingType.lost &&
                'В категории потерянных вещей'}
              {props.thing.Type === ThingType.found &&
                'В категории найденных вещей'}
            </ThingContainerItem>
            <ThingContainerItem
              pathToImage={`${ASSETS_ROUTE}/arrow_circle_right.svg`}
              title={`${props.thing.Name} (статус)`}
            >
              {props.thing.Verified ? (
                props.thing.Found ? (
                  <span class={styles.green_text}>
                    Вещь найдена, объявление снято с публикации
                  </span>
                ) : (
                  <span class={styles.yellow_text}>
                    Вещь не найдена, объявление опубликовано
                  </span>
                )
              ) : (
                <span class={styles.red_text}>
                  Объявление на модерации, не опубликовано
                </span>
              )}
            </ThingContainerItem>
          </>
        )}
        {thingPhotoIsAvailable() && (
          <ThingPhoto
            src={pathToPhoto}
            title={`${props.thing.Name} (изображение)`}
          />
        )}
        {!props.status &&
          props.role !== Role.moderator &&
          (username() === props.thing.NoticeOwner ? (
            <FormButtonsGroup>
              <ThingStatusButton thingId={props.thing.Id} />
              <ThingEditButton thingId={props.thing.Id} />
              <ThingDeleteButton
                thingId={props.thing.Id}
                thingName={props.thing.Name}
                role={props.role}
              />
            </FormButtonsGroup>
          ) : (
            <WriteToUserButton email={noticeOwnerEmail()} />
          ))}
        {props.role === Role.moderator && (
          <>
            {props.thing.Verified == '-1' && (
              <ThingChangeVerificationButton
                thingId={props.thing.Id}
                action={VerificationAction.approve}
              />
            )}
            {props.thing.Verified == '1' && (
              <ThingChangeVerificationButton
                thingId={props.thing.Id}
                action={VerificationAction.reject}
              />
            )}
            {props.thing.Verified == '0' && (
              <FormButtonsGroup>
                <ThingChangeVerificationButton
                  thingId={props.thing.Id}
                  action={VerificationAction.reject}
                />
                <ThingChangeVerificationButton
                  thingId={props.thing.Id}
                  action={VerificationAction.approve}
                />
              </FormButtonsGroup>
            )}
          </>
        )}
      </div>
    </Motion>
  );
};

export default ThingContainer;
