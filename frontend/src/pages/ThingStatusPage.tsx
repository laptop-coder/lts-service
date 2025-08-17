import {
  JSX,
  createResource,
  For,
  Switch,
  Match,
  createSignal,
} from 'solid-js';

import { useSearchParams } from '@solidjs/router';

import Header from '../components/Header/Header';
import Content from '../components/Content/Content';
import List from '../components/List/List';
import ListsGroup from '../components/ListsGroup/ListsGroup';
import Footer from '../components/Footer/Footer';
import Page from '../ui/Page/Page';
import Loading from '../ui/Loading/Loading';
import Error from '../ui/Error/Error';
import SquareImageButton from '../ui/SquareImageButton/SquareImageButton';
import fetchThingData from '../utils/fetchThingData';
import type LostThing from '../types/LostThing';
import type FoundThing from '../types/FoundThing';
import LostThingStatusContainer from '../ui/LostThingStatusContainer/LostThingStatusContainer';
import FoundThingStatusContainer from '../ui/FoundThingStatusContainer/FoundThingStatusContainer';
import type { ResourceReturn } from 'solid-js'; // TODO: is it used correctly?

const ThingStatusPage = (): JSX.Element => {
  const [searchParams, setSearchParams] = useSearchParams();
  const thingId = (searchParams.thing_id || '').toString();
  const thingType = (searchParams.thing_type || '').toString();
  const [thingData, { refetch: reloadThingData }]: ResourceReturn<
    LostThing | FoundThing
  > = createResource({ thingId, thingType }, fetchThingData);
  return (
    <Page>
      <Header>
        <SquareImageButton onclick={reloadThingData}>
          <img src='/src/assets/reload.svg' />
        </SquareImageButton>
      </Header>
      <Content>
        {/*TODO: is it normal to use Loading in the fallback here?*/}
        <List>
          <Switch fallback={<Loading />}>
            <Match when={thingData() === 'Thing not found'}>
              <Error>Объявление не найдено</Error>
            </Match>
            <Match
              when={
                thingData.state === 'unresolved' ||
                thingData.state === 'pending'
              }
            >
              <Loading />
            </Match>
            <Match
              when={
                thingData.state === 'ready' || thingData.state === 'refreshing'
              }
            >
              {thingType === 'lost' && (
                <LostThingStatusContainer {...(thingData() as LostThing)} />
              )}
              {thingType === 'found' && (
                <FoundThingStatusContainer {...(thingData() as FoundThing)} />
              )}
            </Match>
            <Match when={thingData.state === 'errored'}>
              <Error />
            </Match>
          </Switch>
        </List>
      </Content>
      <Footer />
    </Page>
  );
};

export default ThingStatusPage;
