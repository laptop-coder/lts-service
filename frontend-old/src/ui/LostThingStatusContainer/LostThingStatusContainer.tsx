import { JSX, createSignal } from 'solid-js';

import { ASSETS_ROUTE } from '../../utils/consts';
import EmailLink from '../EmailLink/EmailLink';
import ThingPhoto from '../../ui/ThingPhoto/ThingPhoto';
import FormatUTCDatetime from '../FormatUTCDatetime/FormatUTCDatetime';
import Thing from '../Thing/Thing';
import ThingDescriptionTitle from '../ThingDescriptionTitle/ThingDescriptionTitle';
import ThingDescriptionGroup from '../ThingDescriptionGroup/ThingDescriptionGroup';
import ThingDescriptionItem from '../ThingDescriptionItem/ThingDescriptionItem';
import type Email from '../../types/email';
import type LostThing from '../../types/LostThing';
import type UTCDatetime from '../../types/utcDatetime';
import ThingStatus from '../ThingStatus/ThingStatus';
import checkPhotoAvailability from '../../utils/checkPhotoAvailability';
import { STORAGE_ROUTE } from '../../utils/consts';

const LostThingStatusContainer = (props: LostThing): JSX.Element => {
  const pathToPhoto = `${STORAGE_ROUTE}/lost/${props.LostThingId}.jpeg`;
  const [thingPhotoIsAvailable, setThingPhotoIsAvailable] = createSignal(false);
  checkPhotoAvailability({
    pathToPhoto: pathToPhoto,
    success: () => setThingPhotoIsAvailable(true),
  });
  return (
    <Thing>
      <ThingDescriptionTitle
        thingType='lost'
        thingId={props.LostThingId}
        thingName={props.ThingName}
      />
      <ThingDescriptionGroup>
        <ThingDescriptionItem>
          <img
            src={`${ASSETS_ROUTE}/status.svg`}
            title='Статус объявления'
          />
          <ThingStatus
            verified={props.Verified}
            status={props.Status}
          />
        </ThingDescriptionItem>
        <ThingDescriptionItem>
          <img
            src={`${ASSETS_ROUTE}/datetime.svg`}
            title='Дата и время публикации'
          />
          <FormatUTCDatetime
            datetime={props.PublicationDatetime as UTCDatetime}
          />
        </ThingDescriptionItem>
        <ThingDescriptionItem>
          <img
            src={`${ASSETS_ROUTE}/email.svg`}
            title='Email автора объявления'
          />
          <EmailLink userEmail={props.UserEmail as Email} />
        </ThingDescriptionItem>
        {props.CustomText && (
          <ThingDescriptionItem>
            <img
              src={`${ASSETS_ROUTE}/text.svg`}
              title='Сообщение автора объявления'
            />
            {props.CustomText}
          </ThingDescriptionItem>
        )}
        {thingPhotoIsAvailable() && (
          <ThingPhoto
            src={pathToPhoto}
            title='Изображение потерянной вещи'
          />
        )}
      </ThingDescriptionGroup>
    </Thing>
  );
};

export default LostThingStatusContainer;
