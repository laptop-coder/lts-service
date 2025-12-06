import {
  JSX,
  createResource,
  createEffect,
  createSignal,
  Switch,
  Match,
} from 'solid-js';

import type { ResourceReturn } from 'solid-js';
import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, HeaderButton } from '../utils/consts';
import { Thing } from '../types/thing';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import ThingContainer from '../components/ThingContainer/ThingContainer';
import getThingData from '../utils/getThingData';
import Loading from '../ui/Loading/Loading';
import Error from '../ui/Error/Error';

import { useSearchParams } from '@solidjs/router';

const UserThingStatusPage = (): JSX.Element => {
  const [searchParams] = useSearchParams();
  const thingId = (searchParams.thing_id || '').toString();

  const [data, setData] = createSignal();
  const [state, setState] = createSignal();
  createEffect(() => {
    const [thingDataResource]: ResourceReturn<Thing> = createResource(
      {
        thingId: thingId,
      },
      getThingData,
    );
    createEffect(() => {
      setData(thingDataResource());
      setState(thingDataResource.state);
    });
  });

  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);
  return (
    <Page
      role={Role.user}
      authorized={authorized()}
    >
      <Header
        role={Role.user}
        buttons={[authorized() ? HeaderButton.profile : HeaderButton.login]}
      />
      <Content>
        {/*TODO: is it normal to use Loading in the fallback here?*/}
        <Switch fallback={<Loading />}>
          <Match when={state() === 'unresolved' || state() === 'pending'}>
            <Loading />
          </Match>
          <Match when={state() === 'ready' || state() === 'refreshing'}>
            <ThingContainer
              status
              thing={data() as Thing}
              role={Role.none}
            />
          </Match>
          <Match when={state() === 'errored'}>
            <Error />
          </Match>
        </Switch>
      </Content>
      <Footer />
    </Page>
  );
};

export default UserThingStatusPage;
