import type { Component } from "solid-js";
import { createResource, Swith, Match } from "solid-js";
import { createSignal } from "solid-js";

import "./styles.css";
import LostThing from "./components/LostThing";
import FoundThing from "./components/FoundThing";
import AddNewLostThing from "./components/AddNewLostThing";
import d from "./d";
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
        <div class="header__title"></div>
        <div class="header__buttons">
          <button
            class="header_button"
            onClick={() => setAddNewLostThing((prev) => !prev)}
          >
            <SVG d={d.add} />
          </button>
          <button
            class="header_button"
            onClick={() => {
              syncLostThingsList();
              syncFoundThingsList();
            }}
          >
            <SVG d={d.sync} />
          </button>
        </div>
      </div>
      <div class="content">
        <div class="things_list">
          <div class="things_list__title">Потерянные вещи</div>
          <Switch>
            <Match when={lostThingsList.loading}>
              <p>Loading...</p>
            </Match>
            <Match when={lostThingsList()}>
              <div class="things_list__content">
                {lostThingsList().map((lostThing) => (
                  <LostThing props={lostThing} />
                ))}
              </div>
            </Match>
          </Switch>
        </div>
        <div class="things_list">
          <div class="things_list__title">Найденные вещи</div>
          <Switch>
            <Match when={foundThingsList.loading}>
              <p>Loading...</p>
            </Match>
            <Match when={foundThingsList()}>
              <div class="things_list__content">
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
