import { JSX } from 'solid-js';

import EmailLink from '../EmailLink/EmailLink';
import ThingPhoto from '../../ui/ThingPhoto/ThingPhoto';
import FormatUTCDatetime from '../FormatUTCDatetime/FormatUTCDatetime';
import Thing from '../Thing/Thing';
import ThingDescriptionTitle from '../ThingDescriptionTitle/ThingDescriptionTitle';
import ThingDescriptionGroup from '../ThingDescriptionGroup/ThingDescriptionGroup';
import ThingDescriptionItem from '../ThingDescriptionItem/ThingDescriptionItem';
import checkPhotoAvailability from '../../utils/checkPhotoAvailability';
import { STORAGE_ROUTE } from '../../utils/consts';
import type email from '../../types/email';
import type LostThing from '../../types/LostThing';
import type FoundThing from '../../types/FoundThing';
import type utcDatetime from '../../types/utcDatetime';
import ChangeThingStatusButton from '../../ui/ChangeThingStatusButton/ChangeThingStatusButton';

const LostThingContainer = (
  props: LostThing & {
    reloadLostThingsList: (
      info?: unknown,
    ) => LostThing[] | Promise<LostThing[] | undefined> | null | undefined;
    reloadFoundThingsList: (
      info?: unknown,
    ) => FoundThing[] | Promise<FoundThing[] | undefined> | null | undefined;
  },
): JSX.Element => {
  const pathToPhoto = `${STORAGE_ROUTE}/lost/${props.LostThingId}.jpeg`;
  const thingPhotoIsAvailable = checkPhotoAvailability({ pathToPhoto });
  return (
    <Thing>
      <ThingDescriptionTitle
        thingId={props.LostThingId}
        thingName={props.ThingName}
      />
      <ThingDescriptionGroup>
        <ThingDescriptionItem>
          <img
            src='/src/assets/datetime.svg'
            title='Дата и время публикации'
          />
          <FormatUTCDatetime
            datetime={props.PublicationDatetime as utcDatetime}
          />
        </ThingDescriptionItem>
        <ThingDescriptionItem>
          <img
            src='/src/assets/email.svg'
            title='Email автора объявления'
          />
          <EmailLink userEmail={props.UserEmail as email} />
        </ThingDescriptionItem>
        {props.CustomText && (
          <ThingDescriptionItem>
            <img
              src='/src/assets/text.svg'
              title='Сообщение автора объявления'
            />
            {props.CustomText}
          </ThingDescriptionItem>
        )}
        {thingPhotoIsAvailable && (
          <ThingPhoto
            src={pathToPhoto}
            title='Изображение потерянной вещи'
          />
        )}
      </ThingDescriptionGroup>
      <ChangeThingStatusButton
        thingType='lost'
        thingId={props.LostThingId}
        reloadLostThingsList={props.reloadLostThingsList}
        reloadFoundThingsList={props.reloadFoundThingsList}
      />
    </Thing>
  );
};

export default LostThingContainer;
