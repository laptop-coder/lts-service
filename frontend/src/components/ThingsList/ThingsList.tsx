import {
  For,
  JSX,
  createResource,
  createEffect,
  createSignal,
  Switch,
  Match,
} from 'solid-js';
import type { Accessor, ResourceReturn } from 'solid-js';

import styles from './ThingsList.module.css';
import getThingsList from '../../utils/getThingsList';
import getThingsListMy from '../../utils/getThingsListMy';
import getThingsListNotMy from '../../utils/getThingsListNotMy';
import { ThingType, Role, AdvertisementsOwnership } from '../../utils/consts';
import type { LostThing, FoundThing } from '../../types/thing';
import ThingContainer from '../ThingContainer/ThingContainer';
import Loading from '../../ui/Loading/Loading';
import Error from '../../ui/Error/Error';
import NoData from '../../ui/NoData/NoData';

const ThingsList = (props: {
  thingsType: Accessor<ThingType>;
  role: Role;
  advertisementsOwnership: Accessor<AdvertisementsOwnership>;
}): JSX.Element => {
  const [data, setData] = createSignal();
  const [state, setState] = createSignal();
  createEffect(() => {
    const [thingsListResource]: ResourceReturn<LostThing & FoundThing> =
      createResource(
        {
          thingsType: props.thingsType(),
        },
        props.advertisementsOwnership() === AdvertisementsOwnership.my
          ? getThingsListMy
          : props.advertisementsOwnership() === AdvertisementsOwnership.not_my
            ? getThingsListNotMy
            : getThingsList,
      );
    createEffect(() => {
      setData(thingsListResource());
      setState(thingsListResource.state);
    });
  });

  return (
    <div class={styles.things_list}>
      {/*TODO: is it normal to use Loading in the fallback here?*/}
      <Switch fallback={<Loading />}>
        <Match when={state() === 'unresolved' || state() === 'pending'}>
          <Loading />
        </Match>
        <Match when={state() === 'ready' || state() === 'refreshing'}>
          <For
            each={data() as unknown as (LostThing & FoundThing)[]}
            fallback={<NoData />}
          >
            {(item: LostThing & FoundThing) => (
              <ThingContainer
                thing={item}
                thingType={props.thingsType()}
                role={props.role}
              />
            )}
          </For>
        </Match>
        <Match when={state() === 'errored'}>
          <Error />
        </Match>
      </Switch>
    </div>
  );
};

export default ThingsList;
