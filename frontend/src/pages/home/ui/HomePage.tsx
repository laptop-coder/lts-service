import type { Component, Setter, Accessor } from 'solid-js';
import { createSignal, createMemo, Switch, Match } from 'solid-js';

import { Header } from '../../../shared/ui/index';
import type { LostThingProps } from '../../../entities/lostThing/index';
import type { FoundThingProps } from '../../../entities/foundThing/index';
import { fileToBase64 } from '../../../shared/lib/utils/index';
import { POST } from '../../../shared/lib/utils/index';
import '../../../app/styles.css';
import { HomePageLostThing } from '../../../entities/lostThing/index';
import { HomePageFoundThing } from '../../../entities/foundThing/index';
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

const handleLostThingButtonClick = (
  setChooseThingType: Setter<boolean>,
  setAddNewLostThing: Setter<boolean>,
) => {
  setChooseThingType(false);
  setAddNewLostThing(true);
};

const handleFoundThingButtonClick = (
  setChooseThingType: Setter<boolean>,
  setAddNewFoundThing: Setter<boolean>,
) => {
  setChooseThingType(false);
  setAddNewFoundThing(true);
};

const homePageKeyDown = (event: KeyboardEvent) => {
  switch (event.key) {
    case 'a':
      if (!addNewThing()) handleAddButtonClick();
      break;
    case 's':
      if (!addNewThing()) handleSyncButtonClick();
      break;
  }
};

const addNewThingKeyDown = (
  event: KeyboardEvent,
  chooseThingType: Accessor<boolean>,
  setChooseThingType: Setter<boolean>,
  setAddNewLostThing: Setter<boolean>,
  setAddNewFoundThing: Setter<boolean>,
) => {
  switch (event.key) {
    case 'l':
      if (chooseThingType())
        handleLostThingButtonClick(setChooseThingType, setAddNewLostThing);
      break;
    case 'f':
      if (chooseThingType())
        handleFoundThingButtonClick(setChooseThingType, setAddNewFoundThing);
      break;
  }
};

const [tabIndex, setTabIndex] = createSignal('0');
const [rotateAddButton, setRotateAddButton] = createSignal(false);
const [rotateSyncButton, setRotateSyncButton] = createSignal(false);
const [lostThingsListCache, setLostThingsListCache] =
  createSignal<LostThingProps[]>();
const [foundThingsListCache, setFoundThingsListCache] =
  createSignal<FoundThingProps[]>();

const [addNewThing, setAddNewThing] = createSignal(false);

const [chooseThingType, setChooseThingType] = createSignal(true);
const [addNewLostThing, setAddNewLostThing] = createSignal(false);
const [addNewFoundThing, setAddNewFoundThing] = createSignal(false);

const [thingName, setThingName] = createSignal('');
const [email, setEmail] = createSignal('');
const [thingLocation, setThingLocation] = createSignal('');
const [customText, setCustomText] = createSignal('');
const [thingPhoto, setThingPhoto] = createSignal();

const [data, setData] = createSignal({});

const [uploadPhotoFocus, setUploadPhotoFocus] = createSignal(false);

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

