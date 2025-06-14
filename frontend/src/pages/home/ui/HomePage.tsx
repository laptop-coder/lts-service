import type { Component } from 'solid-js';
import { createSignal, createMemo, Switch, Match } from 'solid-js';
import type { LostThingProps } from '../../../entities/lostThing/model/LostThingProps';
import type { FoundThingProps } from '../../../entities/foundThing/model/FoundThingProps';

import '../../../app/styles.css';
import { HomePageLostThing } from '../../../entities/lostThing/index';
import { HomePageFoundThing } from '../../../entities/foundThing/index';
import { AddNewThing } from '../../../features/add-new-thing/index';
import { d } from '../../../shared/assets/index';
import { SVG } from '../../../shared/ui/index';
import { DialogBox } from '../../../shared/ui/index';
import { Loading } from '../../../shared/ui/index';
import {
  lostThingsList,
  foundThingsList,
  syncLostThingsList,
  syncFoundThingsList,
} from '../api/getThingsLists';
import { autofocus } from '@solid-primitives/autofocus';

const [addNewThing, setAddNewThing] = createSignal(false);
const [tabIndex, setTabIndex] = createSignal('0');
const [rotateAddButton, setRotateAddButton] = createSignal(false);
const [rotateSyncButton, setRotateSyncButton] = createSignal(false);
const [lostThingsListCache, setLostThingsListCache] =
  createSignal<LostThingProps[]>();
const [foundThingsListCache, setFoundThingsListCache] =
  createSignal<FoundThingProps[]>();

const handleAddButtonClick = () => {
  setRotateAddButton(true);
  setTimeout(() => {
    setRotateAddButton(false);
  }, 1000);
  setAddNewThing((prev) => !prev);
  setTabIndex('-1');
};

const handleSyncButtonClick = () => {
  setRotateSyncButton(true);
  setTimeout(() => {
    setRotateSyncButton(false);
  }, 1000);
  if (lostThingsList()) {
    setLostThingsListCache(lostThingsList());
  }
  if (foundThingsList()) {
    setFoundThingsListCache(lostThingsList());
  }
  syncLostThingsList();
  syncFoundThingsList();
};

const keyDown = (event: KeyboardEvent) => {
  switch (event.key) {
    case 'a':
      if (!addNewThing()) handleAddButtonClick();
      break;
    case 's':
      if (!addNewThing()) handleSyncButtonClick();
      break;
  }
};

export const HomePage: Component = () => {
  return (
    <div
      class='page'
      tabIndex='1'
      autofocus // required for use:autofocus
      use:autofocus
      onKeyDown={(event) => keyDown(event)}
    >
      {addNewThing() && (
        <DialogBox
          actionToClose={() => {
            setAddNewThing((prev) => !prev);
            setTabIndex('0');
          }}
        >
          <AddNewThing
            syncLostThingsList={syncLostThingsList}
            syncFoundThingsList={syncFoundThingsList}
            setAddNewThing={setAddNewThing}
          />
        </DialogBox>
      )}
      <div class='header'>
        <div class='header__wrapper'>
          <img
            class='header__logo'
            src='/logo.svg'
          />
          <div class='header__title'>Сервис поиска потерянных вещей</div>
        </div>
        <div class='header__buttons'>
          <button
            tabIndex={tabIndex()}
            style='aspect-ratio: 1/1;'
            onClick={() => handleAddButtonClick()}
          >
            <SVG
              d={d.add}
              class={`${rotateAddButton() ? 'rotate' : ''}`}
            />
          </button>
          <button
            tabIndex={tabIndex()}
            style='aspect-ratio: 1/1;'
            onClick={() => handleSyncButtonClick()}
          >
            <SVG
              d={d.sync}
              class={`${rotateSyncButton() ? 'rotate' : ''}`}
            />
          </button>
        </div>
      </div>
      <div class='box'>
        <div
          class='list__wrapper'
          style='margin-left: 5%;'
        >
          <div class='list__title'>Потерянные вещи</div>
          <Switch>
            {/*Data not loaded*/}
            <Match when={!lostThingsList() && !lostThingsListCache()}>
              <Loading />
            </Match>
            {/*New data not loaded, old data loaded*/}
            <Match when={!lostThingsList() && lostThingsListCache()}>
              <div class='list'>
                {createMemo(() => {
                  tabIndex();
                  return lostThingsListCache().map(
                    (lostThing: LostThingProps) => (
                      <HomePageLostThing
                        syncList={syncLostThingsList}
                        tabIndex={tabIndex()}
                        id={lostThing.id}
                        publication_date={lostThing.publication_date}
                        publication_time={lostThing.publication_time}
                        thing_name={lostThing.thing_name}
                        email={lostThing.email}
                        custom_text={lostThing.custom_text}
                      />
                    ),
                  );
                })()}
              </div>
            </Match>
            {/*New data loaded*/}
            <Match when={lostThingsList()}>
              <div class='list'>
                {createMemo(() => {
                  tabIndex();
                  if (lostThingsList().length === 0) {
                    return <p>Данные отсутствуют</p>;
                  } else {
                    return lostThingsList().map((lostThing: LostThingProps) => (
                      <HomePageLostThing
                        syncList={syncLostThingsList}
                        tabIndex={tabIndex()}
                        id={lostThing.id}
                        publication_date={lostThing.publication_date}
                        publication_time={lostThing.publication_time}
                        thing_name={lostThing.thing_name}
                        email={lostThing.email}
                        custom_text={lostThing.custom_text}
                      />
                    ));
                  }
                })()}
              </div>
            </Match>
          </Switch>
        </div>
        <div
          class='list__wrapper'
          style='margin-right: 5%;'
        >
          <div class='list__title'>Найденные вещи</div>
          <Switch>
            {/*Data not loaded*/}
            <Match when={!foundThingsList() && !foundThingsListCache()}>
              <Loading />
            </Match>
            {/*New data not loaded, old data loaded*/}
            <Match when={!foundThingsList() && foundThingsListCache()}>
              <div class='list'>
                {createMemo(() => {
                  tabIndex();
                  return foundThingsListCache().map(
                    (foundThing: FoundThingProps) => (
                      <HomePageFoundThing
                        syncList={syncFoundThingsList}
                        tabIndex={tabIndex()}
                        id={foundThing.id}
                        publication_date={foundThing.publication_date}
                        publication_time={foundThing.publication_time}
                        thing_name={foundThing.thing_name}
                        thing_location={foundThing.thing_location}
                        custom_text={foundThing.custom_text}
                      />
                    ),
                  );
                })()}
              </div>
            </Match>
            {/*New data loaded*/}
            <Match when={foundThingsList()}>
              <div class='list'>
                {createMemo(() => {
                  tabIndex();
                  if (foundThingsList().length === 0) {
                    return <p>Данные отсутствуют</p>;
                  } else {
                    return foundThingsList().map(
                      (foundThing: FoundThingProps) => (
                        <HomePageFoundThing
                          syncList={syncFoundThingsList}
                          tabIndex={tabIndex()}
                          id={foundThing.id}
                          publication_date={foundThing.publication_date}
                          publication_time={foundThing.publication_time}
                          thing_name={foundThing.thing_name}
                          thing_location={foundThing.thing_location}
                          custom_text={foundThing.custom_text}
                        />
                      ),
                    );
                  }
                })()}
              </div>
            </Match>
          </Switch>
        </div>
      </div>
    </div>
  );
};
