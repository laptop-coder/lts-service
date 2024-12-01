import type { Component } from 'solid-js';
import { createResource, Swith, Match } from 'solid-js';

import styles from './app.module.css';


const getThingsList = async (type: string) => {
  const response = await fetch(`http://localhost:8000/get_things_list?type=${type}`);
  return response.json();
}


const App: Component = () => {
  const [lostThingsList] = createResource("lost", getThingsList);
  const [foundThingsList] = createResource("found", getThingsList);
  return (
    <div class={styles.page}>
      <div class={styles.header}>
        <div class={styles.header__title}>
	</div>
        <div class={styles.header__buttons}>
	</div>
      </div>
      <div class={styles.content}>
        <div class={styles.things_list}>
	  <div class={styles.things_list__title}>
	  </div>
	  <Switch>
	    <Match when={lostThingsList.loading}>
	      <p>Loading...</p>
	    </Match>
	    <Match when={lostThingsList()}>
	      <div class={styles.things_list__content}>
	        {JSON.stringify(lostThingsList())}
	      </div>
	    </Match>
	  </Switch>
	</div>
        <div class={styles.things_list}>
	  <div class={styles.things_list__title}>
	  </div>
	  <Switch>
	    <Match when={foundThingsList.loading}>
	      <p>Loading...</p>
	    </Match>
	    <Match when={foundThingsList()}>
	      <div class={styles.things_list__content}>
	        {JSON.stringify(foundThingsList())}
	      </div>
	    </Match>
	  </Switch>
	</div>
      </div>
      <div class={styles.footer}>
      </div>
    </div>
  );
};

export default App;