export const HomePage: Component = () => {
  return (
    <div
      class='page'
      tabIndex='1'
      autofocus // required for use:autofocus
      ref={autofocus}
      onKeyDown={(event) => homePageKeyDown(event)}
    >
      {addNewThing() && (
        <DialogBox
          actionToClose={() => {
            setAddNewThing((prev) => !prev);
            setTabIndex('0');
          }}
        >
          {chooseThingType() && (
            <div
              class='choose_thing_type'
              tabIndex='1'
              autofocus // required for use:autofocus
              use:autofocus
              onKeyDown={(event) =>
                addNewThingKeyDown(
                  event,
                  chooseThingType,
                  setChooseThingType,
                  setAddNewLostThing,
                  setAddNewFoundThing,
                )
              }
            >
              <button
                onClick={() =>
                  handleLostThingButtonClick(
                    setChooseThingType,
                    setAddNewLostThing,
                  )
                }
              >
                Я потерял вещь
              </button>
              <button
                onClick={() =>
                  handleFoundThingButtonClick(
                    setChooseThingType,
                    setAddNewFoundThing,
                  )
                }
              >
                Я нашёл вещь
              </button>
            </div>
          )}
          {addNewLostThing() && (
            <>
              <div class='box_title'>Добавить потерянную вещь</div>
              <form method='post'>
                <input
                  placeholder='Что Вы потеряли?*'
                  value={thingName()}
                  onInput={(event) => setThingName(event.target.value)}
                  autofocus // required for use:autofocus
                  use:autofocus
                />
                <input
                  type='email'
                  placeholder='Email*'
                  value={email()}
                  onInput={(event) => setEmail(event.target.value)}
                />
                <textarea
                  placeholder='Здесь можно оставить сообщение'
                  value={customText()}
                  onInput={(event) => setCustomText(event.target.value)}
                />
                {thingPhoto() ? (
                  <img
                    class='thing__photo'
                    src={thingPhoto()}
                    onClick={(event) => event.target.requestFullscreen()}
                  />
                ) : (
                  ''
                )}
                <input
                  type='file'
                  class='hidden'
                  id='upload-photo__input'
                  accept='image/jpeg'
                  onFocus={() => setUploadPhotoFocus((prev) => !prev)}
                  onBlur={() => setUploadPhotoFocus((prev) => !prev)}
                  tabIndex={thingPhoto() ? '-1' : '0'}
                  onInput={(event) =>
                    event.target.files &&
                    fileToBase64(event.target.files[0]).then((photoBase64) =>
                      setThingPhoto(photoBase64),
                    )
                  }
                />
                <label
                  class={`upload-photo__label${thingPhoto() ? ' hidden' : ''}${uploadPhotoFocus() ? ' focus' : ''}`}
                  for='upload-photo__input'
                >
                  Выберите файл
                </label>
                <button
                  onClick={(event) => {
                    event.preventDefault();
                    if (thingName() !== '' && email() !== '') {
                      setData({
                        thing_name: thingName(),
                        email: email(),
                        custom_text: customText(),
                        thing_photo: thingPhoto(),
                      });
                      POST('add_new_lost_thing', data()).then(() =>
                        syncLostThingsList(),
                      );
                      setAddNewThing(false);
                    } else {
                      alert('Обязательные поля не заполнены');
                    }
                  }}
                >
                  Отправить
                </button>
              </form>
            </>
          )}
          {addNewFoundThing() && (
            <>
              <div class='box_title'>Добавить найденную вещь</div>
              <form method='post'>
                <input
                  placeholder='Что Вы нашли?*'
                  value={thingName()}
                  onInput={(event) => setThingName(event.target.value)}
                  autofocus // required for use:autofocus
                  use:autofocus
                />
                <input
                  placeholder='Где забрать вещь?*'
                  value={thingLocation()}
                  onInput={(event) => setThingLocation(event.target.value)}
                />
                <textarea
                  placeholder='Здесь можно оставить сообщение'
                  value={customText()}
                  onInput={(event) => setCustomText(event.target.value)}
                />
                {thingPhoto() ? (
                  <img
                    class='thing__photo'
                    src={thingPhoto()}
                    onClick={(event) => event.target.requestFullscreen()}
                  />
                ) : (
                  ''
                )}
                <input
                  type='file'
                  class='hidden'
                  id='upload-photo__input'
                  accept='image/jpeg'
                  onFocus={() => setUploadPhotoFocus((prev) => !prev)}
                  onBlur={() => setUploadPhotoFocus((prev) => !prev)}
                  tabIndex={thingPhoto() ? '-1' : '0'}
                  onInput={(event) =>
                    event.target.files &&
                    fileToBase64(event.target.files[0]).then((photoBase64) =>
                      setThingPhoto(photoBase64),
                    )
                  }
                />
                <label
                  class={`upload-photo__label${thingPhoto() ? ' hidden' : ''}${uploadPhotoFocus() ? ' focus' : ''}`}
                  for='upload-photo__input'
                >
                  Выберите файл
                </label>
                <button
                  onClick={(event) => {
                    event.preventDefault();
                    if (thingName() !== '' && thingLocation() !== '') {
                      setData({
                        thing_name: thingName(),
                        thing_location: thingLocation(),
                        custom_text: customText(),
                        thing_photo: thingPhoto(),
                      });
                      POST('add_new_found_thing', data()).then(() =>
                        syncFoundThingsList(),
                      );
                      setAddNewThing(false);
                    } else {
                      alert('Обязательные поля не заполнены');
                    }
                  }}
                >
                  Отправить
                </button>
              </form>
            </>
          )}
        </DialogBox>
      )}
      <Header>
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
      </Header>
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
