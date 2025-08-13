import { JSX } from 'solid-js';

import FormatUTCDatetime from '../FormatUTCDatetime/FormatUTCDatetime';
import Thing from '../Thing/Thing';
import ThingDescriptionTitle from '../ThingDescriptionTitle/ThingDescriptionTitle';
import ThingDescriptionGroup from '../ThingDescriptionGroup/ThingDescriptionGroup';
import ThingDescriptionItem from '../ThingDescriptionItem/ThingDescriptionItem';
import type foundThing from '../../types/foundThing';
import type utcDatetime from '../../types/utcDatetime';

const FoundThing = (props: foundThing): JSX.Element => {
  return (
    <Thing>
      <ThingDescriptionTitle
        thingId={props.FoundThingId}
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

export default FoundThing;
