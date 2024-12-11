import type { Component } from 'solid-js';

import styles from './addnewlostthing.module.css';


interface AddNewLostThingProps {
  onClick: func;
}


const AddNewLostThing: Component = (props: AddNewLostThingProps) => {
  return (
    <div class={styles.wrapper}>
      <div class={styles.box}>
      </div>
      <div class={styles.background} onClick={props.onClick}>
      </div>
    </div>
  );
}


export default AddNewLostThing;
