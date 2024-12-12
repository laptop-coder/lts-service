import type { Component } from 'solid-js';

import Button from './Button';
import styles from './addnewlostthing.module.css';


interface AddNewLostThingProps {
  onClick: func;
}


const AddNewLostThing: Component = (props: AddNewLostThingProps) => {
  return (
    <div class={styles.wrapper}>
      <div class={styles.box}>
        <Button onClick={() => console.log("Lost thing")} type="text" value="Я потерял вещь"/>
        <Button onClick={() => console.log("Found thing")} type="text" value="Я нашёл вещь"/>
      </div>
      <div class={styles.background} onClick={props.onClick}>
      </div>
    </div>
  );
}


export default AddNewLostThing;
