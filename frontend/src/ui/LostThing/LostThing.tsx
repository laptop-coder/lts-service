import { JSX } from 'solid-js';

import EmailLink from '../EmailLink/EmailLink';
import FormatUTCDatetime from '../FormatUTCDatetime/FormatUTCDatetime';
import Thing from '../Thing/Thing';
import ThingDescriptionTitle from '../ThingDescriptionTitle/ThingDescriptionTitle';
import ThingDescriptionGroup from '../ThingDescriptionGroup/ThingDescriptionGroup';
import ThingDescriptionItem from '../ThingDescriptionItem/ThingDescriptionItem';
import type email from '../../types/email';
import type lostThing from '../../types/lostThing';
import type utcDatetime from '../../types/utcDatetime';

const LostThing = (props: lostThing): JSX.Element => {
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
      </ThingDescriptionGroup>
    </Thing>
  );
};

export default LostThing;
