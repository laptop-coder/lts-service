import type { Component } from 'solid-js';

import styles from './app.module.css';


const App: Component = () => {
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
	</div>
        <div class={styles.things_list}>
	  <div class={styles.things_list__title}>
	  </div>
	</div>
      </div>
      <div class={styles.footer}>
      </div>
    </div>
  );
};

export default App;

