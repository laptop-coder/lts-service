import '../../../app/styles.css';
import type { Component } from 'solid-js';
import type {
  LostThingProps,
  FoundThingProps,
} from '../../../entities/thing/index';
import { ButtonHotkeyHint } from '../../../shared/ui/index';
import { DialogBox } from '../../../shared/ui/index';
import { Header } from '../../../shared/ui/index';
import { Loading } from '../../../shared/ui/index';
import { POST } from '../../../shared/lib/utils/index';
import { SVG } from '../../../shared/ui/index';
import { Thing } from '../../../entities/thing/index';
import { autofocus } from '@solid-primitives/autofocus';
import { createSignal, createMemo, Switch, Match } from 'solid-js';
import { d } from '../../../shared/assets/index';
import { fileToBase64 } from '../../../shared/lib/utils/index';
import { getThingsList } from '../api/getThingsList';

const [lostThingsList, reloadLostThingsList] = getThingsList('lost');
const [foundThingsList, reloadFoundThingsList] = getThingsList('found');

const [tabIndex, setTabIndex] = createSignal('0');
const [rotateAddButton, setRotateAddButton] = createSignal(false);
const [rotateReloadButton, setRotateReloadButton] = createSignal(false);
const [lostThingsListCache, setLostThingsListCache] =
  createSignal<LostThingProps[]>();
const [foundThingsListCache, setFoundThingsListCache] =
  createSignal<FoundThingProps[]>();

const [addNewThing, setAddNewThing] = createSignal(false);
enum AddNewThingStatuses { // type of elements is number
  ChooseThingType,
  AddLostThing,
  AddFoundThing,
}
const [currentAddNewThingStatus, setCurrentAddNewThingStatus] = createSignal(
  AddNewThingStatuses.ChooseThingType,
);

const [thingName, setThingName] = createSignal('');
const [email, setEmail] = createSignal('');
const [thingLocation, setThingLocation] = createSignal('');
const [customText, setCustomText] = createSignal('');
const [thingPhoto, setThingPhoto] = createSignal();

const [data, setData] = createSignal({});

const [uploadPhotoFocus, setUploadPhotoFocus] = createSignal(false);

const clear = () => {
  setCurrentAddNewThingStatus(AddNewThingStatuses.ChooseThingType);
  setAddNewThing(false);
  setThingName('');
  setEmail('');
  setThingLocation('');
  setCustomText('');
  setThingPhoto();
  setData({});
  setUploadPhotoFocus(false);
};

const homePageKeyDown = (event: KeyboardEvent) => {
  if (!addNewThing())
    switch (event.key) {
      case 'a':
        handleAddButtonClick();
        break;
      case 'r':
        handleReloadButtonClick();
        break;
    }
};

const addNewThingKeyDown = (event: KeyboardEvent) => {
  if (currentAddNewThingStatus() === AddNewThingStatuses.ChooseThingType)
    switch (event.key) {
      case 'l':
        setCurrentAddNewThingStatus(AddNewThingStatuses.AddLostThing);
        break;
      case 'f':
        setCurrentAddNewThingStatus(AddNewThingStatuses.AddFoundThing);
        break;
    }
};

const handleAddButtonClick = () => {
  setRotateAddButton(true);
  setTimeout(() => {
    setRotateAddButton(false);
  }, 1000);
  setAddNewThing(true);
  setTabIndex('-1');
};

const handleReloadButtonClick = () => {
  setRotateReloadButton(true);
  setTimeout(() => {
    setRotateReloadButton(false);
  }, 1000);
  if (lostThingsList()) {
    setLostThingsListCache(lostThingsList());
  }
  if (foundThingsList()) {
    setFoundThingsListCache(lostThingsList());
  }
  reloadLostThingsList();
  reloadFoundThingsList();
};

