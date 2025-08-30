import { JSX, createSignal } from 'solid-js';

import { ASSETS_ROUTE } from '../../utils/consts';
import FormatUTCDatetime from '../FormatUTCDatetime/FormatUTCDatetime';
import ThingPhoto from '../../ui/ThingPhoto/ThingPhoto';
import Thing from '../Thing/Thing';
import ThingDescriptionTitle from '../ThingDescriptionTitle/ThingDescriptionTitle';
import ThingDescriptionGroup from '../ThingDescriptionGroup/ThingDescriptionGroup';
import ThingDescriptionItem from '../ThingDescriptionItem/ThingDescriptionItem';
import type FoundThing from '../../types/FoundThing';
import type UTCDatetime from '../../types/utcDatetime';
import ThingStatus from '../ThingStatus/ThingStatus';
import checkPhotoAvailability from '../../utils/checkPhotoAvailability';
import { STORAGE_ROUTE } from '../../utils/consts';

const FoundThingStatusContainer = (props: FoundThing): JSX.Element => {
  const pathToPhoto = `${STORAGE_ROUTE}/found/${props.FoundThingId}.jpeg`;
  const [thingPhotoIsAvailable, setThingPhotoIsAvailable] = createSignal(false);
  checkPhotoAvailability({
    pathToPhoto: pathToPhoto,
    success: () => setThingPhotoIsAvailable(true),
  });
  return (
    <Thing>
      <ThingDescriptionTitle
        thingType='found'
        thingId={props.FoundThingId}
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
            src={`${ASSETS_ROUTE}/location.svg`}
            title='Где забрать'
          />
          {props.ThingLocation}
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
            title='Изображение найденной вещи'
          />
        )}
      </ThingDescriptionGroup>
    </Thing>
  );
};

export default FoundThingStatusContainer;
