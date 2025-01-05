import type { Component } from "solid-js";
import { createResource, Swith, Match } from "solid-js";
import { createSignal } from "solid-js";

import "./styles.css";
import LostThing from "./components/LostThing";
import FoundThing from "./components/FoundThing";
import AddNewLostThing from "./components/AddNewLostThing";
import d from "./utils/d";
import SVG from "./components/SVG";

const getThingsList = async (type: string) => {
  const response = await fetch(
    `http://localhost:8000/get_things_list?type=${type}`,
  );
  return response.json();
};

const App: Component = () => {
  const [lostThingsList, { refetch: syncLostThingsList }] = createResource(
    "lost",
    getThingsList,
  );
  const [foundThingsList, { refetch: syncFoundThingsList }] = createResource(
    "found",
    getThingsList,
  );

  const [addNewLostThing, setAddNewLostThing] = createSignal(false);

  return (
    <>
      {addNewLostThing() && (
        <AddNewLostThing onClick={() => setAddNewLostThing((prev) => !prev)} />
      )}
      <div class="header">
        <div class="header__title">Система поиска потерянных вещей</div>
        <div class="header__buttons">
          <button
            style="aspect-ratio: 1/1;"
            onClick={() => setAddNewLostThing((prev) => !prev)}
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
      <div style="display: flex; justify-content: space-evenly;">
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

export default App;
