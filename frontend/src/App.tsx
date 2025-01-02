import type { Component } from "solid-js";
import { createResource, Swith, Match } from "solid-js";
import { createSignal } from "solid-js";

import styles from "./app.module.css";
import button_styles from "./components/button.module.css";
import LostThing from "./components/LostThing";
import FoundThing from "./components/FoundThing";
import AddNewLostThing from "./components/AddNewLostThing";
import Button from "./components/Button";
import d from "./SVG";

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
    <div class={styles.page}>
      {addNewLostThing() && (
        <AddNewLostThing onClick={() => setAddNewLostThing((prev) => !prev)} />
      )}
      <div class={styles.header}>
        <div class={styles.header__title}></div>
        <div class={styles.header__buttons}>
          <Button
            class={button_styles.square_button}
            onClick={() => setAddNewLostThing((prev) => !prev)}
            type="svg"
            value={d.add}
          />
          <Button
            class={button_styles.square_button}
            onClick={() => {
              syncLostThingsList();
              syncFoundThingsList();
            }}
            type="svg"
            value={d.sync}
          />
        </div>
      </div>
      <div class={styles.content}>
        <div class={styles.things_list}>
          <div class={styles.things_list__title}>Потерянные вещи</div>
          <Switch>
            <Match when={lostThingsList.loading}>
              <p>Loading...</p>
            </Match>
            <Match when={lostThingsList()}>
              <div class={styles.things_list__content}>
                {lostThingsList().map((lostThing) => (
                  <LostThing props={lostThing} />
                ))}
              </div>
            </Match>
          </Switch>
        </div>
        <div class={styles.things_list}>
          <div class={styles.things_list__title}>Найденные вещи</div>
          <Switch>
            <Match when={foundThingsList.loading}>
              <p>Loading...</p>
            </Match>
            <Match when={foundThingsList()}>
              <div class={styles.things_list__content}>
                {foundThingsList().map((foundThing) => (
                  <FoundThing props={foundThing} />
                ))}
              </div>
            </Match>
          </Switch>
        </div>
      </div>
      <div class={styles.footer}></div>
    </div>
  );
};

export default App;
