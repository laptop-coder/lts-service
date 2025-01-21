import type { Component } from "solid-js";
import { Swith, Match } from "solid-js";
import { createSignal } from "solid-js";

import "../../../app/styles.css";
import { LostThing } from "../../../entities/lostThing/ui/LostThing";
import { FoundThing } from "../../../entities/foundThing/ui/FoundThing";
import AddNewThing from "../../../components/AddNewThing";
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

  return (
    <>
      {addNewThing() && (
        <DialogBox actionToClose={() => setAddNewThing((prev) => !prev)}>
          <AddNewThing />
        </DialogBox>
      )}
      <div class="header">
        <div class="header__title">Система поиска потерянных вещей</div>
        <div class="header__buttons">
          <button
            style="aspect-ratio: 1/1;"
            onClick={() => setAddNewThing((prev) => !prev)}
          >
            <SVG d={d.add} />
          </button>
          <button
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
                {lostThingsList().map((lostThing) => (
                  <LostThing props={lostThing} />
                ))}
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
                {foundThingsList().map((foundThing) => (
                  <FoundThing props={foundThing} />
                ))}
              </div>
            </Match>
          </Switch>
        </div>
      </div>
    </>
  );
};
