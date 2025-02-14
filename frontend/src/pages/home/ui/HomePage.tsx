import type { Component } from "solid-js";
import { createSignal, createMemo, Swith, Match } from "solid-js";

import "../../../app/styles.css";
import { LostThing } from "../../../entities/lostThing/index";
import { FoundThing } from "../../../entities/foundThing/index";
import { AddNewThing } from "../../../features/add-new-thing/index";
import { d } from "../../../shared/assets/index";
import { SVG } from "../../../shared/ui/index";
import { DialogBox } from "../../../shared/ui/index";
import {
  lostThingsList,
  foundThingsList,
  syncLostThingsList,
  syncFoundThingsList,
} from "../api/getThingsLists";

export const HomePage: Component = () => {
  const [addNewThing, setAddNewThing] = createSignal(false);
  const [tabIndex, setTabIndex] = createSignal("0");
  const [rotateAddButton, setRotateAddButton] = createSignal(false);
  const [rotateSyncButton, setRotateSyncButton] = createSignal(false);
  const [lostThingsListCache, setLostThingsListCache] = createSignal();
  const [foundThingsListCache, setFoundThingsListCache] = createSignal();

  return (
    <div class="page">
      {addNewThing() && (
        <DialogBox
          actionToClose={() => {
            setAddNewThing((prev) => !prev);
            setTabIndex("0");
          }}
        >
          <AddNewThing />
        </DialogBox>
      )}
      <div class="header">
        <div class="header__title">Система поиска потерянных вещей</div>
        <div class="header__buttons">
          <button
            tabindex={tabIndex()}
            style="aspect-ratio: 1/1;"
            onClick={() => {
              setRotateAddButton(true);
              setTimeout(() => {
                setRotateAddButton(false);
              }, 1000);
              setAddNewThing((prev) => !prev);
              setTabIndex("-1");
            }}
          >
            <SVG
              d={d.add}
              class={`${rotateAddButton() ? "rotate" : ""}`}
            />
          </button>
          <button
            tabindex={tabIndex()}
            style="aspect-ratio: 1/1;"
            onClick={() => {
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
            }}
          >
            <SVG
              d={d.sync}
              class={`${rotateSyncButton() ? "rotate" : ""}`}
            />
          </button>
        </div>
      </div>
      <div class="box">
        <div
          class="list__wrapper"
          style="margin-left: 5%;"
        >
          <div class="list__title">Потерянные вещи</div>
          <Switch>
            {/*Data not loaded*/}
            <Match when={!lostThingsList() && !lostThingsListCache()}>
              <p>Loading...</p>
            </Match>
            {/*New data not loaded, old data loaded*/}
            <Match when={!lostThingsList() && lostThingsListCache()}>
              <div class="list">
                {createMemo(() => {
                  tabIndex();
                  return lostThingsListCache().map((lostThing) => (
                    <LostThing
                      tabIndex={tabIndex()}
                      props={lostThing}
                    />
                  ));
                })}
              </div>
            </Match>
            {/*New data loaded*/}
            <Match when={lostThingsList()}>
              <div class="list">
                {createMemo(() => {
                  tabIndex();
                  return lostThingsList().map((lostThing) => (
                    <LostThing
                      tabIndex={tabIndex()}
                      props={lostThing}
                    />
                  ));
                })}
              </div>
            </Match>
          </Switch>
        </div>
        <div
          class="list__wrapper"
          style="margin-right: 5%;"
        >
          <div class="list__title">Найденные вещи</div>
          <Switch>
            {/*Data not loaded*/}
            <Match when={!foundThingsList() && !foundThingsListCache()}>
              <p>Loading...</p>
            </Match>
            {/*New data not loaded, old data loaded*/}
            <Match when={!foundThingsList() && foundThingsListCache()}>
              <div class="list">
                {createMemo(() => {
                  tabIndex();
                  return foundThingsListCache().map((foundThing) => (
                    <FoundThing
                      tabIndex={tabIndex()}
                      props={foundThing}
                    />
                  ));
                })}
              </div>
            </Match>
            {/*New data loaded*/}
            <Match when={foundThingsList()}>
              <div class="list">
                {createMemo(() => {
                  tabIndex();
                  return foundThingsList().map((foundThing) => (
                    <FoundThing
                      tabIndex={tabIndex()}
                      props={foundThing}
                    />
                  ));
                })}
              </div>
            </Match>
          </Switch>
        </div>
      </div>
    </div>
  );
};