export const HomePage: Component = () => {
  return (
    <div
      class='page'
      tabIndex='1'
      autofocus // required for ref={autofocus}
      ref={autofocus}
      onKeyDown={(event) => homePageKeyDown(event)}
    >
      {addNewThing() && (
        <DialogBox
          actionToClose={() => {
            clear();
            setTabIndex('0');
          }}
        >
          {currentAddNewThingStatus() ===
            AddNewThingStatuses.ChooseThingType && (
            <div
              class='choose_thing_type'
              tabIndex='1'
              autofocus // required for ref={autofocus}
              ref={autofocus}
              onKeyDown={(event) => addNewThingKeyDown(event)}
            >
              <button
                onClick={() =>
                  setCurrentAddNewThingStatus(AddNewThingStatuses.AddLostThing)
                }
              >
                Я потерял вещь
                <ButtonHotkeyHint
                  hotkey='L'
                  place='in'
                  side='right'
                />
              </button>
              <button
                onClick={() =>
                  setCurrentAddNewThingStatus(AddNewThingStatuses.AddFoundThing)
                }
              >
                Я нашёл вещь
                <ButtonHotkeyHint
                  hotkey='F'
                  place='in'
                  side='right'
                />
              </button>
            </div>
          )}
          {currentAddNewThingStatus() === AddNewThingStatuses.AddLostThing && (
            <>
              <div class='box_title'>Добавить потерянную вещь</div>
              <form method='post'>
                <input
                  placeholder='Что Вы потеряли?*'
                  value={thingName()}
                  onInput={(event) => setThingName(event.target.value)}
                  autofocus // required for ref={autofocus}
                  ref={autofocus}
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
                      POST('add_new_lost_thing', data()).then(() => {
                        reloadLostThingsList();
                        clear();
                      });
                    } else {
                      alert('Обязательные поля не заполнены');
                    }
                  }}
                >
                  Отправить
                  <ButtonHotkeyHint
                    hotkey='Enter'
                    place='in'
                    side='right'
                  />
                </button>
              </form>
            </>
          )}
          {currentAddNewThingStatus() === AddNewThingStatuses.AddFoundThing && (
            <>
              <div class='box_title'>Добавить найденную вещь</div>
              <form method='post'>
                <input
                  placeholder='Что Вы нашли?*'
                  value={thingName()}
                  onInput={(event) => setThingName(event.target.value)}
                  autofocus // required for ref={autofocus}
                  ref={autofocus}
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
                      POST('add_new_found_thing', data()).then(() => {
                        reloadFoundThingsList();
                        clear;
                      });
                    } else {
                      alert('Обязательные поля не заполнены');
                    }
                  }}
                >
                  Отправить
                  <ButtonHotkeyHint
                    hotkey='Enter'
                    place='in'
                    side='right'
                  />
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
            <ButtonHotkeyHint
              hotkey='A'
              place='out'
              side='bottom'
            />
          </button>
          <button
            tabIndex={tabIndex()}
            style='aspect-ratio: 1/1;'
            onClick={() => handleReloadButtonClick()}
          >
            <SVG
              d={d.reload}
              class={`${rotateReloadButton() ? 'rotate' : ''}`}
            />
            <ButtonHotkeyHint
              hotkey='R'
              place='out'
              side='bottom'
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
                      <Thing
                        custom_text={lostThing.custom_text}
                        email={lostThing.email}
                        id={lostThing.id}
                        page='home'
                        publication_date={lostThing.publication_date}
                        publication_time={lostThing.publication_time}
                        reloadList={reloadLostThingsList}
                        tabIndex={tabIndex()}
                        thing_name={lostThing.thing_name}
                        type='lost'
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
                      <Thing
                        custom_text={lostThing.custom_text}
                        email={lostThing.email}
                        id={lostThing.id}
                        page='home'
                        publication_date={lostThing.publication_date}
                        publication_time={lostThing.publication_time}
                        reloadList={reloadLostThingsList}
                        tabIndex={tabIndex()}
                        thing_name={lostThing.thing_name}
                        type='lost'
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
                      <Thing
                        custom_text={foundThing.custom_text}
                        id={foundThing.id}
                        page='home'
                        publication_date={foundThing.publication_date}
                        publication_time={foundThing.publication_time}
                        reloadList={reloadFoundThingsList}
                        tabIndex={tabIndex()}
                        thing_location={foundThing.thing_location}
                        thing_name={foundThing.thing_name}
                        type='found'
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
                        <Thing
                          custom_text={foundThing.custom_text}
                          id={foundThing.id}
                          page='home'
                          publication_date={foundThing.publication_date}
                          publication_time={foundThing.publication_time}
                          reloadList={reloadFoundThingsList}
                          tabIndex={tabIndex()}
                          thing_location={foundThing.thing_location}
                          thing_name={foundThing.thing_name}
                          type='found'
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
