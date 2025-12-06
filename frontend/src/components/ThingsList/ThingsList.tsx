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
import { ThingType, Role, NoticesOwnership } from '../../utils/consts';
import type { Thing } from '../../types/thing';
import ThingContainer from '../ThingContainer/ThingContainer';
import Loading from '../../ui/Loading/Loading';
import Error from '../../ui/Error/Error';
import NoData from '../../ui/NoData/NoData';

const ThingsList = (props: {
  thingsType: Accessor<ThingType>;
  role: Role;
  noticesOwnership: Accessor<NoticesOwnership>;
}): JSX.Element => {
  const [data, setData] = createSignal();
  const [state, setState] = createSignal();
  createEffect(() => {
    const [thingsListResource]: ResourceReturn<Thing> = createResource(
      {
        thingsType: props.thingsType(),
        noticesOwnership: props.noticesOwnership(),
      },
      getThingsList,
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
            each={data() as unknown as Thing[]}
            fallback={<NoData />}
          >
            {(item: Thing) => (
              <ThingContainer
                thing={item}
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
