import type { Component } from 'solid-js';

import styles from './addnewlostthing.module.css';


interface AddNewLostThingProps {
  onClick: func;
}


const AddNewLostThing: Component = (props: AddNewLostThingProps) => {
  return (
    <div class={styles.background} onClick={props.onClick}>
    </div>
  );
}


export default AddNewLostThing;
