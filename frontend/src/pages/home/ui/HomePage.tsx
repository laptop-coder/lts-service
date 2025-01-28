import type { Component } from "solid-js";
import { createSignal, createMemo, Swith, Match } from "solid-js";

import "../../../app/styles.css";
import { LostThing } from "../../../entities/lostThing/ui/LostThing";
import { FoundThing } from "../../../entities/foundThing/ui/FoundThing";
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

  return (
    <>
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
              setAddNewThing((prev) => !prev);
              setTabIndex("-1");
            }}
          >
            <SVG d={d.add} />
          </button>
          <button
            tabindex={tabIndex()}
            style="aspect-ratio: 1/1;"
            onClick={() => {
              syncLostThingsList();
              syncFoundThingsList();
            }}
          >
            <SVG d={d.sync} />
          </button>
        </div>
      </div>
      <div class="box">
        <div class="list__wrapper">
          <div class="list__title">Потерянные вещи</div>
          <Switch>
            <Match when={lostThingsList.loading}>
              <p>Loading...</p>
            </Match>
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
        <div class="list__wrapper">
          <div class="list__title">Найденные вещи</div>
          <Switch>
            <Match when={foundThingsList.loading}>
              <p>Loading...</p>
            </Match>
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
    </>
  );
};
