import { JSX, Switch, Match } from 'solid-js';

import FormatUTCDatetime from '../FormatUTCDatetime/FormatUTCDatetime';
import Thing from '../Thing/Thing';
import ThingDescriptionTitle from '../ThingDescriptionTitle/ThingDescriptionTitle';
import ThingDescriptionGroup from '../ThingDescriptionGroup/ThingDescriptionGroup';
import ThingDescriptionItem from '../ThingDescriptionItem/ThingDescriptionItem';
import type FoundThing from '../../types/FoundThing';
import type utcDatetime from '../../types/utcDatetime';
import ThingStatus from '../ThingStatus/ThingStatus';

const FoundThingStatusContainer = (props: FoundThing): JSX.Element => {
  return (
    <Thing>
      <ThingDescriptionTitle
        thingId={props.FoundThingId}
        thingName={props.ThingName}
      />
      <ThingDescriptionGroup>
        <ThingDescriptionItem>
          <img
            src='/src/assets/status.svg'
            title='Статус объявления'
          />
          <ThingStatus
            verified={props.Verified}
            status={props.Status}
          />
        </ThingDescriptionItem>
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
            src='/src/assets/location.svg'
            title='Где забрать'
          />
          {props.ThingLocation}
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
      </ThingDescriptionGroup>
    </Thing>
  );
};

export default FoundThingStatusContainer;
