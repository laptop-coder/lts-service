import {
  JSX,
  createResource,
  For,
  Switch,
  Match,
  createSignal,
} from 'solid-js';
import type { Setter } from 'solid-js';

import { ASSETS_ROUTE } from '../utils/consts';
import Header from '../components/Header/Header';
import Content from '../components/Content/Content';
import DoubleList from '../components/DoubleList/DoubleList';
import ListsGroup from '../components/ListsGroup/ListsGroup';
import Footer from '../components/Footer/Footer';
import Page from '../ui/Page/Page';
import Loading from '../ui/Loading/Loading';
import Error from '../ui/Error/Error';
import SquareImageButton from '../ui/SquareImageButton/SquareImageButton';
import fetchThingsList from '../utils/fetchThingsList';
import type LostThing from '../types/LostThing';
import type FoundThing from '../types/FoundThing';
import LostThingContainerModerator from '../ui/LostThingContainerModerator/LostThingContainerModerator';
import FoundThingContainerModerator from '../ui/FoundThingContainerModerator/FoundThingContainerModerator';
import ModeratorUnauthorized from '../ui/ModeratorUnauthorized/ModeratorUnauthorized';
import type { ResourceReturn } from 'solid-js'; // TODO: is it used correctly?
import { ThingsListsSelectionCriteria } from '../enums/thingsListsSelectionCriteria';
import getCookie from '../utils/getCookie';

const updateAuthorizedCookie = (setAuthorized: Setter<boolean>) => {
  var authorizedCookie = getCookie('authorized');
  if (authorizedCookie != undefined) {
    setAuthorized(JSON.parse(authorizedCookie));
  }
};

const ModeratorPage = (): JSX.Element => {
  const [lostThingsList, { refetch: reloadLostThingsList }]: ResourceReturn<
    LostThing[] | undefined
  > = createResource(
    {
      thingsType: 'lost',
      selectBy: ThingsListsSelectionCriteria.NotVerifiedThings,
    },
    fetchThingsList,
  );
  const [foundThingsList, { refetch: reloadFoundThingsList }]: ResourceReturn<
    FoundThing[] | undefined
  > = createResource(
    {
      thingsType: 'found',
      selectBy: ThingsListsSelectionCriteria.NotVerifiedThings,
    },
    fetchThingsList,
  );

  const [authorized, setAuthorized] = createSignal(false);
  updateAuthorizedCookie(setAuthorized);

  return (
    <Page>
      <Header>
        {authorized() && (
          <SquareImageButton
            onclick={() => {
              reloadLostThingsList();
              reloadFoundThingsList();
            }}
          >
            <img src={`${ASSETS_ROUTE}/reload.svg`} />
          </SquareImageButton>
        )}
      </Header>
      <Content>
        {!authorized() && <ModeratorUnauthorized />}
        {authorized() && (
          <ListsGroup>
            <DoubleList title='Потерянные вещи'>
              {/*TODO: is it normal to use Loading in the fallback here?*/}
              <Switch fallback={<Loading />}>
                <Match
                  when={
                    lostThingsList.state === 'unresolved' ||
                    lostThingsList.state === 'pending'
                  }
                >
                  <Loading />
                </Match>
                <Match
                  when={
                    lostThingsList.state === 'ready' ||
                    lostThingsList.state === 'refreshing'
                  }
                >
                  <For
                    each={lostThingsList()}
                    fallback='Данных нет'
                  >
                    {(item: LostThing) => (
                      <LostThingContainerModerator
                        {...item}
                        reloadLostThingsList={reloadLostThingsList}
                        reloadFoundThingsList={reloadFoundThingsList}
                      />
                    )}
                  </For>
                </Match>
                <Match when={lostThingsList.state === 'errored'}>
                  <Error />
                </Match>
              </Switch>
            </DoubleList>
            <DoubleList title='Найденные вещи'>
              {/*TODO: is it normal to use Loading in the fallback here?*/}
              <Switch fallback={<Loading />}>
                <Match
                  when={
                    foundThingsList.state === 'unresolved' ||
                    foundThingsList.state === 'pending'
                  }
                >
                  <Loading />
                </Match>
                <Match
                  when={
                    foundThingsList.state === 'ready' ||
                    foundThingsList.state === 'refreshing'
                  }
                >
                  <For
                    each={foundThingsList()}
                    fallback='Данных нет'
                  >
                    {(item: FoundThing) => (
                      <FoundThingContainerModerator
                        {...item}
                        reloadLostThingsList={reloadLostThingsList}
                        reloadFoundThingsList={reloadFoundThingsList}
                      />
                    )}
                  </For>
                </Match>
                <Match when={foundThingsList.state === 'errored'}>
                  <Error />
                </Match>
              </Switch>
            </DoubleList>
          </ListsGroup>
        )}
      </Content>
      <Footer />
    </Page>
  );
};

export default ModeratorPage;
