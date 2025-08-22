import {
  JSX,
  createResource,
  For,
  Switch,
  Match,
  createSignal,
} from 'solid-js';

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
import LostThingContainer from '../ui/LostThingContainer/LostThingContainer';
import FoundThingContainer from '../ui/FoundThingContainer/FoundThingContainer';
import type { ResourceReturn } from 'solid-js'; // TODO: is it used correctly?
import ToggleSwitch from '../ui/ToggleSwitch/ToggleSwitch';

import { A } from '@solidjs/router';

import { ADD_THING_ROUTE } from '../utils/consts';

const HomePage = (): JSX.Element => {
  const [lostThingsList, { refetch: reloadLostThingsList }]: ResourceReturn<
    LostThing[]
  > = createResource({ thingsType: 'lost' }, fetchThingsList);
  const [foundThingsList, { refetch: reloadFoundThingsList }]: ResourceReturn<
    FoundThing[]
  > = createResource({ thingsType: 'found' }, fetchThingsList);
  const [pagination, setPagination] = createSignal(true);
  return (
    <Page>
      <Header>
        <ToggleSwitch
          title='Пагинация (разбиение списков по страницам)'
          sliderText='П'
          checked={pagination()}
          onchange={() => setPagination((prev) => !prev)}
          id='pagination'
        />
        <div style={{ display: 'flex', gap: '10px' }}>
          <SquareImageButton>
            <A href={ADD_THING_ROUTE}>
              <img src={`${ASSETS_ROUTE}/add.svg`} />
            </A>
          </SquareImageButton>
          <SquareImageButton
            onclick={() => {
              reloadLostThingsList();
              reloadFoundThingsList();
            }}
          >
            <img src={`${ASSETS_ROUTE}/reload.svg`} />
          </SquareImageButton>
        </div>
      </Header>
      <Content>
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
                    <LostThingContainer
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
                    <FoundThingContainer
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
      </Content>
      <Footer />
    </Page>
  );
};

export default HomePage;
